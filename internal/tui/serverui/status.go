// internal/tui/serverui/status.go
package serverui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/server"
	"github.com/skygrime35/mcli/internal/tui"
)

type statusCheckedMsg struct{ online bool }

type statusScreen struct {
	server  config.ServerConfig
	checked bool
	online  bool
}

func newStatusScreen(s config.ServerConfig) *statusScreen {
	return &statusScreen{server: s}
}

func (s *statusScreen) Title() string { return "Status: " + s.server.Name }

func (s *statusScreen) Init() tea.Cmd {
	srv := s.server
	return func() tea.Msg {
		status := server.CheckStatus(context.Background(), srv.Host, srv.SSHPort, 5*time.Second)
		return statusCheckedMsg{online: status.Online}
	}
}

func (s *statusScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusCheckedMsg:
		s.checked = true
		s.online = msg.online
		return s, nil
	case tea.KeyMsg:
		if s.checked && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *statusScreen) View() string {
	if !s.checked {
		return "Checking " + s.server.Name + "..."
	}
	if s.online {
		return tui.SuccessStyle.Render(s.server.Name+" is online.") + "\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
	}
	return tui.ErrorStyle.Render(s.server.Name+" is offline.") + "\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
}

var _ tui.Screen = (*statusScreen)(nil)
