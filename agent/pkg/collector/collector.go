// Package collector gathers system metrics using gopsutil.
package collector

// Metrics represents all system metrics collected by the agent.
type Metrics struct {
	Timestamp int64           `json:"timestamp"`
	Hostname  string          `json:"hostname"`
	OS        string          `json:"os"`
	Platform  string          `json:"platform"`
	Uptime    uint64          `json:"uptime"`
	CPU       *CPUMetrics     `json:"cpu,omitempty"`
	Memory    *MemoryMetrics  `json:"memory,omitempty"`
	Disks     []DiskMetrics   `json:"disks,omitempty"`
	Network   *NetworkMetrics `json:"network,omitempty"`
	Load      *LoadMetrics    `json:"load,omitempty"`
}

// CPUMetrics holds CPU usage information.
type CPUMetrics struct {
	TotalPercent float64   `json:"total_percent"`
	PerCore      []float64 `json:"per_core"`
	CoreCount    int       `json:"core_count"`
}

// MemoryMetrics holds memory usage information.
type MemoryMetrics struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Cached      uint64  `json:"cached"`
	Buffers     uint64  `json:"buffers"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapFree    uint64  `json:"swap_free"`
}

// DiskMetrics holds disk usage and I/O for a mount point.
type DiskMetrics struct {
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	ReadBytes   uint64  `json:"read_bytes"`
	WriteBytes  uint64  `json:"write_bytes"`
	ReadCount   uint64  `json:"read_count"`
	WriteCount  uint64  `json:"write_count"`
}

// NetworkMetrics holds network I/O counters.
type NetworkMetrics struct {
	BytesSent   uint64             `json:"bytes_sent"`
	BytesRecv   uint64             `json:"bytes_recv"`
	PacketsSent uint64             `json:"packets_sent"`
	PacketsRecv uint64             `json:"packets_recv"`
	ErrorsIn    uint64             `json:"errors_in"`
	ErrorsOut   uint64             `json:"errors_out"`
	DropsIn     uint64             `json:"drops_in"`
	DropsOut    uint64             `json:"drops_out"`
	Interfaces  []InterfaceMetrics `json:"interfaces,omitempty"`
}

// InterfaceMetrics holds per-interface network stats.
type InterfaceMetrics struct {
	Name      string `json:"name"`
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
	ErrorsIn  uint64 `json:"errors_in"`
	ErrorsOut uint64 `json:"errors_out"`
}

// LoadMetrics holds system load averages.
type LoadMetrics struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}
