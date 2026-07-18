// internal/cli/network.go
package cli

import (
	"fmt"

	"github.com/skygrime35/mcli/internal/network"
	"github.com/spf13/cobra"
)

func newNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Inspect network interfaces, open ports, and connection speed",
	}
	cmd.AddCommand(newNetworkStatusCmd())
	cmd.AddCommand(newNetworkSpeedtestCmd())
	return cmd
}

func newNetworkStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show IP addresses, listening ports, and internet connectivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			ifaces, err := network.Interfaces()
			if err != nil {
				return err
			}
			fmt.Println("Interfaces:")
			for _, ifc := range ifaces {
				fmt.Printf("  %s: %v\n", ifc.Name, ifc.IPs)
			}

			ports, err := network.ListeningPorts()
			if err != nil {
				return err
			}
			fmt.Printf("\nListening ports (%d):\n", len(ports))
			for _, p := range ports {
				fmt.Printf("  %s/%d\n", p.Protocol, p.Port)
			}

			fmt.Println()
			if network.CheckConnectivity() {
				fmt.Println("Internet connectivity: OK")
			} else {
				fmt.Println("Internet connectivity: UNREACHABLE")
			}
			return nil
		},
	}
}

func newNetworkSpeedtestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "speedtest",
		Short: "Run a download/upload speed test",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Testing...")
			result, err := network.RunSpeedTest()
			if err != nil {
				return err
			}
			fmt.Printf("Server: %s (%s)\n", result.Server, result.Country)
			fmt.Printf("Ping: %.0f ms\n", result.PingMs)
			fmt.Printf("Download: %.2f Mbps\n", result.DownloadMbps)
			fmt.Printf("Upload: %.2f Mbps\n", result.UploadMbps)
			return nil
		},
	}
}
