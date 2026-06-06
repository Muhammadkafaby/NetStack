package collector

import (
	"github.com/shirou/gopsutil/v3/disk"
)

// CollectDisk gathers disk usage and I/O metrics for each mount point.
func CollectDisk() ([]DiskMetrics, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	// Get disk I/O counters
	ioCounters, err := disk.IOCounters()
	if err != nil {
		// I/O counters may not be available in some environments
		ioCounters = nil
	}

	var disks []DiskMetrics

	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		d := DiskMetrics{
			MountPoint:  p.Mountpoint,
			FSType:      p.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		}

		// Match I/O stats by device name
		if ioCounters != nil {
			if io, ok := ioCounters[p.Device]; ok {
				d.ReadBytes = io.ReadBytes
				d.WriteBytes = io.WriteBytes
				d.ReadCount = io.ReadCount
				d.WriteCount = io.WriteCount
			}
		}

		disks = append(disks, d)
	}

	return disks, nil
}
