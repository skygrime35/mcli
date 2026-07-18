// internal/health/memory.go
package health

import (
	"os"
	"strconv"
	"strings"
)

func parseMemInfoFields(content string) map[string]uint64 {
	fields := make(map[string]uint64)
	for _, line := range strings.Split(content, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		valueField := strings.Fields(parts[1])
		if len(valueField) == 0 {
			continue
		}
		v, err := strconv.ParseUint(valueField[0], 10, 64)
		if err != nil {
			continue
		}
		fields[key] = v
	}
	return fields
}

func parseMemInfo(content string) MemoryInfo {
	fields := parseMemInfoFields(content)

	total := fields["MemTotal"]
	available := fields["MemAvailable"]
	used := uint64(0)
	if total > available {
		used = total - available
	}
	percent := 0.0
	if total > 0 {
		percent = float64(used) * 100 / float64(total)
	}

	swapTotal := fields["SwapTotal"]
	swapFree := fields["SwapFree"]
	swapUsed := uint64(0)
	if swapTotal > swapFree {
		swapUsed = swapTotal - swapFree
	}
	swapPercent := 0.0
	swapConfigured := swapTotal > 0
	if swapConfigured {
		swapPercent = float64(swapUsed) * 100 / float64(swapTotal)
	}

	return MemoryInfo{
		TotalKB:        total,
		UsedKB:         used,
		UsagePercent:   percent,
		UsageStatus:    statusForHigh(percent, memWarningPercent, memCriticalPercent),
		CachedKB:       fields["Cached"],
		BuffersKB:      fields["Buffers"],
		SwapTotalKB:    swapTotal,
		SwapUsedKB:     swapUsed,
		SwapPercent:    swapPercent,
		SwapStatus:     statusForHigh(swapPercent, swapWarningPercent, swapCriticalPercent),
		SwapConfigured: swapConfigured,
	}
}

func CollectMemory() MemoryInfo {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}
	}
	return parseMemInfo(string(data))
}
