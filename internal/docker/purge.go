// internal/docker/purge.go
package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func listIDs(args ...string) ([]string, error) {
	out, err := exec.Command("docker", args...).Output()
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			ids = append(ids, line)
		}
	}
	return ids, nil
}

func dedup(ids []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

// parseCustomNetworks filters the `docker network ls --format '{{.ID}} {{.Name}}'`
// output down to networks that aren't Docker's three built-in defaults.
func parseCustomNetworks(output string) []string {
	var ids []string
	for _, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		id, name := fields[0], fields[1]
		if name == "bridge" || name == "host" || name == "none" {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func stopAllContainers(ch chan<- ProgressMsg) error {
	ids, err := listIDs("ps", "-q")
	if err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("listing running containers: %w", err)}
		return err
	}
	if len(ids) == 0 {
		ch <- ProgressMsg{Text: "No running containers."}
		return nil
	}
	ch <- ProgressMsg{Text: "Stopping running containers..."}
	if err := exec.Command("docker", append([]string{"stop"}, ids...)...).Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("stopping containers: %w", err)}
		return err
	}
	ch <- ProgressMsg{Text: fmt.Sprintf("Stopped %d container(s).", len(ids))}
	return nil
}

func removeAllContainers(ch chan<- ProgressMsg) error {
	ids, err := listIDs("ps", "-aq")
	if err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("listing containers: %w", err)}
		return err
	}
	if len(ids) == 0 {
		ch <- ProgressMsg{Text: "No containers to remove."}
		return nil
	}
	ch <- ProgressMsg{Text: "Removing all containers..."}
	if err := exec.Command("docker", append([]string{"rm", "-f"}, ids...)...).Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("removing containers: %w", err)}
		return err
	}
	ch <- ProgressMsg{Text: fmt.Sprintf("Removed %d container(s).", len(ids))}
	return nil
}

func removeAllImages(ch chan<- ProgressMsg) error {
	ids, err := listIDs("images", "-q")
	if err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("listing images: %w", err)}
		return err
	}
	ids = dedup(ids)
	if len(ids) == 0 {
		ch <- ProgressMsg{Text: "No images to remove."}
		return nil
	}
	ch <- ProgressMsg{Text: "Removing all images..."}
	if err := exec.Command("docker", append([]string{"rmi", "-f"}, ids...)...).Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("removing images: %w", err)}
		return err
	}
	ch <- ProgressMsg{Text: "All images removed."}
	return nil
}

func removeAllVolumes(ch chan<- ProgressMsg) error {
	ids, err := listIDs("volume", "ls", "-q")
	if err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("listing volumes: %w", err)}
		return err
	}
	if len(ids) == 0 {
		ch <- ProgressMsg{Text: "No volumes to remove."}
		return nil
	}
	ch <- ProgressMsg{Text: "Removing all volumes..."}
	if err := exec.Command("docker", append([]string{"volume", "rm"}, ids...)...).Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("removing volumes: %w", err)}
		return err
	}
	ch <- ProgressMsg{Text: "All volumes removed."}
	return nil
}

func removeCustomNetworks(ch chan<- ProgressMsg) error {
	out, err := exec.Command("docker", "network", "ls", "--format", "{{.ID}} {{.Name}}").Output()
	if err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("listing networks: %w", err)}
		return err
	}
	ids := parseCustomNetworks(string(out))
	if len(ids) == 0 {
		ch <- ProgressMsg{Text: "No custom networks to remove."}
		return nil
	}
	ch <- ProgressMsg{Text: "Removing custom networks..."}
	if err := exec.Command("docker", append([]string{"network", "rm"}, ids...)...).Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("removing networks: %w", err)}
		return err
	}
	ch <- ProgressMsg{Text: fmt.Sprintf("Removed %d custom network(s).", len(ids))}
	return nil
}

// pruneBuilderCache and systemPrune are TOLERANT of failure (matching
// clear_docker.sh's `|| true` on exactly these two commands) - they
// report an error if one occurs but never halt the calling sequence.
func pruneBuilderCache(ch chan<- ProgressMsg) {
	if err := exec.Command("docker", "builder", "prune", "-af").Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("pruning build cache (non-fatal): %w", err)}
		return
	}
	ch <- ProgressMsg{Text: "Build cache pruned."}
}

func systemPrune(ch chan<- ProgressMsg) {
	if err := exec.Command("docker", "system", "prune", "-a", "--volumes", "-f").Run(); err != nil {
		ch <- ProgressMsg{Err: fmt.Errorf("system prune (non-fatal): %w", err)}
		return
	}
	ch <- ProgressMsg{Text: "System pruned."}
}

// FullPurge stops and removes every container, image, volume, and custom
// network, then prunes the build cache and runs a final system prune - the
// 7-step sequence from clear_docker.sh. Steps 1-5 halt the whole sequence
// on the first error; the last two tolerate failure and continue.
func FullPurge() <-chan ProgressMsg {
	ch := make(chan ProgressMsg)
	go func() {
		defer close(ch)
		if !CheckDaemon() {
			ch <- ProgressMsg{Err: fmt.Errorf("docker daemon is not accessible - is Docker running?")}
			return
		}
		if stopAllContainers(ch) != nil {
			return
		}
		if removeAllContainers(ch) != nil {
			return
		}
		if removeAllImages(ch) != nil {
			return
		}
		if removeAllVolumes(ch) != nil {
			return
		}
		if removeCustomNetworks(ch) != nil {
			return
		}
		pruneBuilderCache(ch)
		systemPrune(ch)
		ch <- ProgressMsg{Text: "Docker purge complete."}
	}()
	return ch
}
