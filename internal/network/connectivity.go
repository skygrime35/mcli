// internal/network/connectivity.go
package network

import (
	"net"
	"time"
)

// CheckConnectivity reports whether this machine can reach the public
// internet, by attempting a TCP handshake against a well-known DNS
// resolver (8.8.8.8:53) - matching the old Python reference's probe
// target. This never sends application data, just completes or fails a
// TCP connect.
func CheckConnectivity() bool {
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
