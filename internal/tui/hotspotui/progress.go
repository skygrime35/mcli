// internal/tui/hotspotui/progress.go
package hotspotui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/hotspot"
	"github.com/skygrime35/mcli/internal/tui"
)

type progressChanMsg struct {
	msg hotspot.ProgressMsg
	ok  bool
}

func waitForNext(ch <-chan hotspot.ProgressMsg) tea.Cmd {
	return func() tea.Msg {
		m, ok := <-ch
		return progressChanMsg{msg: m, ok: ok}
	}
}

type progressScreen struct {
	title  string
	ch     <-chan hotspot.ProgressMsg
	lines  []string
	err    error
	closed bool
}

func newProgressScreen(title string, ch <-chan hotspot.ProgressMsg) tui.Screen {
	return &progressScreen{title: title, ch: ch}
}

func (s *progressScreen) Title() string { return s.title }
func (s *progressScreen) Init() tea.Cmd { return waitForNext(s.ch) }

func (s *progressScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressChanMsg:
		if !msg.ok {
			s.closed = true
			return s, nil
		}
		if msg.msg.Err != nil {
			s.err = msg.msg.Err
		} else {
			s.lines = append(s.lines, msg.msg.Text)
		}
		return s, waitForNext(s.ch)
	case tea.KeyMsg:
		if s.closed && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *progressScreen) View() string {
	out := tui.TitleStyle.Render(s.title) + "\n\n"
	for _, l := range s.lines {
		out += l + "\n"
	}
	if s.err != nil {
		out += tui.ErrorStyle.Render("Error: "+s.err.Error()) + "\n"
	}
	if s.closed {
		out += "\n" + tui.HelpStyle.Render("Done — press enter/esc to go back")
	}
	return out
}

var _ tui.Screen = (*progressScreen)(nil)
