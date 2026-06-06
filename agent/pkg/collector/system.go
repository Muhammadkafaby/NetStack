package collector

import (
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
)

// CollectSystem gathers system-level metrics: load, uptime, host info.
func CollectSystem(hostname string) (uint64, string, string, *LoadMetrics, error) {
	// Host info
	hostInfo, err := host.Info()
	if err != nil {
		return 0, "", "", nil, err
	}

	// Load averages
	loadAvg, err := load.Avg()
	if err != nil {
		return 0, "", "", nil, err
	}

	loadMetrics := &LoadMetrics{
		Load1:  loadAvg.Load1,
		Load5:  loadAvg.Load5,
		Load15: loadAvg.Load15,
	}

	return hostInfo.Uptime, hostInfo.OS, hostInfo.Platform, loadMetrics, nil
}
