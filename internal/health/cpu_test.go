// internal/health/cpu_test.go
package health

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/skygrime35/mcli/internal/platform"
)

func TestParseCPUModel(t *testing.T) {
	cpuinfo := "processor\t: 0\nvendor_id\t: GenuineIntel\nmodel name\t: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz\ncache size\t: 12288 KB\n"
	got := parseCPUModel(cpuinfo)
	want := "Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz"
	if got != want {
		t.Errorf("parseCPUModel() = %q, want %q", got, want)
	}
}

func TestParseLoadAvg(t *testing.T) {
	load1, load5, load15, err := parseLoadAvg("0.52 0.58 0.59 2/1234 56789\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if load1 != 0.52 || load5 != 0.58 || load15 != 0.59 {
		t.Errorf("got %v/%v/%v, want 0.52/0.58/0.59", load1, load5, load15)
	}
}

func TestParseLoadAvg_Malformed(t *testing.T) {
	if _, _, _, err := parseLoadAvg("garbage"); err == nil {
		t.Fatal("expected an error for malformed /proc/loadavg content")
	}
}

// installFakeSensors puts a fake "sensors" script on PATH that reports a
// known, fixed temperature via a "Core 0" line, and points
// thermalZoneTempPath at a nonexistent file so readCPUTemp is forced past
// its first (sysfs) branch and into the sensors-fallback branch. It
// restores both PATH and thermalZoneTempPath via t.Cleanup.
func installFakeSensors(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("fake sensors script targets linux-style exec behavior")
	}

	dir := t.TempDir()
	script := filepath.Join(dir, "sensors")
	// Uses a "Package id:" line (no digit before the temperature) to keep
	// this test focused purely on the caps.Sensors gating behavior under
	// test. The temperature-extraction behavior on "Core 0:"-style lines
	// (which do have a label digit before the real, signed reading) is
	// covered separately by TestParseSensorsTempLine.
	content := "#!/bin/sh\necho 'Package id:    +42.0°C  (high = +80.0°C, crit = +90.0°C)'\n"
	if err := os.WriteFile(script, []byte(content), 0o755); err != nil {
		t.Fatalf("writing fake sensors script: %v", err)
	}

	origPath := os.Getenv("PATH")
	origThermalPath := thermalZoneTempPath
	t.Cleanup(func() {
		os.Setenv("PATH", origPath)
		thermalZoneTempPath = origThermalPath
	})

	os.Setenv("PATH", dir+string(os.PathListSeparator)+origPath)
	thermalZoneTempPath = filepath.Join(dir, "does-not-exist")
}

func TestReadCPUTemp_UsesSensorsOnlyWhenCapabilityIsSet(t *testing.T) {
	installFakeSensors(t)

	// With caps.Sensors true, readCPUTemp should shell out to the (fake)
	// "sensors" binary and report the value it prints.
	temp, ok := readCPUTemp(platform.Capabilities{Sensors: true})
	if !ok {
		t.Fatal("expected available=true when caps.Sensors is true and a sensors binary is on PATH")
	}
	if temp != 42.0 {
		t.Errorf("got temp=%v, want 42.0 (parsed from fake sensors output)", temp)
	}

	// With caps.Sensors false, readCPUTemp must NOT fall back to invoking
	// "sensors" even though the (fake) binary is present on PATH - this is
	// the exact behavior Finding 3 requires: gate on the detected
	// capability, not on exec.LookPath.
	temp, ok = readCPUTemp(platform.Capabilities{Sensors: false})
	if ok {
		t.Errorf("expected available=false when caps.Sensors is false, got temp=%v", temp)
	}
}

func TestParseSensorsTempLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantTemp float64
		wantOK   bool
	}{
		{
			name:     "Core 0 line: label digit must not be matched instead of the real reading",
			line:     "Core 0:        +43.0°C  (high = +80.0°C, crit = +90.0°C)",
			wantTemp: 43.0,
			wantOK:   true,
		},
		{
			name:     "Package id 0 line",
			line:     "Package id 0:  +45.0°C",
			wantTemp: 45.0,
			wantOK:   true,
		},
		{
			name:     "sub-zero reading keeps its sign",
			line:     "Core 0:        -5.5°C  (high = +80.0°C, crit = +90.0°C)",
			wantTemp: -5.5,
			wantOK:   true,
		},
		{
			name:   "no signed number present",
			line:   "Adapter: ISA adapter",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp, ok := parseSensorsTempLine(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("parseSensorsTempLine(%q) ok = %v, want %v", tt.line, ok, tt.wantOK)
			}
			if ok && temp != tt.wantTemp {
				t.Errorf("parseSensorsTempLine(%q) = %v, want %v", tt.line, temp, tt.wantTemp)
			}
		})
	}
}

func TestCollectCPU_AcceptsCapabilities(t *testing.T) {
	// CollectCPU's signature threads platform.Capabilities through to
	// readCPUTemp; this just confirms the wiring compiles and runs cleanly
	// for both capability states without panicking.
	for _, sensors := range []bool{true, false} {
		info := CollectCPU(platform.Capabilities{Sensors: sensors})
		if info.Model == "" {
			t.Error("expected a non-empty CPU model")
		}
	}
}
