// internal/tui/serverui/list.go
package serverui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/tui"
	"github.com/skygrime35/mcli/internal/tui/menu"
)

// listScreen wraps a menu.Screen and rebuilds its items from cfg every
// time it's (re)activated via Init - so returning here after adding a
// server, or after viewing a server's detail screen, always reflects the
// current config rather than a stale snapshot taken at construction time.
type listScreen struct {
	cfg   *config.Config
	inner *menu.Screen
}

// NewListScreen lists the servers configured in cfg, plus an "Add server"
// entry. Selecting a server pushes its detail/actions screen.
func NewListScreen(cfg *config.Config) tui.Screen {
	return &listScreen{cfg: cfg, inner: menu.New("Servers", buildServerItems(cfg))}
}

func buildServerItems(cfg *config.Config) []menu.Item {
	items := make([]menu.Item, 0, len(cfg.Servers)+1)
	for _, s := range cfg.Servers {
		s := s // capture loop variable for the closure below
		items = append(items, menu.Item{
			Title:       s.Name,
			Description: s.SSHUser + "@" + s.Host,
			OnSelect: func() tea.Cmd {
				return tui.PushScreen(NewDetailScreen(cfg, s))
			},
		})
	}
	items = append(items, menu.Item{
		Title:       "+ Add server",
		Description: "Add a new server to the config",
		OnSelect: func() tea.Cmd {
			return tui.PushScreen(NewAddScreen(cfg))
		},
	})
	return items
}

func (s *listScreen) Title() string { return s.inner.Title() }

func (s *listScreen) Init() tea.Cmd {
	s.inner = menu.New("Servers", buildServerItems(s.cfg))
	return s.inner.Init()
}

func (s *listScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := s.inner.Update(msg)
	if m, ok := updated.(*menu.Screen); ok {
		s.inner = m
	}
	return s, cmd
}

func (s *listScreen) View() string { return s.inner.View() }

var _ tui.Screen = (*listScreen)(nil)
