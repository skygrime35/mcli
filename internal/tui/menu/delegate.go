package menu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// delegate is a custom list.ItemDelegate. list.DefaultDelegate.Render reads
// unexported list.Model fields, so it can't be subclassed to add disabled/
// dimmed rendering — this implements the 4-method ItemDelegate interface
// directly instead.
type delegate struct {
	normalTitle   lipgloss.Style
	normalDesc    lipgloss.Style
	selectedTitle lipgloss.Style
	selectedDesc  lipgloss.Style
	dimmedTitle   lipgloss.Style
	dimmedDesc    lipgloss.Style
}

func newDelegate() delegate {
	return delegate{
		normalTitle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Padding(0, 0, 0, 2),
		normalDesc: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
			Padding(0, 0, 0, 2),
		selectedTitle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("212")).
			Foreground(lipgloss.Color("212")).
			Padding(0, 0, 0, 1),
		selectedDesc: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("212")).
			Foreground(lipgloss.Color("240")).
			Padding(0, 0, 0, 1),
		dimmedTitle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 0, 0, 2),
		dimmedDesc:  lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Padding(0, 0, 0, 2),
	}
}

func (d delegate) Height() int                              { return 2 }
func (d delegate) Spacing() int                              { return 1 }
func (d delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(Item)
	if !ok {
		return
	}

	title, desc := it.Title, it.Description
	isSelected := index == m.Index()

	switch {
	case it.Disabled:
		title = d.dimmedTitle.Render(title + " (coming soon)")
		desc = d.dimmedDesc.Render(desc)
	case isSelected:
		title = d.selectedTitle.Render(title)
		desc = d.selectedDesc.Render(desc)
	default:
		title = d.normalTitle.Render(title)
		desc = d.normalDesc.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc) //nolint:errcheck
}
