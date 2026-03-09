package metrics

import (
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo represents information about a process
type ProcessInfo struct {
	PID    int32
	Name   string
	CPU    float64
	MemMB  float64
	MemPct float32
	Status string
	User   string
}

// SortBy specifies the field to sort processes by
type SortBy int

const (
	SortByCPU SortBy = iota
	SortByMem
	SortByPID
	SortByName
)

// ...existing code...

var (
	procCache    []ProcessInfo
	procCacheMu  sync.Mutex
	lastProcTime time.Time
	// procMap persists *process.Process objects across refresh cycles.
	// Re-using the same object lets gopsutil accumulate t0/t1 CPU ticks so
	// Percent(0) returns a real delta instead of 0 on every new fetch.
	procMap = make(map[int32]*process.Process)
)

// procName resolves the best available display name for a process.
// On macOS, p.Name() can fail or return "" for many kernel/opaque processes due
// to SIP and per-user permission restrictions.
// Fallback chain: Name() → basename(Exe()) → "" (caller will skip the entry)
func procName(p *process.Process) string {
	if name, err := p.Name(); err == nil && name != "" {
		return name
	}
	if exe, err := p.Exe(); err == nil && exe != "" {
		return filepath.Base(exe)
	}
	// No readable name or executable — this is a kernel/opaque process.
	// Return "" so the caller's filter drops it rather than showing "[pid]" noise.
	return ""
}

func GetProcesses(sortBy SortBy, limit int) []ProcessInfo {
	procCacheMu.Lock()
	defer procCacheMu.Unlock()

	if time.Since(lastProcTime) > 3*time.Second {
		currentPIDs, err := process.Processes()
		if err == nil {
			// Merge newly seen PIDs into the persistent map.
			// Existing objects already have a t0 CPU measurement, so Percent(0)
			// will return a real delta. Brand-new PIDs get seeded here (returns 0
			// this tick, real value next tick) — same behaviour as htop/top.
			seen := make(map[int32]struct{}, len(currentPIDs))
			for _, p := range currentPIDs {
				seen[p.Pid] = struct{}{}
				if _, exists := procMap[p.Pid]; !exists {
					p.Percent(0) // seed t0 for new PID; result discarded
					procMap[p.Pid] = p
				}
			}
			// Evict processes that are no longer running.
			for pid := range procMap {
				if _, ok := seen[pid]; !ok {
					delete(procMap, pid)
				}
			}

			var list []ProcessInfo
			for _, p := range procMap {
				name := procName(p)
				if name == "" {
					continue
				}
				cpuPct, _ := p.Percent(0) // real delta — object has prior t0
				memInfo, _ := p.MemoryInfo()
				memPct, _ := p.MemoryPercent()
				statuses, _ := p.Status()
				user, _ := p.Username()

				var memMB float64
				if memInfo != nil {
					memMB = float64(memInfo.RSS) / (1024 * 1024)
				}
				stat := "?"
				if len(statuses) > 0 {
					stat = statuses[0]
				}
				list = append(list, ProcessInfo{
					PID:    p.Pid,
					Name:   name,
					CPU:    cpuPct,
					MemMB:  memMB,
					MemPct: memPct,
					Status: stat,
					User:   user,
				})
			}
			procCache = list
			lastProcTime = time.Now()
		}
	}

	result := make([]ProcessInfo, len(procCache))
	copy(result, procCache)

	// Fallback: if all CPU values are 0 (data not warm yet), sort by MEM
	effectiveSort := sortBy
	if effectiveSort == SortByCPU {
		allZero := true
		for _, p := range result {
			if p.CPU > 0 {
				allZero = false
				break
			}
		}
		if allZero {
			effectiveSort = SortByMem
		}
	}

	switch effectiveSort {
	case SortByCPU:
		sort.Slice(result, func(i, j int) bool { return result[i].CPU > result[j].CPU })
	case SortByMem:
		sort.Slice(result, func(i, j int) bool { return result[i].MemMB > result[j].MemMB })
	case SortByPID:
		sort.Slice(result, func(i, j int) bool { return result[i].PID < result[j].PID })
	case SortByName:
		sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}
