// internal/system/analyze_test.go
package system

import (
	"testing"

	"github.com/skygrime35/mcli/internal/platform"
)

func TestAnalyze_UpdateOnly_ProducesUpdateActions(t *testing.T) {
	caps := platform.Capabilities{Apt: true}
	plan, err := Analyze(AnalyzeOptions{Update: true, Clean: false}, caps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := actionNames(plan)
	if !contains(names, "Update package lists") || !contains(names, "Upgrade packages") {
		t.Errorf("expected update actions in plan, got %v", names)
	}
	if contains(names, "Clean APT cache") {
		t.Errorf("expected no clean actions when Clean=false, got %v", names)
	}
}

func TestAnalyze_WithoutApt_ReturnsError(t *testing.T) {
	caps := platform.Capabilities{Apt: false}
	if _, err := Analyze(AnalyzeOptions{Update: true, Clean: true}, caps); err == nil {
		t.Fatal("expected an error when apt-get is not available, got nil")
	}
}

func actionNames(plan Plan) []string {
	var names []string
	for _, a := range plan.Actions {
		names = append(names, a.Name)
	}
	return names
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
