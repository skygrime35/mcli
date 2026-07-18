// internal/hotspot/stats_test.go
package hotspot

import "testing"

func TestFindActiveHotspotInterface(t *testing.T) {
	output := "wlo1:wifi:connected:MyHotspot\ndocker0:bridge:connected:docker0\n"
	if got := findActiveHotspotInterface(output, "MyHotspot"); got != "wlo1" {
		t.Errorf("got %q, want %q", got, "wlo1")
	}
}

func TestFindActiveHotspotInterface_NoMatch(t *testing.T) {
	output := "docker0:bridge:connected:docker0\n"
	if got := findActiveHotspotInterface(output, "MyHotspot"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestParseNeighClients(t *testing.T) {
	output := `192.168.1.5 dev wlo1 lladdr aa:bb:cc:dd:ee:ff REACHABLE
192.168.1.6 dev wlo1 lladdr 11:22:33:44:55:66 STALE
192.168.1.7 dev wlo1  FAILED
192.168.1.8 dev wlo1 lladdr 77:88:99:aa:bb:cc DELAY
`
	got := parseNeighClients(output)
	want := []string{"192.168.1.5", "192.168.1.6", "192.168.1.8"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v (FAILED entries must be excluded)", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseInterfaceBytes(t *testing.T) {
	procNetDev := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo:  123456     100    0    0    0     0          0         0   123456     100    0    0    0     0       0          0
  wlo1: 987654321  654321    0    0    0     0          0         0 123456789   98765    0    0    0     0       0          0
`
	rx, tx, ok := parseInterfaceBytes(procNetDev, "wlo1")
	if !ok {
		t.Fatal("expected ok=true for an interface present in the fixture")
	}
	if rx != 987654321 || tx != 123456789 {
		t.Errorf("got rx=%d tx=%d, want rx=987654321 tx=123456789", rx, tx)
	}
}

func TestParseInterfaceBytes_NotFound(t *testing.T) {
	if _, _, ok := parseInterfaceBytes("lo: 1 2 3 4 5 6 7 8 9\n", "wlo1"); ok {
		t.Error("expected ok=false when the interface isn't present")
	}
}

func TestParseInterfaceBytes_SubstringCollision(t *testing.T) {
	// Regression test: eth0 must not match veth0.
	// This tests the substring-collision bug fix.
	procNetDev := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo:  111111     100    0    0    0     0          0         0   111111     100    0    0    0     0       0          0
  veth0: 222222  200000    0    0    0     0          0         0 333333   200000    0    0    0     0       0          0
  eth0: 444444  300000    0    0    0     0          0         0 555555   300000    0    0    0     0       0          0
`
	// Looking up "eth0" should return eth0's counts (444444 rx, 555555 tx),
	// NOT veth0's (222222 rx, 333333 tx).
	rx, tx, ok := parseInterfaceBytes(procNetDev, "eth0")
	if !ok {
		t.Fatal("expected ok=true for eth0")
	}
	if rx != 444444 || tx != 555555 {
		t.Errorf("got rx=%d tx=%d, want rx=444444 tx=555555 (eth0's counts, not veth0's)", rx, tx)
	}

	// Also verify veth0 returns its own distinct counts.
	rx, tx, ok = parseInterfaceBytes(procNetDev, "veth0")
	if !ok {
		t.Fatal("expected ok=true for veth0")
	}
	if rx != 222222 || tx != 333333 {
		t.Errorf("got rx=%d tx=%d, want rx=222222 tx=333333 (veth0's counts)", rx, tx)
	}
}

func TestFormatBytes(t *testing.T) {
	cases := map[uint64]string{
		500:            "500.0B",
		2048:           "2.0KiB",
		5 * 1024 * 1024: "5.0MiB",
	}
	for input, want := range cases {
		if got := FormatBytes(input); got != want {
			t.Errorf("FormatBytes(%d) = %q, want %q", input, got, want)
		}
	}
}
