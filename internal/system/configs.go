// internal/system/configs.go
package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// parseOrphanedConfigs parses `dpkg --list` output, returning package
// names for lines in "rc" status (removed, config files remain).
func parseOrphanedConfigs(dpkgOutput string) []string {
	var pkgs []string
	for _, line := range strings.Split(dpkgOutput, "\n") {
		if !strings.HasPrefix(line, "rc") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pkgs = append(pkgs, fields[1])
	}
	return pkgs
}

func analyzeOrphanedConfigs() (*Action, error) {
	out, err := exec.Command("dpkg", "--list").Output()
	if err != nil {
		return nil, fmt.Errorf("listing packages: %w", err)
	}
	orphaned := parseOrphanedConfigs(string(out))
	if len(orphaned) == 0 {
		return nil, nil
	}
	return &Action{
		Name:    "Purge orphaned configs",
		Tier:    TierWarning,
		Command: append([]string{"dpkg", "--purge"}, orphaned...),
		Reason:  fmt.Sprintf("%d package config(s) to remove", len(orphaned)),
	}, nil
}
