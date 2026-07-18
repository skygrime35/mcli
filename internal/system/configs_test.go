// internal/system/configs_test.go
package system

import "testing"

const sampleDpkgFullList = `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
||/ Name              Version      Architecture Description
+++-=================-============-============-=================
ii  curl               8.5.0-2ubuntu1  amd64  command line tool for transferring data
rc  old-package-one    1.2.3-1      amd64  some old removed package
rc  old-package-two    4.5.6-2      amd64  another removed package
ii  vim                2:9.1.0016   amd64  Vi IMproved
`

func TestParseOrphanedConfigs(t *testing.T) {
	got := parseOrphanedConfigs(sampleDpkgFullList)
	want := []string{"old-package-one", "old-package-two"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseOrphanedConfigs_None(t *testing.T) {
	if got := parseOrphanedConfigs("ii  curl  8.5.0  amd64  desc\n"); len(got) != 0 {
		t.Errorf("expected no orphaned configs, got %v", got)
	}
}
