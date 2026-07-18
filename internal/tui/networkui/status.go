// internal/tui/networkui/status.go
package networkui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skygrime35/mcli/internal/network"
	"github.com/skygrime35/mcli/internal/tui"
)

type statusScreen struct {
	ifaces    []network.InterfaceInfo
	ports     []network.ListeningPort
	connected bool
	done      bool
	err       error
}

type statusLoadedMsg struct {
	ifaces    []network.InterfaceInfo
	ports     []network.ListeningPort
	connected bool
	err       error
}

func newStatusScreen() tui.Screen {
	return &statusScreen{}
}

func (s *statusScreen) Title() string { return "Network Status" }

func (s *statusScreen) Init() tea.Cmd {
	return func() tea.Msg {
		ifaces, err := network.Interfaces()
		if err != nil {
			return statusLoadedMsg{err: err}
		}
		ports, err := network.ListeningPorts()
		if err != nil {
			return statusLoadedMsg{err: err}
		}
		return statusLoadedMsg{ifaces: ifaces, ports: ports, connected: network.CheckConnectivity()}
	}
}

func (s *statusScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusLoadedMsg:
		s.ifaces = msg.ifaces
		s.ports = msg.ports
		s.connected = msg.connected
		s.err = msg.err
		s.done = true
		return s, nil
	case tea.KeyMsg:
		if s.done && (msg.String() == "esc" || msg.String() == "enter") {
			return s, tui.PopScreen()
		}
	}
	return s, nil
}

func (s *statusScreen) View() string {
	if !s.done {
		return "Loading network status..."
	}
	if s.err != nil {
		return "Error: " + s.err.Error() + "\n\n" + tui.HelpStyle.Render("press enter/esc to go back")
	}

	out := "Interfaces:\n"
	for _, ifc := range s.ifaces {
		out += fmt.Sprintf("  %s: %v\n", ifc.Name, ifc.IPs)
	}

	out += fmt.Sprintf("\nListening ports (%d):\n", len(s.ports))
	for _, p := range s.ports {
		out += fmt.Sprintf("  %s/%d\n", p.Protocol, p.Port)
	}

	out += "\nInternet connectivity: "
	if s.connected {
		out += "OK\n"
	} else {
		out += "UNREACHABLE\n"
	}

	out += "\n" + tui.HelpStyle.Render("press enter/esc to go back")
	return out
}

var _ tui.Screen = (*statusScreen)(nil)
