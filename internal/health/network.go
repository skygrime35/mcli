// internal/health/network.go
package health

import (
	"os"
	"os/exec"
	"strings"
)

func parseInterfaceList(lsOutput string) []string {
	var names []string
	for _, name := range strings.Fields(lsOutput) {
		if name == "lo" {
			continue
		}
		names = append(names, name)
	}
	return names
}

func primaryInterface() string {
	out, err := exec.Command("ip", "route").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "default") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				return fields[4]
			}
		}
	}
	return ""
}

func defaultGateway() string {
	out, err := exec.Command("ip", "route").Output()
	if err != nil {
		return "None"
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "default") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2]
			}
		}
	}
	return "None"
}

func interfaceIP(name string) string {
	out, err := exec.Command("ip", "-4", "addr", "show", name).Output()
	if err != nil {
		return "No IP"
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "inet ") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == "inet" && i+1 < len(fields) {
					addr := fields[i+1]
					if idx := strings.Index(addr, "/"); idx != -1 {
						addr = addr[:idx]
					}
					return addr
				}
			}
		}
	}
	return "No IP"
}

func pingReachable(target string) bool {
	err := exec.Command("ping", "-c", "1", "-W", "2", target).Run()
	return err == nil
}

func CollectNetwork() NetworkInfo {
	entries, err := os.ReadDir("/sys/class/net")
	var names []string
	if err == nil {
		var sb strings.Builder
		for _, e := range entries {
			sb.WriteString(e.Name())
			sb.WriteString("\n")
		}
		names = parseInterfaceList(sb.String())
	}

	primary := primaryInterface()

	var interfaces []NetworkInterface
	for _, name := range names {
		state := "unknown"
		if data, err := os.ReadFile("/sys/class/net/" + name + "/operstate"); err == nil {
			state = strings.TrimSpace(string(data))
		}
		mac := ""
		if data, err := os.ReadFile("/sys/class/net/" + name + "/address"); err == nil {
			mac = strings.TrimSpace(string(data))
		}
		status := StatusInfo
		switch state {
		case "up":
			status = StatusGood
		case "down":
			status = StatusWarning
		}
		interfaces = append(interfaces, NetworkInterface{
			Name:    name,
			IP:      interfaceIP(name),
			State:   state,
			Primary: name == primary,
			MAC:     mac,
			Status:  status,
		})
	}

	internetOK := pingReachable("8.8.8.8")
	dnsOK := pingReachable("google.com")

	internetStatus := StatusCritical
	if internetOK {
		internetStatus = StatusGood
	}
	dnsStatus := StatusWarning
	if dnsOK {
		dnsStatus = StatusGood
	}

	return NetworkInfo{
		Interfaces:     interfaces,
		Gateway:        defaultGateway(),
		InternetOK:     internetOK,
		InternetStatus: internetStatus,
		DNSOK:          dnsOK,
		DNSStatus:      dnsStatus,
	}
}
