// internal/health/types.go
package health

import "time"

type Status string

const (
	StatusGood     Status = "good"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
	StatusInfo     Status = "info"
)

type Snapshot struct {
	Timestamp time.Time    `json:"timestamp"`
	System    SystemInfo   `json:"system"`
	CPU       CPUInfo      `json:"cpu"`
	Memory    MemoryInfo   `json:"memory"`
	Disks     []DiskInfo   `json:"disks"`
	SMART     []SMARTInfo  `json:"smart,omitempty"`
	Network   NetworkInfo  `json:"network"`
	Services  ServicesInfo `json:"services"`
	Battery   *BatteryInfo `json:"battery,omitempty"`
	Errors    ErrorsInfo   `json:"errors"`
	Security  SecurityInfo `json:"security"`
}

type SystemInfo struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Kernel       string `json:"kernel"`
	Architecture string `json:"architecture"`
	Uptime       string `json:"uptime"`
	LastBoot     string `json:"last_boot"`
	LoggedUsers  int    `json:"logged_users"`
}

type CPUInfo struct {
	Model             string  `json:"model"`
	Cores             int     `json:"cores"`
	UsagePercent      float64 `json:"usage_percent"`
	UsageStatus       Status  `json:"usage_status"`
	Load1             float64 `json:"load1"`
	Load5             float64 `json:"load5"`
	Load15            float64 `json:"load15"`
	FreqMHz           float64 `json:"freq_mhz,omitempty"`
	FreqAvailable     bool    `json:"-"`
	TempCelsius       float64 `json:"temp_celsius,omitempty"`
	TempStatus        Status  `json:"temp_status,omitempty"`
	TempAvailable     bool    `json:"-"`
}

type MemoryInfo struct {
	TotalKB        uint64  `json:"total_kb"`
	UsedKB         uint64  `json:"used_kb"`
	UsagePercent   float64 `json:"usage_percent"`
	UsageStatus    Status  `json:"usage_status"`
	CachedKB       uint64  `json:"cached_kb"`
	BuffersKB      uint64  `json:"buffers_kb"`
	SwapTotalKB    uint64  `json:"swap_total_kb"`
	SwapUsedKB     uint64  `json:"swap_used_kb"`
	SwapPercent    float64 `json:"swap_percent"`
	SwapStatus     Status  `json:"swap_status"`
	SwapConfigured bool    `json:"swap_configured"`
}

type DiskInfo struct {
	Filesystem   string `json:"filesystem"`
	Mount        string `json:"mount"`
	SizeHuman    string `json:"size"`
	UsedHuman    string `json:"used"`
	AvailHuman   string `json:"available"`
	UsagePercent int    `json:"usage_percent"`
	Status       Status `json:"status"`
}

type SMARTInfo struct {
	Device  string `json:"device"`
	Healthy bool   `json:"healthy"`
	Status  Status `json:"status"`
	Detail  string `json:"detail"`
}

type NetworkInterface struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	State     string `json:"state"`
	Primary   bool   `json:"primary"`
	MAC       string `json:"mac,omitempty"`
	Status    Status `json:"status"`
}

type NetworkInfo struct {
	Interfaces      []NetworkInterface `json:"interfaces"`
	Gateway         string             `json:"gateway"`
	InternetOK      bool               `json:"internet_ok"`
	InternetStatus  Status             `json:"internet_status"`
	DNSOK           bool               `json:"dns_ok"`
	DNSStatus       Status             `json:"dns_status"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	State  string `json:"state"` // "running" | "stopped_enabled" | "stopped_disabled"
	Status Status `json:"status"`
}

type ServicesInfo struct {
	Services      []ServiceStatus `json:"services"`
	FailedCount   int             `json:"failed_count"`
	FailedNames   []string        `json:"failed_names,omitempty"`
	FailedStatus  Status          `json:"failed_status"`
}

type BatteryInfo struct {
	Present          bool    `json:"present"`
	CapacityPercent  int     `json:"capacity_percent"`
	ChargeStatus     Status  `json:"charge_status"`
	State            string  `json:"state"` // Charging/Discharging/etc
	HealthPercent    float64 `json:"health_percent,omitempty"`
	HealthAvailable  bool    `json:"-"`
	HealthStatus     Status  `json:"health_status,omitempty"`
	Condition        string  `json:"condition,omitempty"`
	CycleCount       int     `json:"cycle_count,omitempty"`
	CycleStatus      Status  `json:"cycle_status,omitempty"`
	VoltageVolts     float64 `json:"voltage_volts,omitempty"`
	PowerWatts       float64 `json:"power_watts,omitempty"`
	TimeRemaining    string  `json:"time_remaining,omitempty"`
	Technology       string  `json:"technology,omitempty"`
	Manufacturer     string  `json:"manufacturer,omitempty"`
	Model            string  `json:"model,omitempty"`
}

type ErrorsInfo struct {
	RecentErrors     int      `json:"recent_errors"`
	ErrorsStatus     Status   `json:"errors_status"`
	RecentWarnings   int      `json:"recent_warnings"`
	RecentDetails    []string `json:"recent_details,omitempty"`
	KernelErrors     int      `json:"kernel_errors"`
	KernelStatus     Status   `json:"kernel_status"`
}

type SecurityInfo struct {
	UpdatesAvailable   int    `json:"updates_available"`
	SecurityUpdates    int    `json:"security_updates"`
	SecurityStatus     Status `json:"security_status"`
	UpdatesChecked     bool   `json:"updates_checked"`
	FirewallState      string `json:"firewall_state"` // "active" | "inactive" | "rules:N" | "unknown"
	FirewallStatus     Status `json:"firewall_status"`
	FailedSSHAttempts  int    `json:"failed_ssh_attempts,omitempty"`
}
