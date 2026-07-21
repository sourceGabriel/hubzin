package firmware

import (
	"net/http/httptest"
	"testing"

	"github.com/sourceGabriel/hubzin/backend"
)

func TestDeviceBootAndNavigation(t *testing.T) {
	s := backend.NewServer()
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	d := NewDevice("dev1", ts.URL)
	d.WiFi = WiFiConfig{SSID: "DeskHub", Password: "p"}
	if err := d.Boot(); err != nil {
		t.Fatal(err)
	}
	if d.State != StateIdle {
		t.Fatalf("expected IDLE got %s", d.State)
	}
	if d.CurrentPageName() == "" {
		t.Fatal("expected non-empty current page")
	}
	first := d.CurrentPageName()
	d.NextPage()
	if d.CurrentPageName() == first && len(d.Config.Pages) > 1 {
		t.Fatal("expected page to change")
	}
}

func TestDeviceOfflineUsesCache(t *testing.T) {
	s := backend.NewServer()
	ts := httptest.NewServer(s.Handler())
	d := NewDevice("dev1", ts.URL)
	if err := d.Boot(); err != nil {
		t.Fatal(err)
	}
	// warm cache
	if _, err := d.FetchSnapshot("weather"); err != nil {
		t.Fatal(err)
	}
	ts.Close()
	w := NewRemoteWidget("weather")
	if err := w.Update(d); err != nil {
		t.Fatalf("expected cached fallback, got %v", err)
	}
}
