// internal/network/interfaces.go
package network

import "net"

// parseIPv4Addrs extracts dotted-decimal IPv4 strings from a slice of
// net.Addr (as returned by net.Interface.Addrs()), skipping IPv6 and any
// address that doesn't carry an IP (defensive - all real net.Addr
// implementations here are *net.IPNet).
func parseIPv4Addrs(addrs []net.Addr) []string {
	var ips []string
	for _, a := range addrs {
		ipNet, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		v4 := ipNet.IP.To4()
		if v4 == nil {
			continue
		}
		ips = append(ips, v4.String())
	}
	return ips
}

// Interfaces reports every network interface with at least one IPv4
// address, skipping the loopback interface (matching the old reference's
// behavior of skipping "lo").
func Interfaces() ([]InterfaceInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []InterfaceInfo
	for _, iface := range ifaces {
		if iface.Name == "lo" {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		ips := parseIPv4Addrs(addrs)
		if len(ips) == 0 {
			continue
		}
		result = append(result, InterfaceInfo{Name: iface.Name, IPs: ips})
	}
	return result, nil
}
