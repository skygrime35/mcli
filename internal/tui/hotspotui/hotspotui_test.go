// internal/tui/hotspotui/hotspotui_test.go
package hotspotui

import "testing"

func TestNewMenuScreen_Title(t *testing.T) {
	screen := NewMenuScreen()
	if screen.Title() != "Hotspot Manager" {
		t.Fatalf("expected title 'Hotspot Manager', got %q", screen.Title())
	}
}
