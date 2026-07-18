// internal/tui/shareui/running.go
package shareui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/share"
	"github.com/skygrime35/mcli/internal/tui"
)

type runningScreen struct {
	srv     *share.Server
	dir     string
	port    int
	hasPass bool
	url     string
	err     error
	stopped bool
}

type serverStartedMsg struct {
	srv  *share.Server
	dir  string
	port int
	pass bool
	url  string
	err  error
}

func newRunningScreen() tui.Screen {
	return &runningScreen{}
}

func (s *runningScreen) Title() string { return "File Sharing" }

func (s *runningScreen) Init() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil {
			return serverStartedMsg{err: err}
		}
		srv, err := share.Start(cfg.Share.Dir, cfg.Share.Port, cfg.Share.Password)
		if err != nil {
			return serverStartedMsg{err: err}
		}
		url := fmt.Sprintf("http://%s:%d", share.LocalIP(), cfg.Share.Port)
		return serverStartedMsg{srv: srv, dir: cfg.Share.Dir, port: cfg.Share.Port, pass: cfg.Share.Password != "", url: url}
	}
}

func (s *runningScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case serverStartedMsg:
		s.srv = msg.srv
		s.dir = msg.dir
		s.port = msg.port
		s.hasPass = msg.pass
		s.url = msg.url
		s.err = msg.err
		return s, nil
	case tea.KeyMsg:
		if s.err != nil && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
		if s.srv != nil && !s.stopped && (msg.String() == "q" || msg.String() == "esc") {
			s.stopped = true
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.srv.Stop(ctx)
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *runningScreen) View() string {
	if s.err != nil {
		return "Error: " + s.err.Error() + "\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
	}
	if s.srv == nil {
		return "Starting server..."
	}

	out := fmt.Sprintf("Sharing '%s' at: %s\n", s.dir, s.url)
	if s.hasPass {
		out += "Password required (username is ignored).\n"
	} else {
		out += "No password set - open access.\n"
	}
	out += "\n" + tui.HelpStyle.Render("press q to stop sharing and go back")
	return out
}

var _ tui.Screen = (*runningScreen)(nil)
