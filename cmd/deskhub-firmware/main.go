package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sourceGabriel/hubzin/firmware"
)

func main() {
	d := firmware.NewDevice("demo-device", "http://localhost:8080")
	d.WiFi = firmware.WiFiConfig{SSID: "DeskHubWiFi", Password: "local-only"}
	if err := d.Boot(); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for i := 0; i < 5; i++ {
		<-ticker.C
		d.Tick()
		fmt.Printf("state=%s page=%s\n", d.State, d.CurrentPageName())
		d.NextPage()
	}
}
