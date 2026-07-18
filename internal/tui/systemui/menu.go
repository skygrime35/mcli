// internal/tui/systemui/menu.go
package systemui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/system"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

// NewMenuScreen is System Update's entry screen: Full Update or
// Selection Update, matching the old app.py's two-option menu.
func NewMenuScreen() tui.Screen {
	if !platform.Detect().Apt {
		return menu.New("System Update", []menu.Item{
			{Title: "apt-get not found", Disabled: true},
		})
	}

	items := []menu.Item{
		{
			Title:       "Full Update",
			Description: "Update, upgrade, clean, and approve ALL warnings/unsafe actions",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newConfirmScreen(
					"Full Update",
					"Start Full System Update (includes unsafe operations)?",
					func() tea.Cmd {
						return tui.PushScreen(newPlanScreen(
							system.AnalyzeOptions{Update: true, Clean: true},
							system.ExecuteOptions{DoWarnings: true, DoUnsafe: true},
							false,
						))
					},
				))
			},
		},
		{
			Title:       "Selection Update",
			Description: "Choose which categories and approval levels to apply",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newSelectionScreen())
			},
		},
	}
	return menu.New("System Update", items)
}
