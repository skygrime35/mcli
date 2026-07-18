// internal/health/system.go
package health

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func parseOSRelease(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			v := strings.TrimPrefix(line, "PRETTY_NAME=")
			return strings.Trim(v, "\"")
		}
	}
	return "Unknown"
}

func CollectSystem() SystemInfo {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	osName := "Unknown"
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		osName = parseOSRelease(string(data))
	}

	kernel := "Unknown"
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		kernel = strings.TrimSpace(string(out))
	}

	uptime := "Unknown"
	if out, err := exec.Command("uptime", "-p").Output(); err == nil {
		uptime = strings.TrimPrefix(strings.TrimSpace(string(out)), "up ")
	}

	lastBoot := "Unknown"
	if out, err := exec.Command("who", "-b").Output(); err == nil {
		fields := strings.Fields(string(out))
		if len(fields) >= 4 {
			lastBoot = fields[2] + " " + fields[3]
		}
	}

	loggedUsers := 0
	if out, err := exec.Command("who").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			loggedUsers = 0
		} else {
			loggedUsers = len(lines)
		}
	}

	return SystemInfo{
		Hostname:     hostname,
		OS:           osName,
		Kernel:       kernel,
		Architecture: runtime.GOARCH,
		Uptime:       uptime,
		LastBoot:     lastBoot,
		LoggedUsers:  loggedUsers,
	}
}
