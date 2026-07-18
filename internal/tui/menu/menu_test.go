package menu

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/tui"
)

func TestScreen_EnterOnEnabledItem_InvokesOnSelect(t *testing.T) {
	called := false
	items := []Item{
		{Title: "one", Description: "first", OnSelect: func() tea.Cmd {
			called = true
			return tui.PopScreen()
		}},
	}
	s := New("Test menu", items)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !called {
		t.Fatal("expected OnSelect to be called")
	}
	if cmd == nil {
		t.Fatal("expected a non-nil cmd from Update")
	}
	if _, ok := cmd().(tui.PopScreenMsg); !ok {
		t.Fatalf("expected the cmd to produce a PopScreenMsg, got %T", cmd())
	}
}

func TestScreen_EnterOnDisabledItem_DoesNothing(t *testing.T) {
	called := false
	items := []Item{
		{Title: "one", Description: "first", Disabled: true, OnSelect: func() tea.Cmd {
			called = true
			return nil
		}},
	}
	s := New("Test menu", items)
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if called {
		t.Fatal("expected OnSelect NOT to be called for a disabled item")
	}
	if cmd != nil {
		t.Fatalf("expected a nil cmd, got %v", cmd)
	}
}

func TestScreen_Esc_PopsScreen(t *testing.T) {
	s := New("Test menu", []Item{{Title: "one"}})
	s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a non-nil cmd from Update")
	}
	if _, ok := cmd().(tui.PopScreenMsg); !ok {
		t.Fatalf("expected the cmd to produce a PopScreenMsg, got %T", cmd())
	}
}

func TestScreen_SatisfiesTuiScreen(t *testing.T) {
	var _ tui.Screen = New("Test menu", nil)
}
