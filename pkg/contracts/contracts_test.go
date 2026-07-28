package contracts

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEnvelopeMarshal(t *testing.T) {
	env := Envelope{SchemaVersion: "1.0", Timestamp: time.Now(), DeviceID: "dev1", Type: "Heartbeat", Data: map[string]interface{}{"ok": true}}
	b, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("expected non-empty json")
	}
}
