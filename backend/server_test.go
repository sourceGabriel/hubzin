package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sourceGabriel/hubzin/pkg/contracts"
)

func TestTokenAndConfigEndpoints(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload, _ := json.Marshal(contracts.TokenRequest{DeviceID: "dev1", DeviceSecret: "x"})
	resp, err := http.Post(ts.URL+"/v1/auth/token", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	cfgResp, err := doAuthorizedRequest(http.MethodGet, ts.URL+"/v1/devices/dev1/config", token, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cfgResp.Body.Close()
	if cfgResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", cfgResp.StatusCode)
	}
}

func TestTokenRejectsInvalidCredentials(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload, _ := json.Marshal(contracts.TokenRequest{DeviceID: "dev1", DeviceSecret: "wrong"})
	resp, err := http.Post(ts.URL+"/v1/auth/token", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", resp.StatusCode)
	}
}

func TestTokenRejectsUnknownFields(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"deviceId":"dev1","deviceSecret":"x","unexpected":"value"}`
	resp, err := http.Post(ts.URL+"/v1/auth/token", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestTokenRejectsTrailingJSON(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"deviceId":"dev1","deviceSecret":"x"}{"extra":true}`
	resp, err := http.Post(ts.URL+"/v1/auth/token", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestTokenRejectsNonJSONContentType(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"deviceId":"dev1","deviceSecret":"x"}`
	resp, err := doRequest(http.MethodPost, ts.URL+"/v1/auth/token", "text/plain", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}

	func TestDeviceEndpointsRequireBearerToken(t *testing.T) {
		s := NewServer()
		ts := httptest.NewServer(s.Handler())
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/v1/devices/dev1/config")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401 got %d", resp.StatusCode)
		}
	}

	func TestDeviceEndpointsRejectTokenFromAnotherDevice(t *testing.T) {
		s := NewServer()
		ts := httptest.NewServer(s.Handler())
		defer ts.Close()

		token := issueTokenForTest(t, ts.URL, "dev1", "x")
		resp, err := doAuthorizedRequest(http.MethodGet, ts.URL+"/v1/devices/demo-device/config", token, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401 got %d", resp.StatusCode)
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 got %d", resp.StatusCode)
	}
}

func TestHeartbeatUpdatesStatus(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()
	token := issueTokenForTest(t, ts.URL, "dev1", "x")

	hb := contracts.Heartbeat{FirmwareVersion: "0.2.0", WiFiConnected: true, MQTTConnected: true, CPUUsage: 10, RAMUsage: 20}
	body, _ := json.Marshal(hb)
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/heartbeat", token, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 got %d", resp.StatusCode)
	}

	statusResp, err := doAuthorizedRequest(http.MethodGet, ts.URL+"/v1/devices/dev1/status", token, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer statusResp.Body.Close()
	if statusResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", statusResp.StatusCode)
	}
}

func TestHeartbeatRejectsUnknownFields(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"firmwareVersion":"0.2.0","wifiConnected":true,"mqttConnected":true,"cpuUsage":10,"ramUsage":20,"unexpected":"x"}`
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/heartbeat", token, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestHeartbeatRejectsInvalidUsageRange(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	hb := contracts.Heartbeat{FirmwareVersion: "0.2.0", WiFiConnected: true, MQTTConnected: true, CPUUsage: 101, RAMUsage: 20}
	body, _ := json.Marshal(hb)
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/heartbeat", token, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestHeartbeatRejectsNonJSONContentType(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"firmwareVersion":"0.2.0","wifiConnected":true,"mqttConnected":true,"cpuUsage":10,"ramUsage":20}`
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/heartbeat", token, "text/plain", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 got %d", resp.StatusCode)
	}
}

func TestEventsEndpointAcceptsValidEvent(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	ev := contracts.Envelope{
		SchemaVersion: "1.0",
		Timestamp:     time.Now(),
		Type:          "widget.update",
		Data:          map[string]interface{}{"ok": true},
	}
	body, _ := json.Marshal(ev)
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/events", token, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 got %d", resp.StatusCode)
	}
}

func TestEventsEndpointRejectsMissingType(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	ev := contracts.Envelope{
		SchemaVersion: "1.0",
		Timestamp:     time.Now(),
		Data:          map[string]interface{}{"ok": true},
	}
	body, _ := json.Marshal(ev)
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/events", token, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestEventsEndpointRejectsUnknownFields(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"schemaVersion":"1.0","timestamp":"2026-07-23T00:00:00Z","type":"widget.update","data":{},"unexpected":"x"}`
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/events", token, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestEventsEndpointRejectsNonJSONContentType(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	payload := `{"schemaVersion":"1.0","type":"widget.update","data":{}}`
	token := issueTokenForTest(t, ts.URL, "dev1", "x")
	resp, err := doAuthorizedRequest(http.MethodPost, ts.URL+"/v1/devices/dev1/events", token, "text/plain", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 got %d", resp.StatusCode)
	}
}

func doRequest(method, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	func doAuthorizedRequest(method, url, token, contentType string, body io.Reader) (*http.Response, error) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		return http.DefaultClient.Do(req)
	}

	func issueTokenForTest(t *testing.T, baseURL, deviceID, deviceSecret string) string {
		t.Helper()
		payload, _ := json.Marshal(contracts.TokenRequest{DeviceID: deviceID, DeviceSecret: deviceSecret})
		resp, err := http.Post(baseURL+"/v1/auth/token", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected token 200 got %d", resp.StatusCode)
		}
		var token contracts.TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
			t.Fatal(err)
		}
		return token.AccessToken
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return http.DefaultClient.Do(req)
}
