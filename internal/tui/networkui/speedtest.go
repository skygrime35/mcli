// internal/tui/networkui/speedtest.go
package networkui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/network"
	"github.com/skygrime35/mcli/internal/tui"
)

type speedtestScreen struct {
	result network.SpeedTestResult
	err    error
	done   bool
}

type speedtestDoneMsg struct {
	result network.SpeedTestResult
	err    error
}

func newSpeedtestScreen() tui.Screen {
	return &speedtestScreen{}
}

func (s *speedtestScreen) Title() string { return "Network Status" }

func (s *speedtestScreen) Init() tea.Cmd {
	return func() tea.Msg {
		result, err := network.RunSpeedTest()
		return speedtestDoneMsg{result: result, err: err}
	}
}

func (s *speedtestScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case speedtestDoneMsg:
		s.result = msg.result
		s.err = msg.err
		s.done = true
		return s, nil
	case tea.KeyMsg:
		if s.done && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *speedtestScreen) View() string {
	if !s.done {
		return "Testing... (this takes 10-30 seconds)"
	}
	if s.err != nil {
		return "Error: " + s.err.Error() + "\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
	}

	out := fmt.Sprintf("Server: %s (%s)\n", s.result.Server, s.result.Country)
	out += fmt.Sprintf("Ping: %.0f ms\n", s.result.PingMs)
	out += fmt.Sprintf("Download: %.2f Mbps\n", s.result.DownloadMbps)
	out += fmt.Sprintf("Upload: %.2f Mbps\n", s.result.UploadMbps)
	out += "\n" + tui.HelpStyle.Render("press enter/esc to go back")
	return out
}

var _ tui.Screen = (*speedtestScreen)(nil)
