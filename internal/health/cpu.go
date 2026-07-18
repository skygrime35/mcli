// internal/health/cpu.go
package health

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/skygrime35/mcli/internal/platform"
)

func parseCPUModel(cpuinfo string) string {
	for _, line := range strings.Split(cpuinfo, "\n") {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "Unknown"
}

func parseLoadAvg(content string) (load1, load5, load15 float64, err error) {
	fields := strings.Fields(content)
	if len(fields) < 3 {
		return 0, 0, 0, fmt.Errorf("malformed /proc/loadavg content: %q", content)
	}
	load1, err = strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parsing load1: %w", err)
	}
	load5, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parsing load5: %w", err)
	}
	load15, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parsing load15: %w", err)
	}
	return load1, load5, load15, nil
}

// cpuStatSample holds the fields of /proc/stat's aggregate "cpu " line
// needed to compute a usage percentage from two samples over time.
type cpuStatSample struct {
	idle, total uint64
}

func readCPUStatSample() (cpuStatSample, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return cpuStatSample{}, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)[1:]
		var total, idle uint64
		for i, f := range fields {
			v, err := strconv.ParseUint(f, 10, 64)
			if err != nil {
				continue
			}
			total += v
			if i == 3 { // idle is the 4th field (index 3) of "cpu  user nice system idle ..."
				idle = v
			}
		}
		return cpuStatSample{idle: idle, total: total}, nil
	}
	return cpuStatSample{}, fmt.Errorf("no \"cpu \" line found in /proc/stat")
}

// measureCPUUsagePercent takes two /proc/stat samples ~300ms apart and
// computes the non-idle fraction over that window. This is a real delta
// measurement, more reliable than the bash reference's single top -bn1
// snapshot (see the plan's Global Constraints for why).
func measureCPUUsagePercent() (float64, error) {
	first, err := readCPUStatSample()
	if err != nil {
		return 0, err
	}
	time.Sleep(300 * time.Millisecond)
	second, err := readCPUStatSample()
	if err != nil {
		return 0, err
	}
	totalDelta := second.total - first.total
	idleDelta := second.idle - first.idle
	if totalDelta == 0 {
		return 0, nil
	}
	usage := float64(totalDelta-idleDelta) / float64(totalDelta) * 100
	return usage, nil
}

// thermalZoneTempRegexp matches a signed temperature reading in `sensors`
// output, e.g. "+43.0" or "-5.0" out of "Core 0:        +43.0°C  (...)".
// The sign is REQUIRED (not "\+?") specifically so this can't match the bare
// label digit that precedes the real reading (e.g. the "0" in "Core 0:"),
// which has no sign of its own. The captured group includes the sign so
// strconv.ParseFloat can parse genuinely sub-zero readings correctly too.
var thermalZoneTempRegexp = regexp.MustCompile(`([-+][0-9]+\.?[0-9]*)`)

// thermalZoneTempPath is the sysfs path read as the primary CPU temperature
// source. It's a package-level var (rather than a literal) so tests can
// point it at a fixture to force the "sensors" fallback path below.
var thermalZoneTempPath = "/sys/class/thermal/thermal_zone0/temp"

// parseSensorsTempLine extracts the signed temperature reading from a single
// line of `sensors` command output, e.g.
// "Core 0:        +43.0°C  (high = +80.0°C, crit = +90.0°C)" -> 43.0, true.
// It's a small pure helper (rather than inline logic in readCPUTemp) so the
// extraction can be unit-tested in isolation without shelling out.
func parseSensorsTempLine(line string) (celsius float64, ok bool) {
	m := thermalZoneTempRegexp.FindStringSubmatch(line)
	if len(m) != 2 {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func readCPUTemp(caps platform.Capabilities) (celsius float64, available bool) {
	if data, err := os.ReadFile(thermalZoneTempPath); err == nil {
		millideg, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err == nil {
			return millideg / 1000, true
		}
	}
	if caps.Sensors {
		out, err := exec.Command("sensors").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "Core 0") || strings.Contains(line, "Package id") {
					if v, ok := parseSensorsTempLine(line); ok {
						return v, true
					}
					break
				}
			}
		}
	}
	return 0, false
}

func readCPUFreqMHz() (mhz float64, available bool) {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		return 0, false
	}
	khz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0, false
	}
	return khz / 1000, true
}

func CollectCPU(caps platform.Capabilities) CPUInfo {
	info := CPUInfo{Model: "Unknown", Cores: runtime.NumCPU()}

	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		info.Model = parseCPUModel(string(data))
	}

	if usage, err := measureCPUUsagePercent(); err == nil {
		info.UsagePercent = usage
		info.UsageStatus = statusForHigh(usage, cpuWarningPercent, cpuCriticalPercent)
	}

	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		if l1, l5, l15, err := parseLoadAvg(string(data)); err == nil {
			info.Load1, info.Load5, info.Load15 = l1, l5, l15
		}
	}

	if mhz, ok := readCPUFreqMHz(); ok {
		info.FreqMHz, info.FreqAvailable = mhz, true
	}

	if temp, ok := readCPUTemp(caps); ok {
		info.TempCelsius, info.TempAvailable = temp, true
		info.TempStatus = statusForHigh(temp, tempWarningCelsius, tempCriticalCelsius)
	}

	return info
}
