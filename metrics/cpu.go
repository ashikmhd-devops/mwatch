package metrics

import (

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

type CPUMetrics struct {
	Overall   float64
	PerCore   []float64
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
	CoreCount int
}

func GetCPU() CPUMetrics {
	overall, err := cpu.Percent(0, false)
	if err != nil || len(overall) == 0 {
		overall = []float64{0}
	}

	perCore, err := cpu.Percent(0, true)
	if err != nil {
		perCore = []float64{}
	}

	avg, err := load.Avg()
	if err != nil {
		avg = &load.AvgStat{}
	}

	count, _ := cpu.Counts(true)

	return CPUMetrics{
		Overall:   overall[0],
		PerCore:   perCore,
		LoadAvg1:  avg.Load1,
		LoadAvg5:  avg.Load5,
		LoadAvg15: avg.Load15,
		CoreCount: count,
	}
}
