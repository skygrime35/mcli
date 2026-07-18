// internal/tui/dockerui/selectpurge.go
package dockerui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/skygrime35/mcli/internal/docker"
	"github.com/skygrime35/mcli/internal/tui"
)

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

type selectPurgeScreen struct {
	form     *huh.Form
	selected []string
	done     bool
}

func newSelectPurgeScreen(containers []docker.Container) tui.Screen {
	s := &selectPurgeScreen{}
	options := make([]huh.Option[string], len(containers))
	for i, c := range containers {
		label := fmt.Sprintf("%s (%s) - %s", c.Name, shortID(c.ID), c.Status)
		options[i] = huh.NewOption(label, c.ID)
	}
	s.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select containers to purge").
				Description("space to toggle, enter to confirm").
				Options(options...).
				Filterable(true).
				Value(&s.selected),
		),
	)
	return s
}

func (s *selectPurgeScreen) Title() string { return "Select Purge" }
func (s *selectPurgeScreen) Init() tea.Cmd { return s.form.Init() }

func (s *selectPurgeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" && s.form.State != huh.StateCompleted {
		return s, tui.PopScreen()
	}

	m, cmd := s.form.Update(msg)
	if f, ok := m.(*huh.Form); ok {
		s.form = f
	}

	if s.form.State == huh.StateCompleted && !s.done {
		s.done = true
		if len(s.selected) == 0 {
			return s, tui.PopScreen()
		}
		selected := s.selected
		return s, tui.PushScreen(newConfirmScreen(
			"Confirm purge",
			fmt.Sprintf("Purge %d selected container(s)? This cannot be undone.", len(selected)),
			func() tea.Cmd {
				return tui.PushScreen(newProgressScreen("Select Purge", docker.SelectPurge(selected)))
			},
		))
	}

	return s, cmd
}

func (s *selectPurgeScreen) View() string { return s.form.View() }

var _ tui.Screen = (*selectPurgeScreen)(nil)
