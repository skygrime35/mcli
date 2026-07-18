package health

import (
	"testing"

	"github.com/skygrime35/mcli/internal/platform"
)

func TestCollect_PopulatesTimestampAndCoreSections(t *testing.T) {
	// This test runs against the REAL machine (reads /proc, /sys, execs
	// real tools per whatever capabilities are actually detected) since
	// Collect is an orchestration function, not a pure parser - this
	// mirrors how the foundation plan treated its own top-level
	// collection functions. It only asserts structural invariants that
	// hold on any machine, not specific values.
	caps := platform.Detect()
	snap := Collect(caps)

	if snap.Timestamp.IsZero() {
		t.Error("expected a non-zero Timestamp")
	}
	if snap.System.Hostname == "" {
		t.Error("expected a non-empty hostname")
	}
	if snap.CPU.Cores <= 0 {
		t.Error("expected a positive core count")
	}
	if snap.Memory.TotalKB == 0 {
		t.Error("expected a non-zero total memory reading")
	}
	// Disks/Network/Services/Errors/Security are allowed to be empty
	// slices/zero values on a minimal or capability-poor machine (e.g. a
	// CI container without systemd) - only assert they don't panic, which
	// this test already does implicitly by reaching this point.
}
