// internal/health/format.go
package health

import (
	"fmt"
	"strings"
	"time"
)

// FormatSnapshot renders a Snapshot as human-readable, section-by-section
// text - the same layout `mcli health summary`/`watch` have always
// printed, now shared with the TUI so both surfaces render identically.
func FormatSnapshot(snap Snapshot) string {
	var b strings.Builder

	fmt.Fprintf(&b, "mcli health — %s\n\n", snap.Timestamp.Format(time.RFC1123))

	fmt.Fprintln(&b, "System")
	fmt.Fprintf(&b, "  Hostname: %s\n", snap.System.Hostname)
	fmt.Fprintf(&b, "  OS: %s\n", snap.System.OS)
	fmt.Fprintf(&b, "  Kernel: %s (%s)\n", snap.System.Kernel, snap.System.Architecture)
	fmt.Fprintf(&b, "  Uptime: %s\n", snap.System.Uptime)
	fmt.Fprintf(&b, "  Logged users: %d\n\n", snap.System.LoggedUsers)

	fmt.Fprintln(&b, "CPU")
	fmt.Fprintf(&b, "  Model: %s (%d cores)\n", snap.CPU.Model, snap.CPU.Cores)
	fmt.Fprintf(&b, "  Usage: %.1f%% [%s]\n", snap.CPU.UsagePercent, snap.CPU.UsageStatus)
	fmt.Fprintf(&b, "  Load: %.2f, %.2f, %.2f\n", snap.CPU.Load1, snap.CPU.Load5, snap.CPU.Load15)
	if snap.CPU.TempAvailable {
		fmt.Fprintf(&b, "  Temp: %.0f°C [%s]\n", snap.CPU.TempCelsius, snap.CPU.TempStatus)
	}
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "Memory")
	fmt.Fprintf(&b, "  Usage: %d/%d KB (%.1f%%) [%s]\n", snap.Memory.UsedKB, snap.Memory.TotalKB, snap.Memory.UsagePercent, snap.Memory.UsageStatus)
	if snap.Memory.SwapConfigured {
		fmt.Fprintf(&b, "  Swap: %d/%d KB (%.1f%%) [%s]\n", snap.Memory.SwapUsedKB, snap.Memory.SwapTotalKB, snap.Memory.SwapPercent, snap.Memory.SwapStatus)
	} else {
		fmt.Fprintln(&b, "  Swap: Not configured")
	}
	fmt.Fprintln(&b)

	if len(snap.Disks) > 0 {
		fmt.Fprintln(&b, "Disks")
		for _, d := range snap.Disks {
			fmt.Fprintf(&b, "  %s: %s/%s (%d%%) [%s]\n", d.Mount, d.UsedHuman, d.SizeHuman, d.UsagePercent, d.Status)
		}
		fmt.Fprintln(&b)
	}

	if len(snap.SMART) > 0 {
		fmt.Fprintln(&b, "S.M.A.R.T.")
		for _, s := range snap.SMART {
			fmt.Fprintf(&b, "  %s: %s [%s]\n", s.Device, s.Detail, s.Status)
		}
		fmt.Fprintln(&b)
	}

	fmt.Fprintln(&b, "Network")
	for _, iface := range snap.Network.Interfaces {
		primary := ""
		if iface.Primary {
			primary = " (primary)"
		}
		fmt.Fprintf(&b, "  %s%s: %s [%s] [%s]\n", iface.Name, primary, iface.IP, iface.State, iface.Status)
	}
	fmt.Fprintf(&b, "  Gateway: %s\n", snap.Network.Gateway)
	fmt.Fprintf(&b, "  Internet: %v [%s]\n", snap.Network.InternetOK, snap.Network.InternetStatus)
	fmt.Fprintf(&b, "  DNS: %v [%s]\n\n", snap.Network.DNSOK, snap.Network.DNSStatus)

	if len(snap.Services.Services) > 0 || snap.Services.FailedCount > 0 {
		fmt.Fprintln(&b, "Services")
		for _, s := range snap.Services.Services {
			fmt.Fprintf(&b, "  %s: %s [%s]\n", s.Name, s.State, s.Status)
		}
		fmt.Fprintf(&b, "  Failed units: %d [%s]\n\n", snap.Services.FailedCount, snap.Services.FailedStatus)
	}

	if snap.Battery != nil {
		fmt.Fprintln(&b, "Battery")
		fmt.Fprintf(&b, "  Charge: %d%% (%s) [%s]\n", snap.Battery.CapacityPercent, snap.Battery.State, snap.Battery.ChargeStatus)
		if snap.Battery.HealthAvailable {
			fmt.Fprintf(&b, "  Health: %.1f%% (%s) [%s]\n", snap.Battery.HealthPercent, snap.Battery.Condition, snap.Battery.HealthStatus)
		}
		if snap.Battery.CycleCount > 0 {
			fmt.Fprintf(&b, "  Cycles: %d [%s]\n", snap.Battery.CycleCount, snap.Battery.CycleStatus)
		}
		if snap.Battery.TimeRemaining != "" {
			fmt.Fprintf(&b, "  Time remaining: %s\n", snap.Battery.TimeRemaining)
		}
		fmt.Fprintln(&b)
	}

	fmt.Fprintln(&b, "Errors")
	fmt.Fprintf(&b, "  Recent errors (1h): %d [%s]\n", snap.Errors.RecentErrors, snap.Errors.ErrorsStatus)
	fmt.Fprintf(&b, "  Kernel errors: %d [%s]\n\n", snap.Errors.KernelErrors, snap.Errors.KernelStatus)

	fmt.Fprintln(&b, "Security")
	if snap.Security.UpdatesChecked {
		fmt.Fprintf(&b, "  Updates available: %d (security: %d) [%s]\n", snap.Security.UpdatesAvailable, snap.Security.SecurityUpdates, snap.Security.SecurityStatus)
	}
	fmt.Fprintf(&b, "  Firewall: %s [%s]\n", snap.Security.FirewallState, snap.Security.FirewallStatus)

	return b.String()
}
