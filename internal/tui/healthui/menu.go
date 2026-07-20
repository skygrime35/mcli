// internal/tui/healthui/menu.go
package healthui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

// NewMenuScreen is the PC Health feature's entry screen: a one-shot
// snapshot, or a continuously-refreshing watch view.
func NewMenuScreen() tui.Screen {
	items := []menu.Item{
		{
			Title:       "Show Health Summary",
			Description: "One-time CPU/memory/disk/network/services/battery/security snapshot",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newSnapshotScreen())
			},
		},
		{
			Title:       "Watch (live refresh)",
			Description: "Continuously refresh every 5 seconds until you stop it",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newWatchScreen())
			},
		},
	}
	return menu.New("PC Health", items)
}
