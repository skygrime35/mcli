// internal/tui/networkui/menu.go
package networkui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

// NewMenuScreen is the Network Status feature's entry screen: view
// interfaces/ports/connectivity, or run a speed test.
func NewMenuScreen() tui.Screen {
	items := []menu.Item{
		{
			Title:       "Show Network Status",
			Description: "IP addresses, listening ports, internet connectivity",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newStatusScreen())
			},
		},
		{
			Title:       "Run Speed Test",
			Description: "Measure download/upload speed (takes ~10-30s)",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newSpeedtestScreen())
			},
		},
	}
	return menu.New("Network Status", items)
}
