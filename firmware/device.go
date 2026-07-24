package firmware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sourceGabriel/hubzin/pkg/contracts"
)

type DeviceState string

const (
	StateBooting  DeviceState = "BOOTING"
	StateSetup    DeviceState = "SETUP"
	StateSyncing  DeviceState = "SYNCING"
	StateIdle     DeviceState = "IDLE"
	StateSleep    DeviceState = "SLEEP"
	StateError    DeviceState = "ERROR"
	StateRecovery DeviceState = "RECOVERY"
)

type WiFiConfig struct {
	SSID     string
	Password string
}

type Widget interface {
	ID() string
	Update(*Device) error
	View() map[string]interface{}
}

type Device struct {
	mu              sync.RWMutex
	ID              string
	BackendURL      string
	State           DeviceState
	FirmwareVersion string
	WiFi            WiFiConfig
	Config          contracts.DeviceConfig
	CurrentPage     int
	DeviceSecret    string
	AccessToken     string
	TokenExpiresAt  time.Time
	Widgets         map[string]Widget
	Cache           map[string]contracts.WidgetSnapshot
	LastSync        time.Time
	BootStartedAt   time.Time
	Offline         bool
	client          *http.Client
}

func NewDevice(id, backendURL string) *Device {
	return &Device{
		ID:              id,
		BackendURL:      backendURL,
		DeviceSecret:    defaultDeviceSecret(id),
		State:           StateBooting,
		FirmwareVersion: "0.2.0",
		Widgets:         map[string]Widget{},
		Cache:           map[string]contracts.WidgetSnapshot{},
		client:          &http.Client{Timeout: 2 * time.Second},
	}

	func defaultDeviceSecret(deviceID string) string {
		switch deviceID {
		case "dev1":
			return "x"
		case "demo-device":
			return "local-only"
		default:
			return "x"
		}
	}
}

func (d *Device) Boot() error {
	d.mu.Lock()
	d.BootStartedAt = time.Now()
	d.State = StateBooting
	d.mu.Unlock()
	steps := []DeviceState{StateSetup, StateSyncing, StateIdle}
	for _, st := range steps {
		time.Sleep(5 * time.Millisecond)
		d.mu.Lock()
		d.State = st
		d.mu.Unlock()
		if st == StateSyncing {
			if err := d.FetchConfig(); err != nil {
				d.enterRecovery(err)
			}
		}
	}
	return nil
}

func (d *Device) FetchConfig() error {
	resp, err := d.doAuthorized(http.MethodGet, fmt.Sprintf("%s/v1/devices/%s/config", d.BackendURL, d.ID), nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected config status: %d", resp.StatusCode)
	}

	func (d *Device) Authenticate() error {
		payload, err := json.Marshal(contracts.TokenRequest{
			DeviceID:     d.ID,
			DeviceSecret: d.DeviceSecret,
		})
		if err != nil {
			return err
		}
		resp, err := d.client.Post(fmt.Sprintf("%s/v1/auth/token", d.BackendURL), "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected auth status: %d", resp.StatusCode)
		}
		var token contracts.TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
			return err
		}
		d.mu.Lock()
		d.AccessToken = token.AccessToken
		d.TokenExpiresAt = token.ExpiresAt
		d.mu.Unlock()
		return nil
	}

	func (d *Device) ensureAuthenticated() error {
		d.mu.RLock()
		hasValidToken := strings.TrimSpace(d.AccessToken) != "" && time.Now().Before(d.TokenExpiresAt.Add(-30*time.Second))
		d.mu.RUnlock()
		if hasValidToken {
			return nil
		}
		return d.Authenticate()
	}

	func (d *Device) doAuthorized(method, url string, body []byte, contentType string) (*http.Response, error) {
		if err := d.ensureAuthenticated(); err != nil {
			return nil, err
		}
		d.mu.RLock()
		token := d.AccessToken
		d.mu.RUnlock()
		resp, err := d.doAuthorizedWithToken(method, url, body, contentType, token)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}
		_ = resp.Body.Close()
		d.mu.Lock()
		d.AccessToken = ""
		d.TokenExpiresAt = time.Time{}
		d.mu.Unlock()
		if err := d.ensureAuthenticated(); err != nil {
			return nil, err
		}
		d.mu.RLock()
		retryToken := d.AccessToken
		d.mu.RUnlock()
		return d.doAuthorizedWithToken(method, url, body, contentType, retryToken)
	}

	func (d *Device) doAuthorizedWithToken(method, url string, body []byte, contentType, token string) (*http.Response, error) {
		var requestBody *bytes.Reader
		if body == nil {
			requestBody = bytes.NewReader([]byte{})
		} else {
			requestBody = bytes.NewReader(body)
		}
		req, err := http.NewRequest(method, url, requestBody)
		if err != nil {
			return nil, err
		}
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		return d.client.Do(req)
	}
	var cfg contracts.DeviceConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return err
	}
	d.mu.Lock()
	d.Config = cfg
	d.mu.Unlock()
	d.ensureWidgets()
	return nil
}

func (d *Device) ensureWidgets() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Widgets["clock"]; !ok {
		d.Widgets["clock"] = NewClockWidget()
	}
	if _, ok := d.Widgets["date"]; !ok {
		d.Widgets["date"] = NewDateWidget()
	}
	if _, ok := d.Widgets["weather"]; !ok {
		d.Widgets["weather"] = NewRemoteWidget("weather")
	}
	if _, ok := d.Widgets["agenda"]; !ok {
		d.Widgets["agenda"] = NewRemoteWidget("agenda")
	}
	if _, ok := d.Widgets["cpu"]; !ok {
		d.Widgets["cpu"] = NewRemoteWidget("cpu")
	}
	if _, ok := d.Widgets["spotify"]; !ok {
		d.Widgets["spotify"] = NewRemoteWidget("spotify")
	}
}

func (d *Device) Tick() {
	d.mu.RLock()
	if d.State != StateIdle {
		d.mu.RUnlock()
		return
	}
	widgets := make([]Widget, 0, len(d.Widgets))
	for _, w := range d.Widgets {
		widgets = append(widgets, w)
	}
	d.mu.RUnlock()

	for _, w := range widgets {
		if err := w.Update(d); err != nil {
			d.enterRecovery(err)
			break
		}
	}
	d.sendHeartbeat()
}

func (d *Device) sendHeartbeat() {
	hb := contracts.Heartbeat{
		UptimeSec:       int64(time.Since(d.BootStartedAt).Seconds()),
		FirmwareVersion: d.FirmwareVersion,
		WiFiConnected:   d.WiFi.SSID != "",
		MQTTConnected:   !d.Offline,
		CPUUsage:        30.0,
		RAMUsage:        40.0,
	}
	body, _ := json.Marshal(hb)
	resp, err := d.doAuthorized(http.MethodPost, fmt.Sprintf("%s/v1/devices/%s/heartbeat", d.BackendURL, d.ID), body, "application/json")
	if err != nil {
		d.enterRecovery(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		d.enterRecovery(fmt.Errorf("unexpected heartbeat status: %d", resp.StatusCode))
	}
}

func (d *Device) NextPage() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.Config.Pages) == 0 {
		return
	}
	d.CurrentPage = (d.CurrentPage + 1) % len(d.Config.Pages)
}

func (d *Device) CurrentPageName() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.Config.Pages) == 0 {
		return ""
	}
	return d.Config.Pages[d.CurrentPage].Name
}

func (d *Device) enterRecovery(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err != nil {
		d.State = StateRecovery
		d.Offline = true
	}
}

func (d *Device) Recover() {
	d.mu.Lock()
	d.State = StateSyncing
	d.mu.Unlock()
	if err := d.FetchConfig(); err != nil {
		d.enterRecovery(err)
		return
	}
	d.mu.Lock()
	d.Offline = false
	d.State = StateIdle
	d.LastSync = time.Now()
	d.mu.Unlock()
}

func (d *Device) FetchSnapshot(widgetID string) (contracts.WidgetSnapshot, error) {
	url := fmt.Sprintf("%s/v1/widgets/%s/snapshot", d.BackendURL, widgetID)
	resp, err := d.client.Get(url)
	if err != nil {
		return contracts.WidgetSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return contracts.WidgetSnapshot{}, fmt.Errorf("snapshot status: %d", resp.StatusCode)
	}
	var snap contracts.WidgetSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&snap); err != nil {
		return contracts.WidgetSnapshot{}, err
	}
	d.mu.Lock()
	d.Cache[widgetID] = snap
	d.mu.Unlock()
	return snap, nil
}

func (d *Device) CachedSnapshot(widgetID string) (contracts.WidgetSnapshot, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	s, ok := d.Cache[widgetID]
	return s, ok
}
