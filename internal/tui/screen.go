// internal/tui/screen.go
package tui

import tea "github.com/charmbracelet/bubbletea"

// Screen is one screen in the navigation stack. Concrete screens (menus,
// forms, progress views) implement this by embedding the standard
// tea.Model methods and adding Title.
type Screen interface {
	tea.Model
	Title() string
}

// Closer is implemented by screens that hold a live resource (e.g. a
// running server) needing explicit cleanup before the program exits.
// RootModel checks for this via a type assertion - most screens don't
// need it and are unaffected.
type Closer interface {
	Close()
}

// PushScreenMsg and PopScreenMsg are the only two messages RootModel
// intercepts itself; every other message is delegated to the top screen.
// Child screens never receive these messages directly — they only ever
// construct and return them via PushScreen/PopScreen as a tea.Cmd.
type PushScreenMsg struct{ Screen Screen }
type PopScreenMsg struct{}

func PushScreen(s Screen) tea.Cmd { return func() tea.Msg { return PushScreenMsg{Screen: s} } }
func PopScreen() tea.Cmd          { return func() tea.Msg { return PopScreenMsg{} } }

// RootModel holds the navigation stack and delegates Update/View to
// whichever screen is on top. It also remembers the last known terminal
// size so newly-activated screens (pushed, or re-exposed by a pop) can be
// sized immediately, instead of waiting for a resize event that may never
// come.
type RootModel struct {
	stack         []Screen
	width, height int
}

func NewRootModel(initial Screen) RootModel { return RootModel{stack: []Screen{initial}} }
func (r RootModel) top() Screen             { return r.stack[len(r.stack)-1] }
func (r RootModel) Init() tea.Cmd           { return r.top().Init() }

// activateTop (re)initializes the current top screen and, if a window
// size is already known, delivers it immediately (synchronously) so the
// screen is correctly sized before its first render. Called after every
// push and every pop.
func (r RootModel) activateTop() tea.Cmd {
	cmd := r.top().Init()
	if r.width > 0 && r.height > 0 {
		updated, sizeCmd := r.top().Update(tea.WindowSizeMsg{Width: r.width, Height: r.height})
		if scr, ok := updated.(Screen); ok {
			r.stack[len(r.stack)-1] = scr
		}
		cmd = tea.Batch(cmd, sizeCmd)
	}
	return cmd
}

func (r RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width, r.height = msg.Width, msg.Height
		var cmds []tea.Cmd
		for i, s := range r.stack {
			updated, cmd := s.Update(msg)
			if scr, ok := updated.(Screen); ok {
				r.stack[i] = scr
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return r, tea.Batch(cmds...)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			for _, scr := range r.stack {
				if c, ok := scr.(Closer); ok {
					c.Close()
				}
			}
			return r, tea.Quit
		}
	case PushScreenMsg:
		r.stack = append(r.stack, msg.Screen)
		return r, r.activateTop()
	case PopScreenMsg:
		if len(r.stack) > 1 {
			r.stack = r.stack[:len(r.stack)-1]
			return r, r.activateTop()
		}
		// Popping the last screen (the main menu) quits the program.
		return r, tea.Quit
	}
	updated, cmd := r.top().Update(msg)
	if scr, ok := updated.(Screen); ok {
		r.stack[len(r.stack)-1] = scr
	}
	return r, cmd
}

func (r RootModel) View() string { return r.top().View() }

// NewProgram wraps the given initial screen in a RootModel and returns a
// ready-to-run bubbletea program, in the alt screen (full-terminal) mode.
func NewProgram(initial Screen) *tea.Program {
	return tea.NewProgram(NewRootModel(initial), tea.WithAltScreen())
}
