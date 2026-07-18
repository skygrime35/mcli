// internal/system/execute_test.go
package system

import "testing"

func TestShouldRun_Safe_AlwaysRuns(t *testing.T) {
	run, _ := ShouldRun(Action{Tier: TierSafe}, ExecuteOptions{})
	if !run {
		t.Error("expected Safe actions to always run")
	}
}

func TestShouldRun_Warning_RequiresDoWarnings(t *testing.T) {
	action := Action{Tier: TierWarning, Reason: "some reason"}

	if run, _ := ShouldRun(action, ExecuteOptions{}); run {
		t.Error("expected Warning action NOT to run without DoWarnings")
	}
	if run, _ := ShouldRun(action, ExecuteOptions{DoWarnings: true}); !run {
		t.Error("expected Warning action to run with DoWarnings=true")
	}
	if run, _ := ShouldRun(action, ExecuteOptions{SkipWarnings: true, DoWarnings: true}); run {
		t.Error("expected SkipWarnings to take precedence and skip even if DoWarnings is also true")
	}
}

func TestShouldRun_Unsafe_RequiresDoUnsafe(t *testing.T) {
	action := Action{Tier: TierUnsafe, Reason: "CRITICAL: current kernel detected in removal list!"}

	if run, _ := ShouldRun(action, ExecuteOptions{}); run {
		t.Error("expected Unsafe action NOT to run without DoUnsafe")
	}
	if run, _ := ShouldRun(action, ExecuteOptions{DoUnsafe: true}); !run {
		t.Error("expected Unsafe action to run with DoUnsafe=true")
	}
	if run, _ := ShouldRun(action, ExecuteOptions{SkipUnsafe: true, DoUnsafe: true}); run {
		t.Error("expected SkipUnsafe to take precedence and skip even if DoUnsafe is also true")
	}
}

func TestShouldRun_NotesExplainSkipReason(t *testing.T) {
	_, note := ShouldRun(Action{Tier: TierUnsafe, Reason: "CRITICAL: x"}, ExecuteOptions{})
	if note == "" {
		t.Error("expected a non-empty note explaining why an unapproved Unsafe action was skipped")
	}
}
