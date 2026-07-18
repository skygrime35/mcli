// internal/tui/systemui/systemui_test.go
package systemui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/system"
	"github.com/skygrime35/mcli/internal/tui"
)

func TestNewMenuScreen_Title(t *testing.T) {
	screen := NewMenuScreen()
	if screen.Title() != "System Update" {
		t.Fatalf("expected title 'System Update', got %q", screen.Title())
	}
}

// TestPlanScreen_DoesNotAutoPushConfirm_RequiresEnter guards against a
// review finding: planScreen used to push the confirmScreen automatically
// as soon as analysis completed (planAnalyzedMsg), meaning the user never
// actually saw the rendered plan - including any Warning/Unsafe reasons
// like "CRITICAL: current kernel detected in removal list!" - before being
// asked to confirm a real system update. Now planAnalyzedMsg must only
// mark the screen as done (no cmd), and only an explicit "enter" key press
// afterwards may push the confirm screen.
func TestPlanScreen_DoesNotAutoPushConfirm_RequiresEnter(t *testing.T) {
	screen := newPlanScreen(system.AnalyzeOptions{Update: true}, system.ExecuteOptions{}, false)
	s, ok := screen.(*planScreen)
	if !ok {
		t.Fatalf("expected *planScreen, got %T", screen)
	}

	plan := system.Plan{Actions: []system.Action{
		{Name: "Remove old kernels", Tier: system.TierUnsafe, Reason: "CRITICAL: current kernel detected in removal list!"},
	}}

	_, cmd := s.Update(planAnalyzedMsg{plan: plan, err: nil})
	if cmd != nil {
		t.Fatalf("expected planAnalyzedMsg to NOT return a cmd (no auto-push to confirm), got a non-nil cmd producing %T", cmd())
	}
	if !s.done {
		t.Fatal("expected s.done to be true after planAnalyzedMsg")
	}

	_, cmd = s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected pressing enter after the plan is shown to return a non-nil cmd")
	}
	if _, ok := cmd().(tui.PushScreenMsg); !ok {
		t.Fatalf("expected enter to push the confirm screen (PushScreenMsg), got %T", cmd())
	}
}

// TestPlanScreen_DryRun_EnterPopsScreen_NeverConfirms confirms that a dry
// run never reaches the confirm/execute path, even via the same "enter"
// key that triggers confirmation in the real (non-dry-run) case above.
func TestPlanScreen_DryRun_EnterPopsScreen_NeverConfirms(t *testing.T) {
	screen := newPlanScreen(system.AnalyzeOptions{Update: true}, system.ExecuteOptions{}, true)
	s, ok := screen.(*planScreen)
	if !ok {
		t.Fatalf("expected *planScreen, got %T", screen)
	}

	plan := system.Plan{Actions: []system.Action{
		{Name: "Remove old kernels", Tier: system.TierUnsafe, Reason: "CRITICAL: current kernel detected in removal list!"},
	}}

	_, cmd := s.Update(planAnalyzedMsg{plan: plan, err: nil})
	if cmd != nil {
		t.Fatalf("expected planAnalyzedMsg to NOT return a cmd, got a non-nil cmd producing %T", cmd())
	}

	_, cmd = s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected pressing enter in dry-run mode to return a non-nil cmd")
	}
	if _, ok := cmd().(tui.PopScreenMsg); !ok {
		t.Fatalf("expected dry-run enter to pop the screen (PopScreenMsg), never push a confirm/execute screen, got %T", cmd())
	}
}

func TestDefaultSelectionToggles(t *testing.T) {
	got := defaultSelectionToggles()
	want := map[string]bool{
		"System Update":         true,
		"System Cleanup":        true,
		"Auto-approve Warnings": true,
		"Auto-approve Unsafe":   false,
		"Dry Run (Simulate)":    false,
	}
	wantCount := 0
	for _, enabled := range want {
		if enabled {
			wantCount++
		}
	}
	if len(got) != wantCount {
		t.Fatalf("expected %d default toggles selected, got %d: %v", wantCount, len(got), got)
	}
	for _, name := range got {
		if !want[name] {
			t.Errorf("toggle %q should not be selected by default", name)
		}
	}
}
