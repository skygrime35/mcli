// internal/cli/health.go
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/skygrime35/mcli/internal/health"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/spf13/cobra"
)

func newHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Show PC health (CPU, memory, disks, network, services, battery, errors, security)",
	}
	cmd.AddCommand(newHealthSummaryCmd())
	cmd.AddCommand(newHealthWatchCmd())
	return cmd
}

func newHealthSummaryCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Print a one-time health snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap := health.Collect(platform.Detect())
			return printSnapshot(snap, jsonOutput)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON")
	return cmd
}

func newHealthWatchCmd() *cobra.Command {
	var jsonOutput bool
	var intervalSeconds int
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Continuously refresh the health snapshot until Ctrl+C",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			ch := health.Watch(ctx, platform.Detect(), time.Duration(intervalSeconds)*time.Second)
			for snap := range ch {
				if !jsonOutput {
					fmt.Print("\033[H\033[2J") // clear screen, matching the bash reference's `clear` between refreshes
				}
				if err := printSnapshot(snap, jsonOutput); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output each snapshot as a JSON object (newline-delimited stream)")
	cmd.Flags().IntVar(&intervalSeconds, "interval", 5, "refresh interval in seconds")
	return cmd
}

func printSnapshot(snap health.Snapshot, jsonOutput bool) error {
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		return enc.Encode(snap)
	}

	fmt.Printf("mcli health — %s\n\n", snap.Timestamp.Format(time.RFC1123))

	fmt.Println("System")
	fmt.Printf("  Hostname: %s\n", snap.System.Hostname)
	fmt.Printf("  OS: %s\n", snap.System.OS)
	fmt.Printf("  Kernel: %s (%s)\n", snap.System.Kernel, snap.System.Architecture)
	fmt.Printf("  Uptime: %s\n", snap.System.Uptime)
	fmt.Printf("  Logged users: %d\n\n", snap.System.LoggedUsers)

	fmt.Println("CPU")
	fmt.Printf("  Model: %s (%d cores)\n", snap.CPU.Model, snap.CPU.Cores)
	fmt.Printf("  Usage: %.1f%% [%s]\n", snap.CPU.UsagePercent, snap.CPU.UsageStatus)
	fmt.Printf("  Load: %.2f, %.2f, %.2f\n", snap.CPU.Load1, snap.CPU.Load5, snap.CPU.Load15)
	if snap.CPU.TempAvailable {
		fmt.Printf("  Temp: %.0f°C [%s]\n", snap.CPU.TempCelsius, snap.CPU.TempStatus)
	}
	fmt.Println()

	fmt.Println("Memory")
	fmt.Printf("  Usage: %d/%d KB (%.1f%%) [%s]\n", snap.Memory.UsedKB, snap.Memory.TotalKB, snap.Memory.UsagePercent, snap.Memory.UsageStatus)
	if snap.Memory.SwapConfigured {
		fmt.Printf("  Swap: %d/%d KB (%.1f%%) [%s]\n", snap.Memory.SwapUsedKB, snap.Memory.SwapTotalKB, snap.Memory.SwapPercent, snap.Memory.SwapStatus)
	} else {
		fmt.Println("  Swap: Not configured")
	}
	fmt.Println()

	if len(snap.Disks) > 0 {
		fmt.Println("Disks")
		for _, d := range snap.Disks {
			fmt.Printf("  %s: %s/%s (%d%%) [%s]\n", d.Mount, d.UsedHuman, d.SizeHuman, d.UsagePercent, d.Status)
		}
		fmt.Println()
	}

	if len(snap.SMART) > 0 {
		fmt.Println("S.M.A.R.T.")
		for _, s := range snap.SMART {
			fmt.Printf("  %s: %s [%s]\n", s.Device, s.Detail, s.Status)
		}
		fmt.Println()
	}

	fmt.Println("Network")
	for _, iface := range snap.Network.Interfaces {
		primary := ""
		if iface.Primary {
			primary = " (primary)"
		}
		fmt.Printf("  %s%s: %s [%s] [%s]\n", iface.Name, primary, iface.IP, iface.State, iface.Status)
	}
	fmt.Printf("  Gateway: %s\n", snap.Network.Gateway)
	fmt.Printf("  Internet: %v [%s]\n", snap.Network.InternetOK, snap.Network.InternetStatus)
	fmt.Printf("  DNS: %v [%s]\n\n", snap.Network.DNSOK, snap.Network.DNSStatus)

	if len(snap.Services.Services) > 0 || snap.Services.FailedCount > 0 {
		fmt.Println("Services")
		for _, s := range snap.Services.Services {
			fmt.Printf("  %s: %s [%s]\n", s.Name, s.State, s.Status)
		}
		fmt.Printf("  Failed units: %d [%s]\n\n", snap.Services.FailedCount, snap.Services.FailedStatus)
	}

	if snap.Battery != nil {
		fmt.Println("Battery")
		fmt.Printf("  Charge: %d%% (%s) [%s]\n", snap.Battery.CapacityPercent, snap.Battery.State, snap.Battery.ChargeStatus)
		if snap.Battery.HealthAvailable {
			fmt.Printf("  Health: %.1f%% (%s) [%s]\n", snap.Battery.HealthPercent, snap.Battery.Condition, snap.Battery.HealthStatus)
		}
		if snap.Battery.CycleCount > 0 {
			fmt.Printf("  Cycles: %d [%s]\n", snap.Battery.CycleCount, snap.Battery.CycleStatus)
		}
		if snap.Battery.TimeRemaining != "" {
			fmt.Printf("  Time remaining: %s\n", snap.Battery.TimeRemaining)
		}
		fmt.Println()
	}

	fmt.Println("Errors")
	fmt.Printf("  Recent errors (1h): %d [%s]\n", snap.Errors.RecentErrors, snap.Errors.ErrorsStatus)
	fmt.Printf("  Kernel errors: %d [%s]\n\n", snap.Errors.KernelErrors, snap.Errors.KernelStatus)

	fmt.Println("Security")
	if snap.Security.UpdatesChecked {
		fmt.Printf("  Updates available: %d (security: %d) [%s]\n", snap.Security.UpdatesAvailable, snap.Security.SecurityUpdates, snap.Security.SecurityStatus)
	}
	fmt.Printf("  Firewall: %s [%s]\n", snap.Security.FirewallState, snap.Security.FirewallStatus)

	return nil
}
