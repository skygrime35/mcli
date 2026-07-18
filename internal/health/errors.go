package health

import (
	"os/exec"
	"strings"

	"github.com/skygrime35/mcli/internal/platform"
)

func parseCount(output string) int {
	output = strings.TrimRight(output, "\n")
	if output == "" {
		return 0
	}
	return len(strings.Split(output, "\n"))
}

func CollectErrors(caps platform.Capabilities) ErrorsInfo {
	var info ErrorsInfo
	if caps.Journalctl {
		errOut, _ := exec.Command("journalctl", "-p", "err", "--since", "1 hour ago", "--no-pager").Output()
		info.RecentErrors = parseCount(string(errOut))
		info.ErrorsStatus = statusForHigh(float64(info.RecentErrors), journalErrorsWarning, journalErrorsCritical)

		warnOut, _ := exec.Command("journalctl", "-p", "warning", "--since", "1 hour ago", "--no-pager").Output()
		info.RecentWarnings = parseCount(string(warnOut))

		if info.RecentErrors > 0 {
			detailOut, _ := exec.Command("journalctl", "-p", "err", "--since", "1 hour ago", "--no-pager", "-n", "5", "--output=short").Output()
			info.RecentDetails = strings.Split(strings.TrimRight(string(detailOut), "\n"), "\n")
		}
	}

	dmesgOut, err := exec.Command("dmesg", "-l", "err,crit,alert,emerg").Output()
	if err == nil {
		lines := strings.Split(strings.TrimRight(string(dmesgOut), "\n"), "\n")
		tailCount := 5
		if len(lines) < tailCount {
			tailCount = len(lines)
		}
		info.KernelErrors = parseCount(strings.Join(lines[len(lines)-tailCount:], "\n"))
	}
	if info.KernelErrors > 0 {
		info.KernelStatus = StatusWarning
	} else {
		info.KernelStatus = StatusGood
	}

	return info
}
