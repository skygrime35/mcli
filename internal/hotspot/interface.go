// internal/hotspot/interface.go
package hotspot

import (
	"os/exec"
	"strings"
)

// parseWifiInterface scans `nmcli device status` output for the first
// line whose TYPE column is exactly "wifi", returning its device name.
func parseWifiInterface(output string) string {
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[1] == "wifi" {
			return fields[0]
		}
	}
	return ""
}

// GetWifiInterface returns the name of the first Wi-Fi device nmcli
// knows about, or "" if none is found or nmcli isn't available.
func GetWifiInterface() string {
	out, err := exec.Command("nmcli", "device", "status").Output()
	if err != nil {
		return ""
	}
	return parseWifiInterface(string(out))
}
