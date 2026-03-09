package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
)

type NetMetrics struct {
	BytesSentPS float64
	BytesRecvPS float64
	TotalSent   float64
	TotalRecv   float64
}

type DiskMetrics struct {
	Path        string
	TotalGB     float64
	UsedGB      float64
	FreeGB      float64
	UsedPercent float64
	ReadMBps    float64
	WriteMBps   float64
}

var (
	lastNetStats []net.IOCountersStat
	lastNetTime  time.Time
	lastDiskIO   []disk.IOCountersStat
	lastDiskTime time.Time
)

func GetNetwork() NetMetrics {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		return NetMetrics{}
	}

	now := time.Now()
	result := NetMetrics{
		TotalSent: float64(stats[0].BytesSent) / (1024 * 1024),
		TotalRecv: float64(stats[0].BytesRecv) / (1024 * 1024),
	}

	if len(lastNetStats) > 0 {
		dt := now.Sub(lastNetTime).Seconds()
		if dt > 0 {
			result.BytesSentPS = float64(stats[0].BytesSent-lastNetStats[0].BytesSent) / dt / 1024
			result.BytesRecvPS = float64(stats[0].BytesRecv-lastNetStats[0].BytesRecv) / dt / 1024
		}
	}

	lastNetStats = stats
	lastNetTime = now
	return result
}

func GetDisk() DiskMetrics {
	usage, err := disk.Usage("/")
	if err != nil {
		return DiskMetrics{}
	}

	gb := float64(1 << 30)
	result := DiskMetrics{
		Path:        "/",
		TotalGB:     float64(usage.Total) / gb,
		UsedGB:      float64(usage.Used) / gb,
		FreeGB:      float64(usage.Free) / gb,
		UsedPercent: usage.UsedPercent,
	}

	now := time.Now()
	ioMap, err := disk.IOCounters()
	if err == nil {
		var curRead, curWrite uint64
		for _, v := range ioMap {
			curRead += v.ReadBytes
			curWrite += v.WriteBytes
		}

		var lastRead, lastWrite uint64
		for _, v := range lastDiskIO {
			lastRead += v.ReadBytes
			lastWrite += v.WriteBytes
		}

		dt := now.Sub(lastDiskTime).Seconds()
		if dt > 0 && len(lastDiskIO) > 0 {
			result.ReadMBps = float64(curRead-lastRead) / dt / (1024 * 1024)
			result.WriteMBps = float64(curWrite-lastWrite) / dt / (1024 * 1024)
		}

		lastDiskTime = now
		// flatten map to slice for next comparison
		lastDiskIO = nil
		for _, v := range ioMap {
			lastDiskIO = append(lastDiskIO, v)
		}
	}

	return result
}
