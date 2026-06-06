package collector

import (
	"github.com/shirou/gopsutil/v3/mem"
)

// CollectMemory gathers memory and swap usage metrics.
func CollectMemory() (*MemoryMetrics, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryMetrics{
		Total:       vmem.Total,
		Available:   vmem.Available,
		Used:        vmem.Used,
		UsedPercent: vmem.UsedPercent,
		Free:        vmem.Free,
		Cached:      vmem.Cached,
		Buffers:     vmem.Buffers,
		SwapTotal:   swap.Total,
		SwapUsed:    swap.Used,
		SwapFree:    swap.Free,
	}, nil
}
