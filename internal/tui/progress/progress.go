// internal/tui/progress/progress.go
package progress

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/server"
	"github.com/skygrime35/mcli/internal/tui"
)

type progressChanMsg struct {
	msg server.ProgressMsg
	ok  bool
}

func waitForNext(ch <-chan server.ProgressMsg) tea.Cmd {
	return func() tea.Msg {
		m, ok := <-ch
		return progressChanMsg{msg: m, ok: ok}
	}
}

// Screen streams messages from a <-chan server.ProgressMsg (as produced by
// server.WakeOnLAN / server.SetupSSHKey) into the terminal live, and lets
// the user return to the previous screen once the channel closes.
type Screen struct {
	title  string
	ch     <-chan server.ProgressMsg
	lines  []string
	err    error
	closed bool
}

func New(title string, ch <-chan server.ProgressMsg) *Screen {
	return &Screen{title: title, ch: ch}
}

func (s *Screen) Title() string { return s.title }
func (s *Screen) Init() tea.Cmd { return waitForNext(s.ch) }

func (s *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s *Screen) View() string {
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

var _ tui.Screen = (*Screen)(nil)
