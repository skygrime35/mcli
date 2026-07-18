// internal/tui/serverui/detail.go
package serverui

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/server"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
	"github.com/skygrime35/mcli/internal/tui/progress"
)

// NewDetailScreen shows the available actions for a single server.
func NewDetailScreen(cfg *config.Config, s config.ServerConfig) tui.Screen {
	return menu.New(fmt.Sprintf("Server: %s", s.Name), detailItems(cfg, s))
}

// detailItems is split out from NewDetailScreen so it's independently
// testable without needing a running bubbletea program.
func detailItems(cfg *config.Config, s config.ServerConfig) []menu.Item {
	return []menu.Item{
		{
			Title:       "Check status",
			Description: "Test SSH port reachability",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newStatusScreen(s))
			},
		},
		{
			Title:       "Wake on LAN",
			Description: "Send WOL packets and wait for the server to come online",
			OnSelect: func() tea.Cmd {
				ch := server.WakeOnLAN(context.Background(), s)
				return tui.PushScreen(progress.New("Wake on LAN: "+s.Name, ch))
			},
		},
		{
			Title:       "Connect (SSH)",
			Description: "Open an interactive SSH session",
			OnSelect: func() tea.Cmd {
				c := exec.Command("ssh", "-p", strconv.Itoa(s.SSHPort), s.SSHUser+"@"+s.Host)
				return tea.ExecProcess(c, nil)
			},
		},
		{
			Title:       "Setup SSH key",
			Description: "Generate (if needed) and copy an SSH key to this server",
			OnSelect: func() tea.Cmd {
				ch := server.SetupSSHKey(s)
				return tui.PushScreen(progress.New("SSH key setup: "+s.Name, ch))
			},
		},
	}
}
