// internal/tui/shareui/shareui_test.go
package shareui

import "testing"

func TestNewMenuScreen_Title(t *testing.T) {
	screen := NewMenuScreen()
	if screen.Title() != "File Sharing" {
		t.Fatalf("expected title 'File Sharing', got %q", screen.Title())
	}
}
