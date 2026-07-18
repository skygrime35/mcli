// internal/health/memory_test.go
package health

import "testing"

const sampleMemInfo = `MemTotal:       16384000 kB
MemFree:         2048000 kB
MemAvailable:    8192000 kB
Buffers:          512000 kB
Cached:          4096000 kB
SwapCached:            0 kB
SwapTotal:       4096000 kB
SwapFree:        3072000 kB
`

func TestParseMemInfo(t *testing.T) {
	info := parseMemInfo(sampleMemInfo)

	if info.TotalKB != 16384000 {
		t.Errorf("TotalKB = %d, want 16384000", info.TotalKB)
	}
	wantUsed := uint64(16384000 - 8192000)
	if info.UsedKB != wantUsed {
		t.Errorf("UsedKB = %d, want %d", info.UsedKB, wantUsed)
	}
	wantPercent := float64(wantUsed) * 100 / 16384000
	if info.UsagePercent != wantPercent {
		t.Errorf("UsagePercent = %v, want %v", info.UsagePercent, wantPercent)
	}
	if info.CachedKB != 4096000 {
		t.Errorf("CachedKB = %d, want 4096000 (must match ^Cached:, not SwapCached)", info.CachedKB)
	}
	if info.BuffersKB != 512000 {
		t.Errorf("BuffersKB = %d, want 512000", info.BuffersKB)
	}
	if !info.SwapConfigured {
		t.Error("expected SwapConfigured=true")
	}
	wantSwapUsed := uint64(4096000 - 3072000)
	if info.SwapUsedKB != wantSwapUsed {
		t.Errorf("SwapUsedKB = %d, want %d", info.SwapUsedKB, wantSwapUsed)
	}
	wantSwapPercent := float64(wantSwapUsed) * 100 / 4096000
	if info.SwapPercent != wantSwapPercent {
		t.Errorf("SwapPercent = %v, want %v", info.SwapPercent, wantSwapPercent)
	}
	if info.SwapStatus != StatusGood {
		// wantSwapPercent here is 25%, below the 50% swap-warning threshold
		t.Errorf("SwapStatus = %v, want %v", info.SwapStatus, StatusGood)
	}
}

func TestParseMemInfo_NoSwap(t *testing.T) {
	content := "MemTotal:       16384000 kB\nMemAvailable:    8192000 kB\nSwapTotal:             0 kB\nSwapFree:              0 kB\n"
	info := parseMemInfo(content)
	if info.SwapConfigured {
		t.Error("expected SwapConfigured=false when SwapTotal is 0")
	}
	if info.SwapPercent != 0 {
		t.Errorf("expected SwapPercent=0 when swap is not configured, got %v", info.SwapPercent)
	}
}
