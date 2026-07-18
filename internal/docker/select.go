// internal/docker/select.go
package docker

import (
	"fmt"
	"os/exec"
)

// SelectPurge stops and removes exactly the given container IDs, one at a
// time - a container is only removed if its stop succeeded, matching the
// old app.py's Select Purge flow. A failure on one container doesn't stop
// the others from being attempted.
func SelectPurge(ids []string) <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)
		if !CheckDaemon() {
			ch <- ProgressMsg{Err: fmt.Errorf("docker daemon is not accessible - is Docker running?")}
			return
		}
		ch <- ProgressMsg{Text: fmt.Sprintf("Purging %d container(s)...", len(ids))}
		for _, id := range ids {
			if err := exec.Command("docker", "stop", id).Run(); err != nil {
				ch <- ProgressMsg{Text: fmt.Sprintf("Failed to stop %s: %v", id, err)}
				continue
			}
			ch <- ProgressMsg{Text: fmt.Sprintf("Stopped %s", id)}

			if err := exec.Command("docker", "rm", "-f", id).Run(); err != nil {
				ch <- ProgressMsg{Text: fmt.Sprintf("Failed to remove %s: %v", id, err)}
				continue
			}
			ch <- ProgressMsg{Text: fmt.Sprintf("Removed %s", id)}
		}
		ch <- ProgressMsg{Text: "Selected purge complete."}
	}()
	return ch
}
