// internal/tui/healthui/watch.go
package healthui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/health"
	"github.com/skygrime35/mcli/internal/platform"
	"github.com/skygrime35/mcli/internal/tui"
)

const watchInterval = 5 * time.Second

type watchScreen struct {
	cancel  context.CancelFunc
	ch      <-chan health.Snapshot
	text    string
	stopped bool
}

type watchChanMsg struct {
	snap health.Snapshot
	ok   bool
}

func newWatchScreen() tui.Screen {
	return &watchScreen{}
}

func (s *watchScreen) Title() string { return "PC Health" }

func waitForSnapshot(ch <-chan health.Snapshot) tea.Cmd {
	return func() tea.Msg {
		snap, ok := <-ch
		return watchChanMsg{snap: snap, ok: ok}
	}
}

func (s *watchScreen) Init() tea.Cmd {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.ch = health.Watch(ctx, platform.Detect(), watchInterval)
	return waitForSnapshot(s.ch)
}

func (s *watchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case watchChanMsg:
		if !msg.ok {
			return s, nil
		}
		s.text = health.FormatSnapshot(msg.snap)
		return s, waitForSnapshot(s.ch)
	case tea.KeyMsg:
		if !s.stopped && (msg.String() == "q" || msg.String() == "esc") {
			s.stop()
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *watchScreen) View() string {
	if s.text == "" {
		return "Collecting health snapshot..."
	}
	return s.text + "\n" + tui.HelpStyle.Render("press q to stop watching and go back")
}

func (s *watchScreen) stop() {
	if s.stopped {
		return
	}
	s.stopped = true
	if s.cancel != nil {
		s.cancel()
	}
}

// Close implements tui.Closer, so RootModel can stop the background
// health.Watch goroutine on ctrl+c, not just q/esc - matching the same
// mechanism shareui.runningScreen uses for its running HTTP server.
func (s *watchScreen) Close() {
	s.stop()
}

var _ tui.Screen = (*watchScreen)(nil)
var _ tui.Closer = (*watchScreen)(nil)
