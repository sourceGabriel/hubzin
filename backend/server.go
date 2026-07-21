package backend

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sourceGabriel/hubzin/pkg/contracts"
)

type Server struct {
	mu        sync.RWMutex
	config    contracts.DeviceConfig
	status    map[string]contracts.DeviceStatus
	snapshots map[string]contracts.WidgetSnapshot
	ota       contracts.OTAInfo
	events    []contracts.Envelope
}

func NewServer() *Server {
	return &Server{
		config: contracts.DeviceConfig{
			SchemaVersion: "1.0",
			Theme:         "dark",
			Timezone:      "UTC",
			Pages: []contracts.PageConfig{
				{Name: "Home", Widgets: []contracts.WidgetConfig{{ID: "clock", Name: "Clock", UpdateIntervalSec: 1}, {ID: "date", Name: "Date", UpdateIntervalSec: 1}}},
				{Name: "Work", Widgets: []contracts.WidgetConfig{{ID: "agenda", Name: "Agenda", UpdateIntervalSec: 300}, {ID: "cpu", Name: "CPU", UpdateIntervalSec: 5}}},
				{Name: "Music", Widgets: []contracts.WidgetConfig{{ID: "spotify", Name: "Spotify", UpdateIntervalSec: 5}}},
				{Name: "Weather", Widgets: []contracts.WidgetConfig{{ID: "weather", Name: "Weather", UpdateIntervalSec: 900}}},
			},
		},
		status: map[string]contracts.DeviceStatus{},
		snapshots: map[string]contracts.WidgetSnapshot{
			"weather": {WidgetID: "weather", Data: map[string]interface{}{"tempC": 24, "condition": "sunny"}, CachedAt: time.Now()},
			"agenda":  {WidgetID: "agenda", Data: []string{"Daily standup 09:00", "Review 14:00"}, CachedAt: time.Now()},
			"cpu":     {WidgetID: "cpu", Data: map[string]interface{}{"usage": 31.4}, CachedAt: time.Now()},
			"spotify": {WidgetID: "spotify", Data: map[string]interface{}{"track": "N/A", "state": "paused"}, CachedAt: time.Now()},
		},
		ota: contracts.OTAInfo{Version: "0.2.0", Checksum: "sha256-demo", Signature: "signed-demo", DownloadURL: "https://example.invalid/ota.bin", Rollout: "stable"},
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/auth/token", s.handleToken)
	mux.HandleFunc("/v1/devices/", s.handleDevices)
	mux.HandleFunc("/v1/widgets/", s.handleWidgets)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	return mux
}

func (s *Server) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "invalid method")
		return
	}
	var req contracts.TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DeviceID == "" {
		writeErr(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid token payload")
		return
	}
	writeJSON(w, http.StatusOK, contracts.TokenResponse{
		AccessToken:  "token-" + req.DeviceID,
		RefreshToken: "refresh-" + req.DeviceID,
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	})
}

func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/v1/devices/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) < 2 {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		return
	}
	deviceID, resource := parts[0], parts[1]
	switch {
	case r.Method == http.MethodGet && resource == "config":
		writeJSON(w, http.StatusOK, s.config)
	case r.Method == http.MethodGet && resource == "status":
		s.mu.RLock()
		status, ok := s.status[deviceID]
		s.mu.RUnlock()
		if !ok {
			status = contracts.DeviceStatus{DeviceID: deviceID, FirmwareVersion: "unknown", UpdatedAt: time.Now()}
		}
		writeJSON(w, http.StatusOK, status)
	case r.Method == http.MethodGet && resource == "ota" && len(parts) >= 3 && parts[2] == "latest":
		writeJSON(w, http.StatusOK, s.ota)
	case r.Method == http.MethodPost && resource == "heartbeat":
		var hb contracts.Heartbeat
		if err := json.NewDecoder(r.Body).Decode(&hb); err != nil {
			writeErr(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid heartbeat")
			return
		}
		s.mu.Lock()
		s.status[deviceID] = contracts.DeviceStatus{
			DeviceID:        deviceID,
			FirmwareVersion: hb.FirmwareVersion,
			CPUUsage:        hb.CPUUsage,
			RAMUsage:        hb.RAMUsage,
			WiFiConnected:   hb.WiFiConnected,
			MQTTConnected:   hb.MQTTConnected,
			UpdatedAt:       time.Now(),
		}
		s.mu.Unlock()
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
	case r.Method == http.MethodPost && resource == "events":
		var ev contracts.Envelope
		if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
			writeErr(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid event")
			return
		}
		ev.DeviceID = deviceID
		if ev.Timestamp.IsZero() {
			ev.Timestamp = time.Now()
		}
		s.mu.Lock()
		s.events = append(s.events, ev)
		s.mu.Unlock()
		writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
	default:
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
	}
}

func (s *Server) handleWidgets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "invalid method")
		return
	}
	trimmed := strings.TrimPrefix(r.URL.Path, "/v1/widgets/")
	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) != 2 || parts[1] != "snapshot" {
		writeErr(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		return
	}
	id := parts[0]
	s.mu.RLock()
	snap, ok := s.snapshots[id]
	s.mu.RUnlock()
	if !ok {
		writeErr(w, http.StatusNotFound, "WIDGET_NOT_FOUND", "widget snapshot not found")
		return
	}
	writeJSON(w, http.StatusOK, snap)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, contracts.ErrorResponse{Error: contracts.APIError{Code: code, Message: msg}})
}
