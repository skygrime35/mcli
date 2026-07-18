// internal/hotspot/interface_test.go
package hotspot

import "testing"

const sampleDeviceStatus = `DEVICE           TYPE      STATE                   CONNECTION
wlo1             wifi      connected               Le Super Telephone
br-641d3af8cb9d  bridge    connected (externally)  br-641d3af8cb9d
lo               loopback  connected (externally)  lo
docker0          bridge    connected (externally)  docker0
p2p-dev-wlo1     wifi-p2p  disconnected            --
enp3s0           ethernet  unavailable             --
`

func TestParseWifiInterface(t *testing.T) {
	if got := parseWifiInterface(sampleDeviceStatus); got != "wlo1" {
		t.Errorf("parseWifiInterface() = %q, want %q", got, "wlo1")
	}
}

func TestParseWifiInterface_NoWifi(t *testing.T) {
	noWifi := "DEVICE   TYPE      STATE                   CONNECTION\nlo       loopback  connected (externally)  lo\n"
	if got := parseWifiInterface(noWifi); got != "" {
		t.Errorf("expected empty string when no wifi device present, got %q", got)
	}
}
