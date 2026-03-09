package metrics

import "github.com/shirou/gopsutil/v3/mem"

type MemMetrics struct {
	TotalGB     float64
	UsedGB      float64
	AvailGB     float64
	UsedPercent float64
	SwapTotalGB float64
	SwapUsedGB  float64
	SwapPercent float64
}

func GetMemory() MemMetrics {
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemMetrics{}
	}

	s, err := mem.SwapMemory()
	if err != nil {
		s = &mem.SwapMemoryStat{}
	}

	gb := float64(1 << 30)
	return MemMetrics{
		TotalGB:     float64(v.Total) / gb,
		UsedGB:      float64(v.Used) / gb,
		AvailGB:     float64(v.Available) / gb,
		UsedPercent: v.UsedPercent,
		SwapTotalGB: float64(s.Total) / gb,
		SwapUsedGB:  float64(s.Used) / gb,
		SwapPercent: s.UsedPercent,
	}
}
