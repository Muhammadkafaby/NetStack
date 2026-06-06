package collector

import (
	"github.com/shirou/gopsutil/v3/net"
)

// CollectNetwork gathers network I/O metrics.
func CollectNetwork() (*NetworkMetrics, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	metrics := &NetworkMetrics{}

	var interfaces []InterfaceMetrics

	for _, c := range counters {
		// Skip loopback
		if c.Name == "lo" {
			continue
		}

		metrics.BytesSent += c.BytesSent
		metrics.BytesRecv += c.BytesRecv
		metrics.PacketsSent += c.PacketsSent
		metrics.PacketsRecv += c.PacketsRecv
		metrics.ErrorsIn += c.Errin
		metrics.ErrorsOut += c.Errout
		metrics.DropsIn += c.Dropin
		metrics.DropsOut += c.Dropout

		interfaces = append(interfaces, InterfaceMetrics{
			Name:      c.Name,
			BytesSent: c.BytesSent,
			BytesRecv: c.BytesRecv,
			ErrorsIn:  c.Errin,
			ErrorsOut: c.Errout,
		})
	}

	metrics.Interfaces = interfaces

	return metrics, nil
}
