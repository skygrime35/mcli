// internal/health/thresholds_test.go
package health

import "testing"

func TestStatusForHigh(t *testing.T) {
	cases := []struct {
		value, warn, crit float64
		want              Status
	}{
		{50, 80, 90, StatusGood},
		{85, 80, 90, StatusWarning},
		{95, 80, 90, StatusCritical},
		{90, 80, 90, StatusCritical}, // exactly at critical
		{80, 80, 90, StatusWarning},  // exactly at warning
	}
	for _, c := range cases {
		if got := statusForHigh(c.value, c.warn, c.crit); got != c.want {
			t.Errorf("statusForHigh(%v, %v, %v) = %v, want %v", c.value, c.warn, c.crit, got, c.want)
		}
	}
}

func TestStatusForLow(t *testing.T) {
	cases := []struct {
		value, warn, crit float64
		want              Status
	}{
		{90, 50, 20, StatusGood},
		{40, 50, 20, StatusWarning},
		{10, 50, 20, StatusCritical},
		{20, 50, 20, StatusCritical},
		{50, 50, 20, StatusWarning},
	}
	for _, c := range cases {
		if got := statusForLow(c.value, c.warn, c.crit); got != c.want {
			t.Errorf("statusForLow(%v, %v, %v) = %v, want %v", c.value, c.warn, c.crit, got, c.want)
		}
	}
}

func TestThresholdConstants_MatchBashReference(t *testing.T) {
	// One value from each independent threshold set in pc_health.sh, to
	// catch accidental unification/typos rather than re-deriving every
	// number (the full set is exercised implicitly via collector tests
	// in later tasks).
	if cpuWarningPercent != 80 || cpuCriticalPercent != 90 {
		t.Errorf("cpu thresholds: got %v/%v, want 80/90", cpuWarningPercent, cpuCriticalPercent)
	}
	if swapWarningPercent != 50 || swapCriticalPercent != 80 {
		t.Errorf("swap thresholds: got %v/%v, want 50/80", swapWarningPercent, swapCriticalPercent)
	}
	if batteryChargeWarningPercent != 50 || batteryChargeCriticalPercent != 20 {
		t.Errorf("battery charge thresholds: got %v/%v, want 50/20", batteryChargeWarningPercent, batteryChargeCriticalPercent)
	}
	if batteryHealthWarningPercent != 80 || batteryHealthCriticalPercent != 50 {
		t.Errorf("battery health thresholds: got %v/%v, want 80/50", batteryHealthWarningPercent, batteryHealthCriticalPercent)
	}
	if batteryCyclesWarning != 500 || batteryCyclesCritical != 1000 {
		t.Errorf("battery cycle thresholds: got %v/%v, want 500/1000", batteryCyclesWarning, batteryCyclesCritical)
	}
	if journalErrorsWarning != 10 || journalErrorsCritical != 50 {
		t.Errorf("journal error thresholds: got %v/%v, want 10/50", journalErrorsWarning, journalErrorsCritical)
	}
	if securityUpdatesWarning != 0 || securityUpdatesCritical != 10 {
		t.Errorf("security update thresholds: got %v/%v, want 0/10", securityUpdatesWarning, securityUpdatesCritical)
	}
}
