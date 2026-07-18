// internal/tui/progress/progress_test.go
package progress

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/server"
	"github.com/skygrime35/mcli/internal/tui"
)

func TestScreen_ConsumesChannelUntilClose(t *testing.T) {
	ch := make(chan server.ProgressMsg, 2)
	ch <- server.ProgressMsg{Text: "step 1"}
	ch <- server.ProgressMsg{Text: "step 2"}
	close(ch)

	s := New("Test progress", ch)
	cmd := s.Init()

	for i := 0; i < 2; i++ {
		msg := cmd()
		model, next := s.Update(msg)
		s = model.(*Screen)
		cmd = next
	}
	// Third receive observes the closed channel.
	msg := cmd()
	model, next := s.Update(msg)
	s = model.(*Screen)

	if len(s.lines) != 2 || s.lines[0] != "step 1" || s.lines[1] != "step 2" {
		t.Fatalf("expected 2 lines [step 1, step 2], got %v", s.lines)
	}
	if !s.closed {
		t.Fatal("expected screen to be marked closed after channel close")
	}
	if next != nil {
		t.Fatal("expected no further re-issued cmd once the channel is closed")
	}
}

func TestScreen_EnterAfterClosed_Pops(t *testing.T) {
	ch := make(chan server.ProgressMsg)
	close(ch)

	s := New("Test progress", ch)
	msg := s.Init()()
	model, _ := s.Update(msg)
	s = model.(*Screen)

	if !s.closed {
		t.Fatal("expected screen to be closed")
	}

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a non-nil cmd")
	}
	if _, ok := cmd().(tui.PopScreenMsg); !ok {
		t.Fatalf("expected PopScreenMsg, got %T", cmd())
	}
}

func TestScreen_EnterBeforeClosed_DoesNotPop(t *testing.T) {
	ch := make(chan server.ProgressMsg)
	s := New("Test progress", ch)

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("expected no cmd while still running (not yet closed)")
	}
}

func TestScreen_CapturesErr(t *testing.T) {
	ch := make(chan server.ProgressMsg, 1)
	testErr := &testError{"boom"}
	ch <- server.ProgressMsg{Err: testErr}
	close(ch)

	s := New("Test progress", ch)
	cmd := s.Init()
	msg := cmd()
	model, _ := s.Update(msg)
	s = model.(*Screen)

	if s.err == nil || s.err.Error() != "boom" {
		t.Fatalf("expected err 'boom', got %v", s.err)
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
