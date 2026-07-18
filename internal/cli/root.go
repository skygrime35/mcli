// internal/cli/root.go
package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/dockerui"
	"github.com/skygrime35/mcli/internal/tui/hotspotui"
	"github.com/skygrime35/mcli/internal/tui/menu"
	"github.com/skygrime35/mcli/internal/tui/networkui"
	"github.com/skygrime35/mcli/internal/tui/serverui"
	"github.com/skygrime35/mcli/internal/tui/shareui"
	"github.com/skygrime35/mcli/internal/tui/systemui"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "mcli",
		Short: "Personal PC/server management CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInteractive()
		},
	}
	root.AddCommand(newServerCmd())
	root.AddCommand(newHealthCmd())
	root.AddCommand(newDockerCmd())
	root.AddCommand(newSystemCmd())
	root.AddCommand(newHotspotCmd())
	root.AddCommand(newNetworkCmd())
	root.AddCommand(newShareCmd())
	return root
}

func runInteractive() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	main := menu.New("mcli", mainMenuItems(cfg))
	_, err = tui.NewProgram(main).Run()
	return err
}

func mainMenuItems(cfg *config.Config) []menu.Item {
	return []menu.Item{
		{
			Title:       "Server",
			Description: "Wake-on-LAN, SSH connect, SSH key setup",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(serverui.NewListScreen(cfg))
			},
		},
		{Title: "PC Health", Description: "Coming soon", Disabled: true},
		{
			Title:       "Docker",
			Description: "Full purge, clear all containers, select purge",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(dockerui.NewMenuScreen())
			},
		},
		{
			Title:       "System Update",
			Description: "apt update/upgrade/cleanup with safe/warning/unsafe tiers",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(systemui.NewMenuScreen())
			},
		},
		{
			Title:       "Hotspot",
			Description: "Activate/deactivate a Wi-Fi hotspot, view stats",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(hotspotui.NewMenuScreen())
			},
		},
		{
			Title:       "Network Status",
			Description: "IP addresses, ports, connectivity, speed test",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(networkui.NewMenuScreen())
			},
		},
		{
			Title:       "File Sharing",
			Description: "Share a local directory over HTTP",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(shareui.NewMenuScreen())
			},
		},
	}
}

func loadConfig() (*config.Config, error) {
	return config.Load()
}
