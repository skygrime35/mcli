// internal/tui/dockerui/dockerui_test.go
package dockerui

import "testing"

func TestShortID(t *testing.T) {
	if got := shortID("7d35e839c844abcdef0123456789"); got != "7d35e839c844" {
		t.Errorf("shortID() = %q, want %q", got, "7d35e839c844")
	}
	if got := shortID("abc"); got != "abc" {
		t.Errorf("shortID() on a short id should return it unchanged, got %q", got)
	}
}

func TestNewMenuScreen_Title(t *testing.T) {
	// Whether or not docker is actually available on this machine, the
	// screen must always have the same title.
	screen := NewMenuScreen()
	if screen.Title() != "Docker Manager" {
		t.Fatalf("expected title 'Docker Manager', got %q", screen.Title())
	}
}
