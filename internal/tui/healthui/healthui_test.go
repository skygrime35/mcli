// internal/tui/healthui/healthui_test.go
package healthui

import "testing"

func TestNewMenuScreen_Title(t *testing.T) {
	screen := NewMenuScreen()
	if screen.Title() != "PC Health" {
		t.Fatalf("expected title 'PC Health', got %q", screen.Title())
	}
}
