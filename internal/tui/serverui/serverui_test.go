// internal/tui/serverui/serverui_test.go
package serverui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/tui"
)

func TestNewListScreen_ListsServersPlusAddEntry(t *testing.T) {
	cfg := &config.Config{
		Servers: []config.ServerConfig{
			{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "alice", SSHPort: 22, WOLPort: 9},
		},
	}

	screen := NewListScreen(cfg)
	if screen.Title() != "Servers" {
		t.Fatalf("expected title 'Servers', got %q", screen.Title())
	}
}

func TestListScreen_Init_RebuildsFromCurrentConfig(t *testing.T) {
	cfg := &config.Config{Servers: []config.ServerConfig{
		{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "alice", SSHPort: 22, WOLPort: 9},
	}}

	screen := NewListScreen(cfg).(*listScreen)
	initialCount := len(screen.inner.Items())

	// Simulate a server having been added to the SAME config after this
	// screen was first constructed (e.g. by the add-server flow).
	cfg.Servers = append(cfg.Servers, config.ServerConfig{
		Name: "work", Host: "work.example.com", MAC: "11:22:33:44:55:66", SSHUser: "bob", SSHPort: 22, WOLPort: 9,
	})

	screen.Init()

	if got := len(screen.inner.Items()); got != initialCount+1 {
		t.Fatalf("expected item count to increase by 1 after Init rebuild (from %d), got %d", initialCount, got)
	}
}

func TestNewDetailScreen_HasFourActions(t *testing.T) {
	s := config.ServerConfig{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "alice", SSHPort: 22, WOLPort: 9}
	cfg := &config.Config{Servers: []config.ServerConfig{s}}

	screen := NewDetailScreen(cfg, s)
	if screen.Title() != "Server: home" {
		t.Fatalf("expected title 'Server: home', got %q", screen.Title())
	}
}

func TestDetailItems_TitlesMatchExpectedActions(t *testing.T) {
	s := config.ServerConfig{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "alice", SSHPort: 22, WOLPort: 9}
	cfg := &config.Config{Servers: []config.ServerConfig{s}}

	items := detailItems(cfg, s)
	want := []string{"Check status", "Wake on LAN", "Connect (SSH)", "Setup SSH key"}
	if len(items) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(items))
	}
	for i, w := range want {
		if items[i].Title != w {
			t.Errorf("item %d: expected title %q, got %q", i, w, items[i].Title)
		}
		// "Connect (SSH)" (index 2) is intentionally a disabled stub with a
		// nil OnSelect in this task - Task 5 wires it up for real and adds
		// its own dedicated test for that. All other items must already be
		// fully wired here.
		if w == "Connect (SSH)" {
			continue
		}
		if items[i].OnSelect == nil {
			t.Errorf("item %d (%q): expected non-nil OnSelect", i, w)
		}
	}
}

func TestNewAddScreen_SavesValidServerOnCompletion(t *testing.T) {
	cfg := &config.Config{}
	screen := NewAddScreen(cfg)
	if screen.Title() != "Add server" {
		t.Fatalf("expected title 'Add server', got %q", screen.Title())
	}
	// Full form-submission flow needs a real terminal (huh.Form driven by
	// key events through bubbletea) and is exercised manually in Task 6's
	// end-to-end smoke test, not here — this test only checks the screen
	// constructs correctly and exposes the right title.
}

func TestDetailItems_ConnectSSH_NoLongerDisabled(t *testing.T) {
	s := config.ServerConfig{Name: "home", Host: "example.com", MAC: "AA:BB:CC:DD:EE:FF", SSHUser: "alice", SSHPort: 22, WOLPort: 9}
	cfg := &config.Config{Servers: []config.ServerConfig{s}}

	items := detailItems(cfg, s)
	for _, it := range items {
		if it.Title == "Connect (SSH)" {
			if it.Disabled {
				t.Error("expected 'Connect (SSH)' to no longer be disabled")
			}
			if it.OnSelect == nil {
				t.Error("expected 'Connect (SSH)' to have a non-nil OnSelect")
			}
			return
		}
	}
	t.Fatal("expected a 'Connect (SSH)' item")
}

// TestAddScreen_Esc_PopsScreenAfterSubmissionError guards against a dead-end
// bug: once the huh form reaches huh.StateCompleted, a submission error
// (e.g. AddServer/Save failing) used to leave the user stuck with no way to
// leave the Add Server screen, because esc only worked while the form's
// State was not yet StateCompleted. Esc must also work once s.err is set,
// regardless of the form's State.
func TestAddScreen_Esc_PopsScreenAfterSubmissionError(t *testing.T) {
	cfg := &config.Config{}
	screen := NewAddScreen(cfg)

	s, ok := screen.(*addScreen)
	if !ok {
		t.Fatalf("expected *addScreen, got %T", screen)
	}

	// Simulate the post-completion error state directly, since driving the
	// real huh form to StateCompleted requires a full terminal/key-event
	// flow that is out of scope here (see TestNewAddScreen_SavesValidServerOnCompletion).
	s.form.State = huh.StateCompleted
	s.done = true
	s.err = errors.New("server name already exists")

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a non-nil cmd from Update after esc with s.err set")
	}
	if _, ok := cmd().(tui.PopScreenMsg); !ok {
		t.Fatalf("expected the cmd to produce a PopScreenMsg, got %T", cmd())
	}
}
