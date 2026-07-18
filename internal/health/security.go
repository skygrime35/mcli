// internal/health/security.go
package health

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/skygrime35/mcli/internal/platform"
)

func parseAptUpgradable(output string) (total, security int) {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Listing") {
			continue
		}
		total++
		if strings.Contains(strings.ToLower(line), "security") {
			security++
		}
	}
	return total, security
}

// firewallFallback reports firewall state when ufw is unavailable, or when
// ufw is present but failed to run (e.g. requires root). It tries iptables
// if available, otherwise reports "unknown".
func firewallFallback(caps platform.Capabilities) (string, Status) {
	if caps.Iptables {
		out, err := exec.Command("iptables", "-L").Output()
		if err == nil {
			count := parseCount(string(out))
			return "rules:" + strconv.Itoa(count), StatusInfo
		}
	}
	return "unknown", StatusInfo
}

func CollectSecurity(caps platform.Capabilities) SecurityInfo {
	var info SecurityInfo

	if caps.Apt {
		out, err := exec.Command("apt", "list", "--upgradable").Output()
		if err == nil {
			info.UpdatesChecked = true
			info.UpdatesAvailable, info.SecurityUpdates = parseAptUpgradable(string(out))
			info.SecurityStatus = statusForHigh(float64(info.SecurityUpdates), securityUpdatesWarning, securityUpdatesCritical)
		}
	}

	switch {
	case caps.Ufw:
		out, err := exec.Command("ufw", "status").Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 0 && strings.Contains(strings.ToLower(lines[0]), "active") {
				info.FirewallState = "active"
				info.FirewallStatus = StatusGood
				break
			}
			info.FirewallState = "inactive"
			info.FirewallStatus = StatusWarning
			break
		}
		// ufw is installed but failed to run (e.g. requires root) - fall
		// back to iptables/unknown detection instead of leaving fields blank.
		info.FirewallState, info.FirewallStatus = firewallFallback(caps)
	case caps.Iptables:
		info.FirewallState, info.FirewallStatus = firewallFallback(caps)
	default:
		info.FirewallState, info.FirewallStatus = firewallFallback(caps)
	}

	if caps.Journalctl {
		out, err := exec.Command("journalctl", "-u", "sshd", "--since", "24 hours ago").Output()
		if err == nil {
			count := 0
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(strings.ToLower(line), "failed") {
					count++
				}
			}
			info.FailedSSHAttempts = count
		}
	}

	return info
}
