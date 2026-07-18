// internal/tui/shareui/menu.go
package shareui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

// NewMenuScreen is the File Sharing feature's entry screen.
func NewMenuScreen() tui.Screen {
	items := []menu.Item{
		{
			Title:       "Start Sharing",
			Description: "Start the HTTP file server using your configured directory",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newRunningScreen())
			},
		},
	}
	return menu.New("File Sharing", items)
}
