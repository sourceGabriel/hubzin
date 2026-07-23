package backend

import (
	"bytes"
	"encoding/json"
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

	cfgResp, err := http.Get(ts.URL + "/v1/devices/dev1/config")
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

func TestHeartbeatUpdatesStatus(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	hb := contracts.Heartbeat{FirmwareVersion: "0.2.0", WiFiConnected: true, MQTTConnected: true, CPUUsage: 10, RAMUsage: 20}
	body, _ := json.Marshal(hb)
	resp, err := http.Post(ts.URL+"/v1/devices/dev2/heartbeat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 got %d", resp.StatusCode)
	}

	statusResp, err := http.Get(ts.URL + "/v1/devices/dev2/status")
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
	resp, err := http.Post(ts.URL+"/v1/devices/dev2/heartbeat", "application/json", bytes.NewBufferString(payload))
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
	resp, err := http.Post(ts.URL+"/v1/devices/dev2/heartbeat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
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
	resp, err := http.Post(ts.URL+"/v1/devices/dev3/events", "application/json", bytes.NewBuffer(body))
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
	resp, err := http.Post(ts.URL+"/v1/devices/dev3/events", "application/json", bytes.NewBuffer(body))
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
	resp, err := http.Post(ts.URL+"/v1/devices/dev3/events", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}
