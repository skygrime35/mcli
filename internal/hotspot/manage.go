// internal/hotspot/manage.go
package hotspot

import (
	"fmt"
	"os/exec"
)

// Activate brings up a Wi-Fi hotspot named ssid/password on the first
// detected Wi-Fi interface. If a connection named ssid is already
// active, this is a no-op that just reports the fact - a deliberate
// improvement over the old Python reference, which always deactivated
// and recreated the hotspot even when already active, unnecessarily
// disconnecting existing clients.
func Activate(ssid, password string) <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)

		if IsActive(ssid) {
			ch <- ProgressMsg{Text: fmt.Sprintf("Hotspot '%s' is already active.", ssid)}
			return
		}

		iface := GetWifiInterface()
		if iface == "" {
			ch <- ProgressMsg{Err: fmt.Errorf("no Wi-Fi interface found")}
			return
		}
		ch <- ProgressMsg{Text: fmt.Sprintf("Configuring hotspot on %s...", iface)}

		// Remove any stale connection profile with the same name first;
		// ignored if it doesn't exist, matching the old reference.
		_ = exec.Command("nmcli", "connection", "delete", ssid).Run()

		cmd := exec.Command("nmcli", "device", "wifi", "hotspot",
			"ifname", iface, "con-name", ssid, "ssid", ssid, "password", password)
		if err := cmd.Run(); err != nil {
			ch <- ProgressMsg{Err: fmt.Errorf("activating hotspot: %w", err)}
			return
		}
		ch <- ProgressMsg{Text: fmt.Sprintf("Hotspot activated: SSID=%s", ssid)}
	}()
	return ch
}

// Deactivate brings down the connection named ssid.
func Deactivate(ssid string) <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)
		ch <- ProgressMsg{Text: fmt.Sprintf("Deactivating hotspot '%s'...", ssid)}
		if err := exec.Command("nmcli", "connection", "down", ssid).Run(); err != nil {
			ch <- ProgressMsg{Err: fmt.Errorf("deactivating hotspot: %w", err)}
			return
		}
		ch <- ProgressMsg{Text: "Hotspot deactivated."}
	}()
	return ch
}
