// internal/system/kernel_test.go
package system

import "testing"

func TestParseKernelBase(t *testing.T) {
	cases := map[string]string{
		"5.15.0-91-generic": "5.15.0-91",
		"6.8.0-45-generic":  "6.8.0-45",
		"5.4.0-42-aws":      "5.4.0-42",
		"5.15.0-91":         "5.15.0-91", // no flavor suffix to strip
	}
	for input, want := range cases {
		if got := parseKernelBase(input); got != want {
			t.Errorf("parseKernelBase(%q) = %q, want %q", input, got, want)
		}
	}
}

const sampleDpkgKernelList = `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
||/ Name              Version      Architecture Description
+++-=================-============-============-=================
ii  linux-image-5.15.0-88-generic  5.15.0-88.98  amd64  Signed kernel image
ii  linux-image-5.15.0-91-generic  5.15.0-91.101 amd64  Signed kernel image
rc  linux-image-5.15.0-76-generic  5.15.0-76.83  amd64  Signed kernel image
`

func TestParseInstalledKernels(t *testing.T) {
	got := parseInstalledKernels(sampleDpkgKernelList)
	want := []string{"linux-image-5.15.0-88-generic", "linux-image-5.15.0-91-generic"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v (rc/removed lines must be excluded)", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestOldKernels(t *testing.T) {
	installed := []string{"linux-image-5.15.0-88-generic", "linux-image-5.15.0-91-generic", "linux-image-generic"}
	got := oldKernels(installed, "5.15.0-91")
	want := []string{"linux-image-5.15.0-88-generic"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("got %v, want %v (current kernel and non-versioned meta-package must be excluded)", got, want)
	}
}

func TestClassifyKernelRemoval_Warning(t *testing.T) {
	old := []string{"linux-image-5.15.0-88-generic"}
	if got := classifyKernelRemoval(old, "5.15.0-91-generic"); got != TierWarning {
		t.Errorf("got %v, want TierWarning", got)
	}
}

func TestClassifyKernelRemoval_Unsafe(t *testing.T) {
	// Defensive-net case: the running kernel somehow still ended up in the
	// candidate list (e.g. a filtering edge case) - must be reclassified.
	old := []string{"linux-image-5.15.0-88-generic", "linux-image-5.15.0-91-generic"}
	if got := classifyKernelRemoval(old, "5.15.0-91-generic"); got != TierUnsafe {
		t.Errorf("got %v, want TierUnsafe", got)
	}
}
