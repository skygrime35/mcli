// internal/tui/healthui/snapshot.go
package healthui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/health"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/tui"
)

type snapshotScreen struct {
	text string
	done bool
}

type snapshotLoadedMsg struct{ text string }

func newSnapshotScreen() tui.Screen {
	return &snapshotScreen{}
}

func (s *snapshotScreen) Title() string { return "PC Health" }

func (s *snapshotScreen) Init() tea.Cmd {
	return func() tea.Msg {
		snap := health.Collect(platform.Detect())
		return snapshotLoadedMsg{text: health.FormatSnapshot(snap)}
	}
}

func (s *snapshotScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case snapshotLoadedMsg:
		s.text = msg.text
		s.done = true
		return s, nil
	case tea.KeyMsg:
		if s.done && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *snapshotScreen) View() string {
	if !s.done {
		return "Collecting health snapshot..."
	}
	return s.text + "\n" + tui.HelpStyle.Render("press enter/esc to go back")
}

var _ tui.Screen = (*snapshotScreen)(nil)
