package health

import (
	"strings"
	"testing"
	"time"
)

func TestFormatSnapshot_ContainsKeySections(t *testing.T) {
	snap := Snapshot{
		Timestamp: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
		System:    SystemInfo{Hostname: "test-host", OS: "Linux", Kernel: "6.17.0", Architecture: "x86_64", Uptime: "3 days", LoggedUsers: 1},
		CPU:       CPUInfo{Model: "Test CPU", Cores: 8, UsagePercent: 12.5, UsageStatus: StatusGood, Load1: 0.5, Load5: 0.4, Load15: 0.3},
		Memory:    MemoryInfo{TotalKB: 16000000, UsedKB: 8000000, UsagePercent: 50.0, UsageStatus: StatusGood, SwapConfigured: false},
		Network:   NetworkInfo{InternetOK: true, InternetStatus: StatusGood, DNSOK: true, DNSStatus: StatusGood, Gateway: "192.168.1.1"},
		Errors:    ErrorsInfo{RecentErrors: 0, ErrorsStatus: StatusGood, KernelErrors: 0, KernelStatus: StatusGood},
		Security:  SecurityInfo{FirewallState: "active", FirewallStatus: StatusGood},
	}

	out := FormatSnapshot(snap)

	for _, want := range []string{
		"test-host",
		"Test CPU",
		"12.5%",
		"good",
		"192.168.1.1",
		"Not configured", // swap not configured
	} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatSnapshot() output missing %q\nfull output:\n%s", want, out)
		}
	}
}

func TestFormatSnapshot_OmitsAbsentOptionalSections(t *testing.T) {
	snap := Snapshot{
		Timestamp: time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC),
		System:    SystemInfo{Hostname: "test-host"},
		Battery:   nil, // no battery present - section must be omitted entirely
	}

	out := FormatSnapshot(snap)

	if strings.Contains(out, "Battery") {
		t.Errorf("expected no 'Battery' section when snap.Battery is nil, got:\n%s", out)
	}
}
