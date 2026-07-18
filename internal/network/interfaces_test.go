// internal/network/interfaces_test.go
package network

import (
	"net"
	"testing"
)

func TestParseIPv4Addrs(t *testing.T) {
	addrs := []net.Addr{
		&net.IPNet{IP: net.ParseIP("192.168.1.42"), Mask: net.CIDRMask(24, 32)},
		&net.IPNet{IP: net.ParseIP("fe80::1"), Mask: net.CIDRMask(64, 128)},
		&net.IPNet{IP: net.ParseIP("10.0.0.5"), Mask: net.CIDRMask(8, 32)},
	}
	got := parseIPv4Addrs(addrs)
	want := []string{"192.168.1.42", "10.0.0.5"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v (IPv6 addr must be excluded)", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseIPv4Addrs_NoIPv4(t *testing.T) {
	addrs := []net.Addr{
		&net.IPNet{IP: net.ParseIP("fe80::1"), Mask: net.CIDRMask(64, 128)},
	}
	if got := parseIPv4Addrs(addrs); len(got) != 0 {
		t.Errorf("expected no addresses, got %v", got)
	}
}
