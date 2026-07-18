// internal/health/smart.go
package health

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/skygrime35/mcli/internal/platform"
)

type smartctlJSON struct {
	SmartStatus struct {
		Passed bool `json:"passed"`
	} `json:"smart_status"`
}

func parseSmartctlJSON(data []byte) (healthy bool, ok bool) {
	var v smartctlJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return false, false
	}
	return v.SmartStatus.Passed, true
}

var smartHealthLineRegexp = regexp.MustCompile(`(?i)SMART overall-health self-assessment test result:\s*(\S+)`)

// nvmeNamespaceRegexp matches a whole NVMe namespace device path such as
// /dev/nvme0n1, but NOT a partition of that namespace such as /dev/nvme0n1p1.
var nvmeNamespaceRegexp = regexp.MustCompile(`^/dev/nvme\d+n\d+$`)

// isNVMeNamespace reports whether path is a whole NVMe namespace device
// (e.g. /dev/nvme0n1) as opposed to one of its partitions (e.g.
// /dev/nvme0n1p1). SMART health is a per-physical-disk concept, so only
// namespace devices should be probed - partitions matched by the
// /dev/nvme*n* glob must be filtered out.
func isNVMeNamespace(path string) bool {
	return nvmeNamespaceRegexp.MatchString(path)
}

func parseSmartctlText(output string) (healthy bool, ok bool) {
	m := smartHealthLineRegexp.FindStringSubmatch(output)
	if len(m) != 2 {
		return false, false
	}
	result := strings.ToUpper(strings.TrimSuffix(m[1], "!"))
	return result == "PASSED", true
}

// smartDevices returns candidate block device paths to probe: SATA/SCSI
// (/dev/sd*) as covered by the bash reference, PLUS NVMe (/dev/nvme*n*),
// which the bash reference's `/dev/sd?` glob misses entirely - a
// deliberate improvement over the bash script (see plan's Global
// Constraints).
//
// The /dev/nvme*n* glob matches both whole namespace devices (/dev/nvme0n1)
// and their partitions (/dev/nvme0n1p1, /dev/nvme0n1p2, ...); only the
// former are physical disks with their own SMART health, so
// isNVMeNamespace filters out anything with a trailing partition suffix.
func smartDevices() []string {
	var devices []string
	sdMatches, _ := filepath.Glob("/dev/sd?")
	devices = append(devices, sdMatches...)
	nvmeMatches, _ := filepath.Glob("/dev/nvme*n*")
	for _, m := range nvmeMatches {
		if isNVMeNamespace(m) {
			devices = append(devices, m)
		}
	}
	return devices
}

// CollectSMART probes each candidate device with smartctl and reports its
// SMART health status.
//
// smartctl's process exit code is a bitmask, and bit 3 (value 8) is set
// whenever the SMART overall-health self-assessment result is FAILED - this
// is true for BOTH the "-j" (JSON) invocation and the plain-text fallback
// below. That means a genuinely failing disk is exactly the case where
// exec.Cmd.Output() returns a non-nil *exec.ExitError, even though stdout
// was populated with perfectly valid, parseable output reporting the
// failure. Output() returns whatever was written to stdout regardless of
// the process's exit code; the only ways it fails to do so here are if the
// binary can't be started at all (ruled out by the caps.Smartctl gate
// above) or if the command genuinely wrote nothing to stdout. So the exec
// error is deliberately discarded (`_`) and every attempt below always
// tries to parse whatever bytes came back. Do NOT reintroduce an
// `err == nil` gate before parsing - that would silently downgrade the
// most important case (a disk that is actually critically failing) to
// StatusInfo "N/A (requires root)".
func CollectSMART(caps platform.Capabilities) []SMARTInfo {
	if !caps.Smartctl {
		return nil
	}
	var results []SMARTInfo
	for _, device := range smartDevices() {
		info := SMARTInfo{Device: device}

		jsonOut, _ := exec.Command("smartctl", "-H", "-j", device).Output()
		if healthy, ok := parseSmartctlJSON(jsonOut); ok {
			info.Healthy = healthy
			if healthy {
				info.Status = StatusGood
				info.Detail = "PASSED"
			} else {
				info.Status = StatusCritical
				info.Detail = "FAILED"
			}
			results = append(results, info)
			continue
		}

		textOut, _ := exec.Command("smartctl", "-H", device).Output()
		if healthy, ok := parseSmartctlText(string(textOut)); ok {
			info.Healthy = healthy
			if healthy {
				info.Status = StatusGood
				info.Detail = "PASSED"
			} else {
				info.Status = StatusCritical
				info.Detail = "FAILED"
			}
		} else {
			info.Status = StatusInfo
			info.Detail = "N/A (requires root)"
		}
		results = append(results, info)
	}
	return results
}
