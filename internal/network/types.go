// internal/network/types.go
package network

// InterfaceInfo is one network interface's name and its IPv4 addresses.
type InterfaceInfo struct {
	Name string
	IPs  []string
}

// SpeedTestResult reports a completed download/upload speed measurement.
type SpeedTestResult struct {
	DownloadMbps float64
	UploadMbps   float64
	PingMs       float64
	Server       string
	Country      string
}
