package collector

import (
	"time"
)

// CollectAll gathers all system metrics in a single call.
func CollectAll(hostname string) (*Metrics, error) {
	now := time.Now().Unix()

	// Collect all metrics in parallel-ish sequence
	uptime, os, platform, load, err := CollectSystem(hostname)
	if err != nil {
		return nil, err
	}

	cpu, err := CollectCPU()
	if err != nil {
		return nil, err
	}

	memory, err := CollectMemory()
	if err != nil {
		return nil, err
	}

	disks, err := CollectDisk()
	if err != nil {
		return nil, err
	}

	network, err := CollectNetwork()
	if err != nil {
		return nil, err
	}

	return &Metrics{
		Timestamp: now,
		Hostname:  hostname,
		OS:        os,
		Platform:  platform,
		Uptime:    uptime,
		CPU:       cpu,
		Memory:    memory,
		Disks:     disks,
		Network:   network,
		Load:      load,
	}, nil
}
