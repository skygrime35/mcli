// internal/docker/list.go
package docker

import (
	"os/exec"
	"strings"
)

func parseContainerList(output string) []Container {
	var containers []Container
	for _, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		containers = append(containers, Container{
			ID:     parts[0],
			Name:   parts[1],
			Status: parts[2],
			Image:  parts[3],
		})
	}
	return containers
}

// ListContainers lists all containers (running and stopped), matching
// `docker ps -a`.
func ListContainers() ([]Container, error) {
	out, err := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Status}}|{{.Image}}").Output()
	if err != nil {
		return nil, err
	}
	return parseContainerList(string(out)), nil
}
