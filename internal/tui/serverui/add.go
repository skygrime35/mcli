// internal/tui/serverui/add.go
package serverui

import (
	"fmt"
	"net"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/skygrime35/mcli/internal/config"
	"github.com/skygrime35/mcli/internal/tui"
)

type addScreen struct {
	cfg  *config.Config
	form *huh.Form

	name, host, mac, sshUser, sshPort, wolPort string
	done                                       bool
	err                                        error
}

// NewAddScreen presents a form to collect a new server's details and
// saves it to cfg on successful submission.
func NewAddScreen(cfg *config.Config) tui.Screen {
	s := &addScreen{cfg: cfg, sshPort: "22", wolPort: "9"}
	s.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&s.name).Validate(requireNonEmpty("name")),
			huh.NewInput().Title("Host").Value(&s.host).Validate(requireNonEmpty("host")),
			huh.NewInput().Title("MAC address").Placeholder("AA:BB:CC:DD:EE:FF").Value(&s.mac).
				Validate(func(v string) error {
					if _, err := net.ParseMAC(v); err != nil {
						return fmt.Errorf("invalid MAC address: %w", err)
					}
					return nil
				}),
			huh.NewInput().Title("SSH user").Value(&s.sshUser).Validate(requireNonEmpty("ssh user")),
			huh.NewInput().Title("SSH port").Value(&s.sshPort).Validate(requirePort("ssh port")),
			huh.NewInput().Title("WOL port").Value(&s.wolPort).Validate(requirePort("wol port")),
		),
	)
	return s
}

func requireNonEmpty(field string) func(string) error {
	return func(v string) error {
		if v == "" {
			return fmt.Errorf("%s cannot be empty", field)
		}
		return nil
	}
}

func requirePort(field string) func(string) error {
	return func(v string) error {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 || p > 65535 {
			return fmt.Errorf("%s must be a number between 1 and 65535", field)
		}
		return nil
	}
}

func (s *addScreen) Title() string { return "Add server" }
func (s *addScreen) Init() tea.Cmd { return s.form.Init() }

func (s *addScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" && (s.form.State != huh.StateCompleted || s.err != nil) {
		return s, tui.PopScreen()
	}

	m, cmd := s.form.Update(msg)
	if f, ok := m.(*huh.Form); ok {
		s.form = f
	}

	if s.form.State == huh.StateCompleted && !s.done {
		s.done = true
		sshPort, _ := strconv.Atoi(s.sshPort)
		wolPort, _ := strconv.Atoi(s.wolPort)
		entry := config.ServerConfig{
			Name: s.name, Host: s.host, MAC: s.mac,
			SSHUser: s.sshUser, SSHPort: sshPort, WOLPort: wolPort,
		}
		if err := config.AddServer(s.cfg, entry); err != nil {
			s.err = err
			return s, nil
		}
		if err := config.Save(s.cfg); err != nil {
			s.err = err
			return s, nil
		}
		return s, tui.PopScreen()
	}

	return s, cmd
}

func (s *addScreen) View() string {
	if s.err != nil {
		return s.form.View() + "\n" + tui.ErrorStyle.Render("Error: "+s.err.Error()) +
			"\n" + tui.HelpStyle.Render("press esc to go back")
	}
	return s.form.View()
}

var _ tui.Screen = (*addScreen)(nil)
