// internal/tui/systemui/selection.go
package systemui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/system"
	"github.com/skygrime35/mcli/internal/tui"
)

const (
	toggleUpdate     = "System Update"
	toggleCleanup    = "System Cleanup"
	toggleDoWarnings = "Auto-approve Warnings"
	toggleDoUnsafe   = "Auto-approve Unsafe"
	toggleDryRun     = "Dry Run (Simulate)"
)

// defaultSelectionToggles mirrors the old Python menu's initial checkbox
// state: System Update, System Cleanup, and Auto-approve Warnings are
// enabled by default; Auto-approve Unsafe and Dry Run are not.
func defaultSelectionToggles() []string {
	return []string{toggleUpdate, toggleCleanup, toggleDoWarnings}
}

type selectionScreen struct {
	form     *huh.Form
	selected []string
	done     bool
}

func newSelectionScreen() tui.Screen {
	s := &selectionScreen{selected: defaultSelectionToggles()}
	s.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("System Update — Select options").
				Options(
					huh.NewOption(toggleUpdate, toggleUpdate),
					huh.NewOption(toggleCleanup, toggleCleanup),
					huh.NewOption(toggleDoWarnings, toggleDoWarnings),
					huh.NewOption(toggleDoUnsafe, toggleDoUnsafe),
					huh.NewOption(toggleDryRun, toggleDryRun),
				).
				Value(&s.selected),
		),
	)
	return s
}

func has(selected []string, name string) bool {
	for _, s := range selected {
		if s == name {
			return true
		}
	}
	return false
}

func (s *selectionScreen) Title() string { return "System Update — Selection" }
func (s *selectionScreen) Init() tea.Cmd { return s.form.Init() }

func (s *selectionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" && s.form.State != huh.StateCompleted {
		return s, tui.PopScreen()
	}

	m, cmd := s.form.Update(msg)
	if f, ok := m.(*huh.Form); ok {
		s.form = f
	}

	if s.form.State == huh.StateCompleted && !s.done {
		s.done = true

		analyzeOpts := system.AnalyzeOptions{
			Update: has(s.selected, toggleUpdate),
			Clean:  has(s.selected, toggleCleanup),
		}
		execOpts := system.ExecuteOptions{
			DoWarnings: has(s.selected, toggleDoWarnings),
			DoUnsafe:   has(s.selected, toggleDoUnsafe),
		}
		dryRun := has(s.selected, toggleDryRun)

		return s, tui.PushScreen(newPlanScreen(analyzeOpts, execOpts, dryRun))
	}

	return s, cmd
}

func (s *selectionScreen) View() string { return s.form.View() }

var _ tui.Screen = (*selectionScreen)(nil)

// planScreen runs Analyze (always safe/read-only) and shows the plan,
// then either stops (dry run) or asks for a final confirmation before
// pushing the real executeScreen.
type planScreen struct {
	analyzeOpts system.AnalyzeOptions
	execOpts    system.ExecuteOptions
	dryRun      bool
	plan        system.Plan
	err         error
	done        bool
}

type planAnalyzedMsg struct {
	plan system.Plan
	err  error
}

func newPlanScreen(analyzeOpts system.AnalyzeOptions, execOpts system.ExecuteOptions, dryRun bool) tui.Screen {
	return &planScreen{analyzeOpts: analyzeOpts, execOpts: execOpts, dryRun: dryRun}
}

func (s *planScreen) Title() string { return "System Update — Plan" }

func (s *planScreen) Init() tea.Cmd {
	return func() tea.Msg {
		plan, err := system.Analyze(s.analyzeOpts, platform.Detect())
		return planAnalyzedMsg{plan: plan, err: err}
	}
}

func (s *planScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case planAnalyzedMsg:
		s.plan = msg.plan
		s.err = msg.err
		s.done = true
		// Do NOT auto-push a confirm screen here anymore - the user must
		// see the rendered plan (including any Warning/Unsafe reasons)
		// and explicitly press enter before proceeding. This is the fix
		// for a review finding: silently jumping straight to a bare
		// "Apply N actions?" confirm let a user approve a real update
		// (potentially including kernel removal) without ever reading
		// what was actually flagged and why.
		return s, nil
	case tea.KeyMsg:
		if !s.done {
			return s, nil
		}
		switch msg.String() {
		case "esc":
			return s, tui.PopScreen()
		case "enter":
			if s.err != nil || s.dryRun {
				return s, tui.PopScreen()
			}
			return s, tui.PushScreen(newConfirmScreen(
				"Confirm System Update",
				fmt.Sprintf("Apply %d planned action(s)? (see plan above)", len(s.plan.Actions)),
				func() tea.Cmd {
					return tui.PushScreen(newExecuteScreen(s.plan, s.execOpts, platform.Detect()))
				},
			))
		}
	}
	return s, nil
}

func (s *planScreen) View() string {
	if !s.done {
		return "Analyzing system..."
	}
	out := tui.TitleStyle.Render("Planned actions") + "\n\n"
	if s.err != nil {
		out += tui.ErrorStyle.Render("Error: "+s.err.Error()) + "\n"
		return out + "\n" + tui.HelpStyle.Render("press enter/esc to go back")
	}
	for _, a := range s.plan.Actions {
		out += fmt.Sprintf("  [%s] %s\n", a.Tier, a.Name)
		if a.Reason != "" {
			out += "        " + a.Reason + "\n"
		}
	}
	if s.dryRun {
		out += "\n" + tui.HelpStyle.Render("Dry run — press enter/esc to go back")
	} else {
		out += "\n" + tui.HelpStyle.Render("press enter to review confirmation, esc to go back")
	}
	return out
}

var _ tui.Screen = (*planScreen)(nil)
