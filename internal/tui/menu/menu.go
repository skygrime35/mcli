package menu

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/tui"
)

// Item is one entry in a menu. Disabled entries render dimmed and cannot
// be activated. OnSelect returns the tea.Cmd to run when this item is
// chosen (typically tui.PushScreen(...) to navigate).
type Item struct {
	Title       string
	Description string
	Disabled    bool
	OnSelect    func() tea.Cmd
}

// Item must satisfy list.Item, whose only requirement is FilterValue.
func (i Item) FilterValue() string { return i.Title }

// Screen is a generic, reusable list-based menu implementing tui.Screen.
type Screen struct {
	title string
	list  list.Model
}

func New(title string, items []Item) *Screen {
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = it
	}

	l := list.New(listItems, newDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &Screen{title: title, list: l}
}

func (s *Screen) Title() string { return s.title }
func (s *Screen) Init() tea.Cmd { return nil }

// Items exposes the underlying list's items, mainly so callers outside
// this package (e.g. tests verifying a screen rebuilds itself) can
// inspect the current item count without reaching into the unexported
// list field.
func (s *Screen) Items() []list.Item { return s.list.Items() }

func (s *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width, msg.Height)
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := s.list.SelectedItem().(Item)
			if !ok || selected.Disabled || selected.OnSelect == nil {
				return s, nil
			}
			return s, selected.OnSelect()
		case "esc":
			return s, tui.PopScreen()
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s *Screen) View() string { return s.list.View() }

var _ tui.Screen = (*Screen)(nil)
