// internal/tui/dockerui/loading.go
package dockerui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/docker"
	"github.com/skygrime35/mcli/internal/tui"
)

type containersLoadedMsg struct {
	containers []docker.Container
	err        error
}

type loadingContainersScreen struct {
	done    bool
	message string
	isError bool
}

func newLoadingContainersScreen() tui.Screen {
	return &loadingContainersScreen{}
}

func (s *loadingContainersScreen) Title() string { return "Loading containers..." }

func (s *loadingContainersScreen) Init() tea.Cmd {
	return func() tea.Msg {
		containers, err := docker.ListContainers()
		return containersLoadedMsg{containers: containers, err: err}
	}
}

func (s *loadingContainersScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case containersLoadedMsg:
		if msg.err != nil {
			s.done, s.isError, s.message = true, true, msg.err.Error()
			return s, nil
		}
		if len(msg.containers) == 0 {
			s.done, s.isError, s.message = true, false, "No containers found."
			return s, nil
		}
		return s, tui.PushScreen(newSelectPurgeScreen(msg.containers))
	case tea.KeyMsg:
		if s.done && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *loadingContainersScreen) View() string {
	if !s.done {
		return "Fetching containers..."
	}
	line := s.message
	if s.isError {
		line = tui.ErrorStyle.Render("Error: " + s.message)
	}
	return line + "\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
}

var _ tui.Screen = (*loadingContainersScreen)(nil)
