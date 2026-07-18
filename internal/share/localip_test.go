package share

import "testing"

func TestLocalIP_ReturnsNonEmpty(t *testing.T) {
	// This is a real, safe, read-only check (opens a UDP socket, never
	// sends a packet) - it must always return SOME usable address.
	ip := LocalIP()
	if ip == "" {
		t.Error("expected LocalIP to return a non-empty address (falling back to 127.0.0.1 at worst)")
	}
}
