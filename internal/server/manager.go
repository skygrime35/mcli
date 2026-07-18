// internal/server/manager.go
package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/skygrime35/mcli/internal/config"
)

type Status struct {
	Online bool
}

func CheckStatus(ctx context.Context, host string, port int, timeout time.Duration) Status {
	d := net.Dialer{Timeout: timeout}
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return Status{Online: false}
	}
	conn.Close()
	return Status{Online: true}
}

type ProgressMsg struct {
	Text string
	Err  error
}

const (
	wolAttempts    = 5
	wolInterval    = 2 * time.Second
	statusAttempts = 36
	statusInterval = 10 * time.Second
	statusTimeout  = 5 * time.Second
)

func WakeOnLAN(ctx context.Context, s config.ServerConfig) <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)
		if CheckStatus(ctx, s.Host, s.SSHPort, statusTimeout).Online {
			ch <- ProgressMsg{Text: "The server is already online and accessible via SSH."}
			return
		}
		ch <- ProgressMsg{Text: "Sending WOL requests..."}
		for i := 1; i <= wolAttempts; i++ {
			if err := SendMagicPacket(s.Host, s.WOLPort, s.MAC); err != nil {
				ch <- ProgressMsg{Err: fmt.Errorf("sending magic packet: %w", err)}
				return
			}
			ch <- ProgressMsg{Text: fmt.Sprintf("  WOL %d sent", i)}
			if i < wolAttempts {
				select {
				case <-time.After(wolInterval):
				case <-ctx.Done():
					return
				}
			}
		}
		ch <- ProgressMsg{Text: "Waiting for the server to start..."}
		for attempt := 1; attempt <= statusAttempts; attempt++ {
			ch <- ProgressMsg{Text: fmt.Sprintf("  Attempt %d/%d", attempt, statusAttempts)}
			if CheckStatus(ctx, s.Host, s.SSHPort, statusTimeout).Online {
				ch <- ProgressMsg{Text: "The server is online and accessible via SSH."}
				return
			}
			if attempt < statusAttempts {
				select {
				case <-time.After(statusInterval):
				case <-ctx.Done():
					return
				}
			}
		}
		ch <- ProgressMsg{Text: fmt.Sprintf("Reached maximum attempts (%d). The server might be offline.", statusAttempts)}
	}()
	return ch
}

func Connect(s config.ServerConfig) error {
	cmd := exec.Command("ssh", "-p", strconv.Itoa(s.SSHPort), fmt.Sprintf("%s@%s", s.SSHUser, s.Host))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SetupSSHKey(s config.ServerConfig) <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)
		home, err := os.UserHomeDir()
		if err != nil {
			ch <- ProgressMsg{Err: fmt.Errorf("resolving home dir: %w", err)}
			return
		}
		sshDir := filepath.Join(home, ".ssh")
		keyPath := filepath.Join(sshDir, "id_ed25519")
		pubPath := keyPath + ".pub"

		if _, err := os.Stat(pubPath); err == nil {
			ch <- ProgressMsg{Text: fmt.Sprintf("SSH key already exists: %s", pubPath)}
		} else {
			ch <- ProgressMsg{Text: "Generating SSH key (ed25519)..."}
			if err := os.MkdirAll(sshDir, 0o700); err != nil {
				ch <- ProgressMsg{Err: fmt.Errorf("creating ssh dir: %w", err)}
				return
			}
			cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "")
			if out, err := cmd.CombinedOutput(); err != nil {
				ch <- ProgressMsg{Err: fmt.Errorf("generating key: %w (%s)", err, out)}
				return
			}
			ch <- ProgressMsg{Text: "SSH key generated successfully."}
		}

		ch <- ProgressMsg{Text: fmt.Sprintf("Copying key to %s@%s...", s.SSHUser, s.Host)}
		cmd := exec.Command("ssh-copy-id", "-p", strconv.Itoa(s.SSHPort), fmt.Sprintf("%s@%s", s.SSHUser, s.Host))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			ch <- ProgressMsg{Err: fmt.Errorf("copying ssh key: %w", err)}
			return
		}
		ch <- ProgressMsg{Text: "SSH key copied successfully."}
	}()
	return ch
}
