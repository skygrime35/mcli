// internal/tui/systemui/execute.go
package systemui

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/system"
	"github.com/skygrime35/mcli/internal/tui"
)

type actionDoneMsg struct{ err error }

// executeScreen sequences through plan's approved actions one at a time
// using tea.ExecProcess, so each privileged command gets a real terminal
// (sudo password prompts, live apt output) - the same mechanism already
// used for "Connect (SSH)" in the Server Manager TUI screens.
type executeScreen struct {
	plan  system.Plan
	opts  system.ExecuteOptions
	caps  platform.Capabilities
	index int
	lines []string
	err   error
	done  bool
}

func newExecuteScreen(plan system.Plan, opts system.ExecuteOptions, caps platform.Capabilities) tui.Screen {
	return &executeScreen{plan: plan, opts: opts, caps: caps}
}

func (s *executeScreen) Title() string { return "System Update" }
func (s *executeScreen) Init() tea.Cmd { return s.runNext() }

func (s *executeScreen) runNext() tea.Cmd {
	if !s.caps.Apt {
		s.err = fmt.Errorf("apt-get not found - System Update is only supported on Debian/Ubuntu-based systems")
		s.done = true
		return nil
	}
	for s.index < len(s.plan.Actions) {
		action := s.plan.Actions[s.index]
		run, note := system.ShouldRun(action, s.opts)
		s.lines = append(s.lines, fmt.Sprintf("-> %s%s", action.Name, note))
		s.index++
		if !run {
			continue
		}
		args := append([]string{action.Command[0]}, action.Command[1:]...)
		cmd := exec.Command("sudo", args...)
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			return actionDoneMsg{err: err}
		})
	}
	s.done = true
	return nil
}

func (s *executeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case actionDoneMsg:
		if msg.err != nil {
			s.err = msg.err
			s.lines = append(s.lines, "Error: "+msg.err.Error())
			s.done = true
			return s, nil
		}
		s.lines = append(s.lines, "   done.")
		return s, s.runNext()
	case tea.KeyMsg:
		if s.done && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *executeScreen) View() string {
	out := tui.TitleStyle.Render("System Update") + "\n\n"
	for _, l := range s.lines {
		out += l + "\n"
	}
	if s.err != nil {
		out += "\n" + tui.ErrorStyle.Render("Error: "+s.err.Error()) + "\n"
	}
	if s.done {
		out += "\n" + tui.HelpStyle.Render("Done — press enter/esc to go back")
	}
	return out
}

var _ tui.Screen = (*executeScreen)(nil)
