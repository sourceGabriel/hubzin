package contracts

import "time"

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

type TokenRequest struct {
	DeviceID     string `json:"deviceId"`
	DeviceSecret string `json:"deviceSecret"`
}

type TokenResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type WidgetConfig struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	UpdateIntervalSec int    `json:"updateIntervalSec"`
}

type PageConfig struct {
	Name    string         `json:"name"`
	Widgets []WidgetConfig `json:"widgets"`
}

type DeviceConfig struct {
	SchemaVersion string       `json:"schemaVersion"`
	Theme         string       `json:"theme"`
	Timezone      string       `json:"timezone"`
	Pages         []PageConfig `json:"pages"`
}

type DeviceStatus struct {
	DeviceID        string    `json:"deviceId"`
	FirmwareVersion string    `json:"firmwareVersion"`
	CPUUsage        float64   `json:"cpuUsage"`
	RAMUsage        float64   `json:"ramUsage"`
	WiFiConnected   bool      `json:"wifiConnected"`
	MQTTConnected   bool      `json:"mqttConnected"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type OTAInfo struct {
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
	Signature   string `json:"signature"`
	DownloadURL string `json:"downloadUrl"`
	Rollout     string `json:"rollout"`
}

type Envelope struct {
	SchemaVersion string      `json:"schemaVersion"`
	Timestamp     time.Time   `json:"timestamp"`
	DeviceID      string      `json:"deviceId"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
}

type Heartbeat struct {
	UptimeSec       int64   `json:"uptimeSec"`
	FirmwareVersion string  `json:"firmwareVersion"`
	WiFiConnected   bool    `json:"wifiConnected"`
	MQTTConnected   bool    `json:"mqttConnected"`
	CPUUsage        float64 `json:"cpuUsage"`
	RAMUsage        float64 `json:"ramUsage"`
}

type WidgetSnapshot struct {
	WidgetID string      `json:"widgetId"`
	Data     interface{} `json:"data"`
	CachedAt time.Time   `json:"cachedAt"`
}
