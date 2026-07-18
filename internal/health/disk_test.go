// internal/health/disk_test.go
package health

import "testing"

const sampleDFOutput = `Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        98G   45G   48G  49% /
/dev/sda2       450G  200G  227G  47% /home
tmpfs           7.8G     0  7.8G   0% /dev/shm
`

func TestParseDF(t *testing.T) {
	disks := parseDF(sampleDFOutput)
	if len(disks) != 2 {
		t.Fatalf("expected 2 disks (tmpfs already excluded by df flags upstream), got %d", len(disks))
	}
	if disks[0].Filesystem != "/dev/sda1" || disks[0].Mount != "/" {
		t.Errorf("disk 0 = %+v, want filesystem /dev/sda1 mount /", disks[0])
	}
	if disks[0].UsagePercent != 49 {
		t.Errorf("disk 0 usage = %d, want 49", disks[0].UsagePercent)
	}
	if disks[0].Status != StatusGood {
		// 49% is below the 80 warning threshold, so status should be StatusGood.
		t.Errorf("disk 0 status = %v", disks[0].Status)
	}
	if disks[1].UsagePercent != 47 {
		t.Errorf("disk 1 usage = %d, want 47", disks[1].UsagePercent)
	}
}
