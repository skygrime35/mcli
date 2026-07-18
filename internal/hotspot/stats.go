// internal/hotspot/stats.go
package hotspot

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// findActiveHotspotInterface scans `nmcli -t -f DEVICE,TYPE,STATE,CONNECTION
// device` output for a connected Wi-Fi device whose connection name is
// ssid, returning the device name (e.g. "wlo1").
func findActiveHotspotInterface(deviceOutput string, ssid string) string {
	for _, line := range strings.Split(deviceOutput, "\n") {
		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}
		device, devType, state, conn := fields[0], fields[1], fields[2], fields[3]
		if devType == "wifi" && state == "connected" && conn == ssid {
			return device
		}
	}
	return ""
}

// parseNeighClients parses `ip neigh show dev <iface>` output, returning
// the IP addresses of entries that are REACHABLE, STALE, or DELAY (i.e.
// recently-seen clients) - matching the old reference's filter exactly.
func parseNeighClients(output string) []string {
	var clients []string
	for _, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
		if line == "" {
			continue
		}
		if strings.Contains(line, "REACHABLE") || strings.Contains(line, "STALE") || strings.Contains(line, "DELAY") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				clients = append(clients, fields[0])
			}
		}
	}
	return clients
}

// parseInterfaceBytes parses /proc/net/dev content for the given
// interface, returning (rxBytes, txBytes, ok). After splitting on the
// interface's ":", rx bytes is the 1st field and tx bytes is the 9th
// field of the remaining columns.
func parseInterfaceBytes(procNetDev string, iface string) (rx, tx uint64, ok bool) {
	for _, line := range strings.Split(procNetDev, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) != iface {
			continue
		}
		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			continue
		}
		rxVal, err1 := strconv.ParseUint(fields[0], 10, 64)
		txVal, err2 := strconv.ParseUint(fields[8], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		return rxVal, txVal, true
	}
	return 0, 0, false
}

// FormatBytes renders a byte count in human-readable units (B/KiB/MiB/
// GiB/TiB), matching the old reference's formatting.
func FormatBytes(n uint64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	v := float64(n)
	for _, u := range units[:len(units)-1] {
		if v < 1024 {
			return fmt.Sprintf("%.1f%s", v, u)
		}
		v /= 1024
	}
	return fmt.Sprintf("%.1f%s", v, units[len(units)-1])
}

// GetStats reports the current hotspot's active interface, connected
// clients, and data volume. Returns Stats{Active: false} if no
// connection named ssid is currently active.
func GetStats(ssid string) Stats {
	if !IsActive(ssid) {
		return Stats{Active: false}
	}

	iface := ""
	if out, err := exec.Command("nmcli", "-t", "-f", "DEVICE,TYPE,STATE,CONNECTION", "device").Output(); err == nil {
		iface = findActiveHotspotInterface(string(out), ssid)
	}
	if iface == "" {
		iface = GetWifiInterface()
	}
	if iface == "" {
		return Stats{Active: true}
	}

	stats := Stats{Active: true, Interface: iface}

	if out, err := exec.Command("ip", "neigh", "show", "dev", iface).Output(); err == nil {
		stats.Clients = parseNeighClients(string(out))
	}

	if data, err := os.ReadFile("/proc/net/dev"); err == nil {
		if rxVal, txVal, ok := parseInterfaceBytes(string(data), iface); ok {
			stats.RXBytes = rxVal
			stats.TXBytes = txVal
		}
	}

	return stats
}
