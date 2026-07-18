// internal/health/thresholds.go
package health

const (
	cpuWarningPercent  = 80
	cpuCriticalPercent = 90

	tempWarningCelsius  = 70
	tempCriticalCelsius = 85

	memWarningPercent  = 80
	memCriticalPercent = 90

	swapWarningPercent = 50
	swapCriticalPercent = 80

	diskWarningPercent  = 80
	diskCriticalPercent = 90

	batteryChargeWarningPercent  = 50
	batteryChargeCriticalPercent = 20

	batteryHealthWarningPercent  = 80
	batteryHealthCriticalPercent = 50

	batteryCyclesWarning  = 500
	batteryCyclesCritical = 1000

	journalErrorsWarning  = 10
	journalErrorsCritical = 50

	securityUpdatesWarning  = 0
	securityUpdatesCritical = 10
)

// statusForHigh returns the status for a metric where a HIGHER value is
// worse (CPU/memory/disk/temp/swap usage, battery cycle count, journal
// error count, pending security updates).
func statusForHigh(value, warn, crit float64) Status {
	switch {
	case value >= crit:
		return StatusCritical
	case value >= warn:
		return StatusWarning
	default:
		return StatusGood
	}
}

// statusForLow returns the status for a metric where a LOWER value is
// worse (battery charge percent, battery health percent).
func statusForLow(value, warn, crit float64) Status {
	switch {
	case value <= crit:
		return StatusCritical
	case value <= warn:
		return StatusWarning
	default:
		return StatusGood
	}
}
