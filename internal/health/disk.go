// internal/health/disk.go
package health

import (
	"os/exec"
	"strconv"
	"strings"
)

// pseudoFilesystems lists virtual filesystems that df may report but that
// don't correspond to real disks. CollectDisks also excludes these via
// `df -x`, but parseDF filters them too so it behaves correctly as a
// standalone pure function regardless of how its input was produced.
var pseudoFilesystems = map[string]bool{
	"tmpfs":    true,
	"devtmpfs": true,
	"squashfs": true,
}

func parseDF(output string) []DiskInfo {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 2 {
		return nil
	}
	var disks []DiskInfo
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		if pseudoFilesystems[fields[0]] {
			continue
		}
		percentStr := strings.TrimSuffix(fields[4], "%")
		percent, err := strconv.Atoi(percentStr)
		if err != nil {
			continue
		}
		disks = append(disks, DiskInfo{
			Filesystem:   fields[0],
			SizeHuman:    fields[1],
			UsedHuman:    fields[2],
			AvailHuman:   fields[3],
			UsagePercent: percent,
			Mount:        fields[5],
			Status:       statusForHigh(float64(percent), diskWarningPercent, diskCriticalPercent),
		})
	}
	return disks
}

func CollectDisks() []DiskInfo {
	out, err := exec.Command("df", "-h",
		"--output=source,size,used,avail,pcent,target",
		"-x", "tmpfs", "-x", "devtmpfs", "-x", "squashfs",
	).Output()
	if err != nil {
		return nil
	}
	return parseDF(string(out))
}
