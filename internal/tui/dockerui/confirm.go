// internal/tui/dockerui/confirm.go
package dockerui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/skygrime35/mcli/internal/tui"
)

type confirmScreen struct {
	title string
	form  *huh.Form
	value bool
	onYes func() tea.Cmd
	done  bool
}

func newConfirmScreen(title, question string, onYes func() tea.Cmd) tui.Screen {
	s := &confirmScreen{title: title, onYes: onYes}
	s.form = huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(question).Value(&s.value)))
	return s
}

func (s *confirmScreen) Title() string { return s.title }
func (s *confirmScreen) Init() tea.Cmd { return s.form.Init() }

func (s *confirmScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" && s.form.State != huh.StateCompleted {
		return s, tui.PopScreen()
	}

	m, cmd := s.form.Update(msg)
	if f, ok := m.(*huh.Form); ok {
		s.form = f
	}

	if s.form.State == huh.StateCompleted && !s.done {
		s.done = true
		if s.value {
			return s, s.onYes()
		}
		return s, tui.PopScreen()
	}

	return s, cmd
}

func (s *confirmScreen) View() string { return s.form.View() }

var _ tui.Screen = (*confirmScreen)(nil)
