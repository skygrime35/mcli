package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type fakeScreen struct {
	name          string
	initCmd       tea.Cmd
	updates       int
	sawWindowSize bool
}

func (f *fakeScreen) Init() tea.Cmd { return f.initCmd }
func (f *fakeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	f.updates++
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		f.sawWindowSize = true
	}
	return f, nil
}
func (f *fakeScreen) View() string  { return f.name }
func (f *fakeScreen) Title() string { return f.name }

func TestRootModel_PushPop(t *testing.T) {
	root := NewRootModel(&fakeScreen{name: "main"})

	updated, _ := root.Update(PushScreenMsg{Screen: &fakeScreen{name: "detail"}})
	root = updated.(RootModel)
	if root.top().Title() != "detail" {
		t.Fatalf("expected top screen to be 'detail', got %q", root.top().Title())
	}

	updated, _ = root.Update(PopScreenMsg{})
	root = updated.(RootModel)
	if root.top().Title() != "main" {
		t.Fatalf("expected top screen to be 'main' after pop, got %q", root.top().Title())
	}
}

func TestRootModel_PopAtRootQuits(t *testing.T) {
	root := NewRootModel(&fakeScreen{name: "main"})

	_, cmd := root.Update(PopScreenMsg{})
	if cmd == nil {
		t.Fatal("expected a tea.Quit command when popping the root screen, got nil")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestRootModel_UpdateDelegatesToTopScreen(t *testing.T) {
	top := &fakeScreen{name: "top"}
	root := NewRootModel(top)

	root.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if top.updates != 1 {
		t.Fatalf("expected top screen's Update to be called once, got %d", top.updates)
	}
}

func TestRootModel_PushForwardsKnownWindowSize(t *testing.T) {
	root := NewRootModel(&fakeScreen{name: "main"})

	updated, _ := root.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	root = updated.(RootModel)

	pushed := &fakeScreen{name: "detail"}
	updated, _ = root.Update(PushScreenMsg{Screen: pushed})
	root = updated.(RootModel)

	if !pushed.sawWindowSize {
		t.Fatal("expected the newly-pushed screen to be sized immediately with the known window size")
	}
}

func TestRootModel_PushBeforeAnySize_DoesNotPanic(t *testing.T) {
	// No WindowSizeMsg has been delivered yet (width/height are zero) —
	// activateTop must not try to deliver a bogus zero-size message.
	root := NewRootModel(&fakeScreen{name: "main"})
	pushed := &fakeScreen{name: "detail"}

	updated, _ := root.Update(PushScreenMsg{Screen: pushed})
	root = updated.(RootModel)

	if pushed.sawWindowSize {
		t.Fatal("expected no WindowSizeMsg to be delivered when none is known yet")
	}
	if root.top().Title() != "detail" {
		t.Fatalf("expected push to still work, got top=%q", root.top().Title())
	}
}
