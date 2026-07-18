// internal/tui/dockerui/menu.go
package dockerui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/docker"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

// NewMenuScreen is the Docker Manager's entry screen: Select Purge, Full
// Purge, Clear All. If docker isn't available, a single disabled item
// explains why.
func NewMenuScreen() tui.Screen {
	if !docker.IsAvailable() {
		return menu.New("Docker Manager", []menu.Item{
			{Title: "Docker command not found", Disabled: true},
		})
	}

	items := []menu.Item{
		{
			Title:       "Select Purge",
			Description: "Choose specific containers to stop and remove",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newLoadingContainersScreen())
			},
		},
		{
			Title:       "Full Purge",
			Description: "Remove ALL containers, images, volumes, networks, and build cache",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newConfirmScreen(
					"Full Purge",
					"WARNING: this will delete EVERYTHING. Are you sure?",
					func() tea.Cmd {
						return tui.PushScreen(newProgressScreen("Full Purge", docker.FullPurge()))
					},
				))
			},
		},
		{
			Title:       "Clear All Containers",
			Description: "Stop and remove all containers (images/volumes untouched)",
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(newConfirmScreen(
					"Clear All Containers",
					"WARNING: this will stop and remove ALL containers. Are you sure?",
					func() tea.Cmd {
						return tui.PushScreen(newProgressScreen("Clear All Containers", docker.ClearAll()))
					},
				))
			},
		},
	}
	return menu.New("Docker Manager", items)
}
