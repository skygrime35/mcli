// internal/tui/networkui/networkui_test.go
package networkui

import "testing"

func TestNewMenuScreen_Title(t *testing.T) {
	screen := NewMenuScreen()
	if screen.Title() != "Network Status" {
		t.Fatalf("expected title 'Network Status', got %q", screen.Title())
	}
}
