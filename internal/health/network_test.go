// internal/health/network_test.go
package health

import "testing"

func TestParseInterfaceList(t *testing.T) {
	got := parseInterfaceList("eth0\nlo\nwlan0\n")
	want := []string{"eth0", "wlan0"} // lo is always excluded
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}
