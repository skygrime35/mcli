// internal/tui/hotspotui/stats.go
package hotspotui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/hotspot"
	"github.com/skygrime35/mcli/internal/tui"
)

type statsScreen struct {
	ssid  string
	stats hotspot.Stats
	done  bool
}

type statsLoadedMsg struct{ stats hotspot.Stats }

func newStatsScreen(ssid string) tui.Screen {
	return &statsScreen{ssid: ssid}
}

func (s *statsScreen) Title() string { return "Hotspot Statistics" }

func (s *statsScreen) Init() tea.Cmd {
	ssid := s.ssid
	return func() tea.Msg {
		return statsLoadedMsg{stats: hotspot.GetStats(ssid)}
	}
}

func (s *statsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsLoadedMsg:
		s.stats = msg.stats
		s.done = true
		return s, nil
	case tea.KeyMsg:
		if s.done && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *statsScreen) View() string {
	if !s.done {
		return "Fetching stats..."
	}
	if !s.stats.Active {
		return "Hotspot is not active.\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
	}
	out := fmt.Sprintf("Active interface: %s\n", s.stats.Interface)
	out += fmt.Sprintf("Connected clients (%d):\n", len(s.stats.Clients))
	for _, c := range s.stats.Clients {
		out += "  - " + c + "\n"
	}
	out += fmt.Sprintf("Data usage: TX: %s / RX: %s\n", hotspot.FormatBytes(s.stats.TXBytes), hotspot.FormatBytes(s.stats.RXBytes))
	out += "\n" + tui.HelpStyle.Render("press enter/esc to go back")
	return out
}

var _ tui.Screen = (*statsScreen)(nil)
