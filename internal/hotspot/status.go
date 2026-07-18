// internal/hotspot/status.go
package hotspot

import (
	"os/exec"
	"strings"
)

// parseActiveConnectionNames parses `nmcli -t -f NAME connection show
// --active` output (one connection name per line) into a slice.
func parseActiveConnectionNames(output string) []string {
	var names []string
	for _, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
		if line != "" {
			names = append(names, line)
		}
	}
	return names
}

func containsName(names []string, target string) bool {
	for _, n := range names {
		if n == target {
			return true
		}
	}
	return false
}

// IsActive reports whether a connection named ssid is currently active.
// Uses machine-readable terse output and an exact-match comparison in Go
// - unlike the old Python reference's `grep -qw` through a shell string,
// this has no shell-injection surface at all.
func IsActive(ssid string) bool {
	out, err := exec.Command("nmcli", "-t", "-f", "NAME", "connection", "show", "--active").Output()
	if err != nil {
		return false
	}
	return containsName(parseActiveConnectionNames(string(out)), ssid)
}
