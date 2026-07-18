// internal/health/security_test.go
package health

import (
	"testing"

	"github.com/skygrime35/mcli/internal/platform"
)

const sampleAptUpgradable = `Listing...
package-one/jammy 2.0 amd64 [upgradable from: 1.0]
package-two/jammy-security 3.1 amd64 [upgradable from: 3.0]
package-three/jammy 1.5 amd64 [upgradable from: 1.4]
`

func TestParseAptUpgradable(t *testing.T) {
	total, security := parseAptUpgradable(sampleAptUpgradable)
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if security != 1 {
		t.Errorf("security = %d, want 1", security)
	}
}

func TestParseAptUpgradable_NoneUpgradable(t *testing.T) {
	total, security := parseAptUpgradable("Listing...\n")
	if total != 0 || security != 0 {
		t.Errorf("got total=%d security=%d, want 0/0", total, security)
	}
}

// TestFirewallFallback_NoIptables proves the fix for the bug where ufw is
// present but fails at runtime (e.g. permission denied without root): the
// CollectSecurity switch now delegates to firewallFallback in that case,
// exactly as it does for the caps.Iptables and default cases. This test
// covers the deterministic branch of that fallback - no firewall tooling
// available at all - which must report "unknown"/StatusInfo rather than
// leaving FirewallState/FirewallStatus blank.
//
// The caps.Iptables=true branch of firewallFallback shells out to a real
// "iptables" binary, whose availability/output/permissions vary across
// test environments (e.g. requires root), so it is not asserted here to
// keep the test deterministic - matching the precedent set by the
// smartctl exit-code fix, which also documents its exec-dependent path via
// comment rather than asserting it in a test.
func TestFirewallFallback_NoIptables(t *testing.T) {
	state, status := firewallFallback(platform.Capabilities{Iptables: false})
	if state != "unknown" {
		t.Errorf("state = %q, want %q", state, "unknown")
	}
	if status != StatusInfo {
		t.Errorf("status = %q, want %q", status, StatusInfo)
	}
}
