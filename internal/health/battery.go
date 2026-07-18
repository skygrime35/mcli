// internal/health/battery.go
package health

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/skygrime35/mcli/internal/platform"
)

func readIntFile(path string) (int, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	v, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, false
	}
	return v, true
}

func readFloatFile(path string) (float64, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func readStringFile(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(data)), true
}

func batteryConditionLabel(healthPercent float64) string {
	switch {
	case healthPercent >= 80:
		return "Excellent"
	case healthPercent >= 60:
		return "Good"
	case healthPercent >= 40:
		return "Fair"
	default:
		return "Poor"
	}
}

func CollectBattery(caps platform.Capabilities) *BatteryInfo {
	base := "/sys/class/power_supply/BAT0"
	if _, err := os.Stat("/sys/class/power_supply/BAT1"); err == nil {
		base = "/sys/class/power_supply/BAT1"
	} else if _, err := os.Stat(base); err != nil {
		return nil
	}

	info := &BatteryInfo{Present: true}

	if capacity, ok := readIntFile(base + "/capacity"); ok {
		info.CapacityPercent = capacity
		info.ChargeStatus = statusForLow(float64(capacity), batteryChargeWarningPercent, batteryChargeCriticalPercent)
	}
	if state, ok := readStringFile(base + "/status"); ok {
		info.State = state
	}

	energyFull, energyOK := readFloatFile(base + "/energy_full")
	energyFullDesign, energyDesignOK := readFloatFile(base + "/energy_full_design")
	if !energyOK || !energyDesignOK {
		energyFull, energyOK = readFloatFile(base + "/charge_full")
		energyFullDesign, energyDesignOK = readFloatFile(base + "/charge_full_design")
	}
	if energyOK && energyDesignOK && energyFullDesign > 0 {
		healthPercent := energyFull * 100 / energyFullDesign
		info.HealthPercent = healthPercent
		info.HealthAvailable = true
		info.HealthStatus = statusForLow(healthPercent, batteryHealthWarningPercent, batteryHealthCriticalPercent)
		info.Condition = batteryConditionLabel(healthPercent)
	}

	if cycles, ok := readIntFile(base + "/cycle_count"); ok && cycles != 0 {
		info.CycleCount = cycles
		info.CycleStatus = statusForHigh(float64(cycles), batteryCyclesWarning, batteryCyclesCritical)
	}

	if voltage, ok := readFloatFile(base + "/voltage_now"); ok {
		info.VoltageVolts = voltage / 1000000
	}
	if power, ok := readFloatFile(base + "/power_now"); ok {
		info.PowerWatts = power / 1000000
	}
	if tech, ok := readStringFile(base + "/technology"); ok {
		info.Technology = tech
	}
	if mfr, ok := readStringFile(base + "/manufacturer"); ok {
		info.Manufacturer = mfr
	}
	if model, ok := readStringFile(base + "/model_name"); ok {
		info.Model = model
	}

	if caps.Upower {
		info.TimeRemaining = batteryTimeRemaining(info.State)
	}

	return info
}

func batteryTimeRemaining(state string) string {
	devOut, err := exec.Command("upower", "-e").Output()
	if err != nil {
		return ""
	}
	var device string
	for _, line := range strings.Split(string(devOut), "\n") {
		if strings.Contains(line, "BAT") {
			device = strings.TrimSpace(line)
			break
		}
	}
	if device == "" {
		return ""
	}

	field := "time to empty"
	if state == "Charging" {
		field = "time to full"
	}
	infoOut, err := exec.Command("upower", "-i", device).Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(infoOut), "\n") {
		if strings.Contains(line, field) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return strings.Join(fields[len(fields)-2:], " ")
			}
		}
	}
	return ""
}
