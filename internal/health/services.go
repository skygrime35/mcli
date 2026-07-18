package health

import (
	"os/exec"
	"strings"

	"github.com/skygrime35/mcli/internal/platform"
)

var knownServices = []string{"ssh", "sshd", "systemd-resolved", "NetworkManager", "cron", "rsyslog"}

func parseFailedUnits(output string) []string {
	var names []string
	for _, line := range strings.Split(strings.TrimRight(output, "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			names = append(names, fields[0])
		}
	}
	return names
}

// parseUnitFileNames parses the output of `systemctl list-unit-files` into
// a set of unit names with the ".service" suffix stripped, so membership
// can be checked for each known service without re-running the (slow)
// enumeration once per service.
func parseUnitFileNames(output string) map[string]bool {
	names := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			names[strings.TrimSuffix(fields[0], ".service")] = true
		}
	}
	return names
}

func CollectServices(caps platform.Capabilities) ServicesInfo {
	if !caps.Systemd {
		return ServicesInfo{}
	}

	unitFilesOut, _ := exec.Command("systemctl", "list-unit-files").Output()
	existingUnits := parseUnitFileNames(string(unitFilesOut))

	var services []ServiceStatus
	for _, name := range knownServices {
		if !existingUnits[name] {
			continue
		}
		state, status := "stopped_disabled", StatusInfo
		if exec.Command("systemctl", "is-active", "--quiet", name).Run() == nil {
			state, status = "running", StatusGood
		} else if exec.Command("systemctl", "is-enabled", "--quiet", name).Run() == nil {
			state, status = "stopped_enabled", StatusWarning
		}
		services = append(services, ServiceStatus{Name: name, State: state, Status: status})
	}

	failedOut, _ := exec.Command("systemctl", "--failed", "--no-legend").Output()
	failedNames := parseFailedUnits(string(failedOut))
	failedStatus := StatusGood
	if len(failedNames) > 0 {
		failedStatus = StatusCritical
	}

	return ServicesInfo{
		Services:     services,
		FailedCount:  len(failedNames),
		FailedNames:  failedNames,
		FailedStatus: failedStatus,
	}
}
