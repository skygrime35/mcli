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
	fmt.Print(health.FormatSnapshot(snap))
	return nil
}
