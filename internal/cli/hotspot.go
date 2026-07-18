// internal/cli/hotspot.go
package cli

import (
	"fmt"

	"github.com/skygrime35/mcli/internal/hotspot"
	"github.com/spf13/cobra"
)

func newHotspotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hotspot",
		Short: "Activate, deactivate, or inspect a Wi-Fi hotspot (nmcli)",
	}
	cmd.AddCommand(newHotspotOnCmd())
	cmd.AddCommand(newHotspotOffCmd())
	cmd.AddCommand(newHotspotStatsCmd())
	return cmd
}

func hotspotCredentials() (ssid, password string, err error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", "", err
	}
	ssid = cfg.Hotspot.SSID
	if ssid == "" {
		ssid = "Hotspot"
	}
	password = cfg.Hotspot.Password
	if password == "" {
		password = "123456789"
	}
	return ssid, password, nil
}

func runHotspotProgress(ch <-chan hotspot.ProgressMsg) error {
	var firstErr error
	for msg := range ch {
		if msg.Err != nil {
			fmt.Println("Error:", msg.Err)
			if firstErr == nil {
				firstErr = msg.Err
			}
			continue
		}
		fmt.Println(msg.Text)
	}
	return firstErr
}

func newHotspotOnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "on",
		Short: "Activate the configured hotspot",
		RunE: func(cmd *cobra.Command, args []string) error {
			ssid, password, err := hotspotCredentials()
			if err != nil {
				return err
			}
			return runHotspotProgress(hotspot.Activate(ssid, password))
		},
	}
}

func newHotspotOffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "off",
		Short: "Deactivate the configured hotspot",
		RunE: func(cmd *cobra.Command, args []string) error {
			ssid, _, err := hotspotCredentials()
			if err != nil {
				return err
			}
			return runHotspotProgress(hotspot.Deactivate(ssid))
		},
	}
}

func newHotspotStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show hotspot connection statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			ssid, _, err := hotspotCredentials()
			if err != nil {
				return err
			}
			stats := hotspot.GetStats(ssid)
			if !stats.Active {
				fmt.Println("Hotspot is not active.")
				return nil
			}
			fmt.Printf("Active interface: %s\n", stats.Interface)
			fmt.Printf("Connected clients (%d):\n", len(stats.Clients))
			for _, c := range stats.Clients {
				fmt.Printf("  - %s\n", c)
			}
			fmt.Printf("Data usage: TX: %s / RX: %s\n", hotspot.FormatBytes(stats.TXBytes), hotspot.FormatBytes(stats.RXBytes))
			return nil
		},
	}
}
