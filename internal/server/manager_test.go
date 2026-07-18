// internal/server/manager_test.go
package server

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/skygrime35/mcli/internal/config"
)

func TestCheckStatus_Online(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	defer ln.Close()

	host, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)

	status := CheckStatus(context.Background(), host, port, 2*time.Second)
	if !status.Online {
		t.Error("expected Online=true for a reachable port")
	}
}

func TestCheckStatus_Offline(t *testing.T) {
	// Bind and immediately close: nothing listens on the port right after,
	// almost never valid on this host, so the connection reliably fails.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate a port: %v", err)
	}
	host, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	ln.Close()

	status := CheckStatus(context.Background(), host, port, 1*time.Second)
	if status.Online {
		t.Error("expected Online=false for a closed port")
	}
}

func TestWakeOnLAN_AlreadyOnline(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	defer ln.Close()

	host, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)

	s := config.ServerConfig{
		Name: "test", Host: host, SSHPort: port,
		MAC: "AA:BB:CC:DD:EE:FF", WOLPort: 9,
	}

	ch := WakeOnLAN(context.Background(), s)
	var messages []string
	for msg := range ch {
		if msg.Err != nil {
			t.Fatalf("unexpected error: %v", msg.Err)
		}
		messages = append(messages, msg.Text)
	}

	if len(messages) != 1 || !strings.Contains(messages[0], "already online") {
		t.Errorf("expected a single 'already online' message, got %v", messages)
	}
}
