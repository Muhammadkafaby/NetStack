package collector

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CollectCPU gathers CPU usage metrics.
func CollectCPU() (*CPUMetrics, error) {
	// Total CPU usage
	totalPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	// Per-core CPU usage
	perCorePercent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	return &CPUMetrics{
		TotalPercent: totalPercent[0],
		PerCore:      perCorePercent,
		CoreCount:    runtime.NumCPU(),
	}, nil
}
