package share

import "net"

// LocalIP returns this machine's outbound-facing local IP address, for
// display purposes (so the user can share "http://<ip>:<port>" with
// someone else on the LAN). It opens a UDP "connection" to a
// non-routable address purely to ask the OS which local interface would
// be used - no packet is actually sent. Falls back to "127.0.0.1" if
// this fails for any reason (matching the old Python reference's
// fallback).
func LocalIP() string {
	conn, err := net.Dial("udp", "10.255.255.255:1")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "127.0.0.1"
	}
	return addr.IP.String()
}
