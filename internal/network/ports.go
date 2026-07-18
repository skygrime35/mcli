package network

import (
	"os"
	"strconv"
	"strings"
)

// ListeningPort is one locally-bound listening (TCP) or bound (UDP) port.
type ListeningPort struct {
	Protocol string
	Port     uint16
}

const tcpListenState = "0A"

// parseProcNet parses the content of /proc/net/tcp, /proc/net/tcp6,
// /proc/net/udp, or /proc/net/udp6 (they share the same column layout).
// When listeningOnly is true, only rows whose state column equals "0A"
// (TCP_LISTEN) are included - used for TCP. UDP has no listen state, so
// callers pass listeningOnly=false to include every bound row.
func parseProcNet(content string, protocol string, listeningOnly bool) []ListeningPort {
	var ports []ListeningPort
	lines := strings.Split(content, "\n")
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		if listeningOnly && !strings.EqualFold(fields[3], tcpListenState) {
			continue
		}
		localAddr := fields[1]
		parts := strings.Split(localAddr, ":")
		if len(parts) != 2 {
			continue
		}
		portVal, err := strconv.ParseUint(parts[1], 16, 16)
		if err != nil {
			continue
		}
		ports = append(ports, ListeningPort{Protocol: protocol, Port: uint16(portVal)})
	}
	return ports
}

// ListeningPorts reports every TCP (LISTEN state) and UDP (bound) local
// port found in /proc/net/{tcp,tcp6,udp,udp6}. Missing files (e.g. no
// IPv6 support) are silently skipped rather than treated as an error.
func ListeningPorts() ([]ListeningPort, error) {
	sources := []struct {
		path          string
		protocol      string
		listeningOnly bool
	}{
		{"/proc/net/tcp", "tcp", true},
		{"/proc/net/tcp6", "tcp", true},
		{"/proc/net/udp", "udp", false},
		{"/proc/net/udp6", "udp", false},
	}

	var all []ListeningPort
	for _, src := range sources {
		data, err := os.ReadFile(src.path)
		if err != nil {
			continue
		}
		all = append(all, parseProcNet(string(data), src.protocol, src.listeningOnly)...)
	}
	return all, nil
}
