package ui

import (
	"fmt"
	"strings"

	"github.com/ashikmhd/mwatch/metrics"
	"github.com/charmbracelet/lipgloss"
)

const barWidth = 20

func RenderCPU(m metrics.CPUMetrics, width int) string {
	title := StylePanelTitle.Render("▸ CPU")

	overall := fmt.Sprintf("%.1f%%", m.Overall)
	overallColored := lipgloss.NewStyle().Foreground(BarColor(m.Overall)).Bold(true).Render(overall)
	overallBar := Bar(m.Overall, barWidth)

	lines := []string{
		title,
		fmt.Sprintf("  %s %s %s",
			StyleLabel.Render("Overall"),
			overallBar,
			overallColored,
		),
		fmt.Sprintf("  %s  %s  %s  %s",
			StyleLabel.Render("Load"),
			StyleValue.Render(fmt.Sprintf("%.2f", m.LoadAvg1)),
			StyleMuted.Render(fmt.Sprintf("%.2f", m.LoadAvg5)),
			StyleMuted.Render(fmt.Sprintf("%.2f", m.LoadAvg15)),
		),
		"",
	}

	// Per-core in 2-column grid
	cols := 2
	for i := 0; i < len(m.PerCore); i += cols {
		row := "  "
		for c := 0; c < cols && i+c < len(m.PerCore); c++ {
			idx := i + c
			pct := m.PerCore[idx]
			label := StyleDim.Render(fmt.Sprintf("P%-2d", idx))
			bar := Bar(pct, 10)
			val := lipgloss.NewStyle().Foreground(BarColor(pct)).Render(fmt.Sprintf("%5.1f%%", pct))
			cell := fmt.Sprintf("%s %s %s", label, bar, val)
			if c == 0 {
				row += fmt.Sprintf("%-36s", cell)
			} else {
				row += cell
			}
		}
		lines = append(lines, row)
	}

	return panel(strings.Join(lines, "\n"), width)
}

func RenderMemory(m metrics.MemMetrics, width int) string {
	title := StylePanelTitle.Render("▸ Memory")

	memBar := Bar(m.UsedPercent, barWidth)
	memVal := fmt.Sprintf("%.1f / %.1f GB", m.UsedGB, m.TotalGB)

	swapBar := Bar(m.SwapPercent, barWidth)
	swapVal := fmt.Sprintf("%.1f / %.1f GB", m.SwapUsedGB, m.SwapTotalGB)

	content := strings.Join([]string{
		title,
		fmt.Sprintf("  %s  %s %s",
			StyleLabel.Render("RAM  "),
			memBar,
			lipgloss.NewStyle().Foreground(BarColor(m.UsedPercent)).Bold(true).Render(memVal),
		),
		fmt.Sprintf("  %s  %s %s",
			StyleLabel.Render("Swap "),
			swapBar,
			lipgloss.NewStyle().Foreground(BarColor(m.SwapPercent)).Render(swapVal),
		),
		fmt.Sprintf("  %s %s",
			StyleLabel.Render("Avail"),
			StyleValue.Render(fmt.Sprintf("%.1f GB", m.AvailGB)),
		),
	}, "\n")

	return panel(content, width)
}

func RenderNetwork(m metrics.NetMetrics, width int) string {
	title := StylePanelTitle.Render("▸ Network")

	upColor := lipgloss.NewStyle().Foreground(ColorAmber).Bold(true)
	downColor := lipgloss.NewStyle().Foreground(ColorCyan).Bold(true)

	content := strings.Join([]string{
		title,
		fmt.Sprintf("  %s  %s",
			StyleLabel.Render("↑ Upload  "),
			upColor.Render(fmt.Sprintf("%.1f KB/s", m.BytesSentPS)),
		),
		fmt.Sprintf("  %s  %s",
			StyleLabel.Render("↓ Download"),
			downColor.Render(fmt.Sprintf("%.1f KB/s", m.BytesRecvPS)),
		),
		fmt.Sprintf("  %s  %s  %s  %s",
			StyleLabel.Render("Total"),
			StyleMuted.Render(fmt.Sprintf("↑%.1fMB", m.TotalSent)),
			StyleDim.Render("/"),
			StyleMuted.Render(fmt.Sprintf("↓%.1fMB", m.TotalRecv)),
		),
	}, "\n")

	return panel(content, width)
}

func RenderDisk(m metrics.DiskMetrics, width int) string {
	title := StylePanelTitle.Render("▸ Disk  /")

	diskBar := Bar(m.UsedPercent, barWidth)
	diskVal := fmt.Sprintf("%.0f / %.0f GB", m.UsedGB, m.TotalGB)

	content := strings.Join([]string{
		title,
		fmt.Sprintf("  %s  %s %s",
			StyleLabel.Render("Used "),
			diskBar,
			lipgloss.NewStyle().Foreground(BarColor(m.UsedPercent)).Bold(true).Render(diskVal),
		),
		fmt.Sprintf("  %s  %s   %s  %s",
			StyleLabel.Render("I/O  "),
			lipgloss.NewStyle().Foreground(ColorCyan).Render(fmt.Sprintf("R: %.2f MB/s", m.ReadMBps)),
			StyleDim.Render("|"),
			lipgloss.NewStyle().Foreground(ColorAmber).Render(fmt.Sprintf("W: %.2f MB/s", m.WriteMBps)),
		),
	}, "\n")

	return panel(content, width)
}

func RenderProcesses(procs []metrics.ProcessInfo, sortBy metrics.SortBy, width int) string {
	headers := []string{"PID", "NAME", "CPU%", "MEM MB", "STATUS", "USER"}
	widths := []int{7, 24, 8, 9, 8, 14}

	titleMap := map[metrics.SortBy]string{
		metrics.SortByCPU:  "CPU%",
		metrics.SortByMem:  "MEM MB",
		metrics.SortByPID:  "PID",
		metrics.SortByName: "NAME",
	}
	activeSort := titleMap[sortBy]

	// Header row
	header := "  "
	for i, h := range headers {
		s := StyleTableHeader
		if h == activeSort {
			s = StyleSelected
		}
		header += s.Render(pad(h, widths[i]))
	}

	lines := []string{
		StylePanelTitle.Render("▸ Processes"),
		header,
		"  " + StyleDim.Render(strings.Repeat("─", width-6)),
	}

	for _, p := range procs {
		cpuStyle := lipgloss.NewStyle().Foreground(BarColor(p.CPU))
		memStyle := lipgloss.NewStyle().Foreground(ColorText)

		row := "  " +
			StyleDim.Render(pad(fmt.Sprintf("%d", p.PID), widths[0])) +
			StyleValue.Render(pad(truncate(p.Name, widths[1]-1), widths[1])) +
			cpuStyle.Render(pad(fmt.Sprintf("%.1f", p.CPU), widths[2])) +
			memStyle.Render(pad(fmt.Sprintf("%.0f", p.MemMB), widths[3])) +
			StyleMuted.Render(pad(p.Status, widths[4])) +
			StyleDim.Render(pad(truncate(p.User, widths[5]-1), widths[5]))
		lines = append(lines, row)
	}

	return panel(strings.Join(lines, "\n"), width)
}

func RenderFooter(width int) string {
	keys := []struct{ key, desc string }{
		{"c", "sort CPU"},
		{"m", "sort MEM"},
		{"p", "sort PID"},
		{"n", "sort NAME"},
		{"q", "quit"},
	}

	parts := []string{}
	for _, k := range keys {
		parts = append(parts,
			StyleKeyBind.Render(k.key)+StyleKeyHint.Render(":"+k.desc),
		)
	}

	footer := strings.Join(parts, StyleDim.Render("  │  "))
	return lipgloss.NewStyle().
		Foreground(ColorDim).
		PaddingLeft(2).
		PaddingTop(1).
		Render(footer)
}

// --- helpers ---

func panel(content string, width int) string {
	return StyleBorder.Width(width - 2).Render(content)
}

func pad(s string, n int) string {
	runes := []rune(s)
	if len(runes) >= n {
		return string(runes[:n])
	}
	return s + strings.Repeat(" ", n-len(runes))
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n-1]) + "~" // single-byte safe ellipsis
	}
	return s
}

var StyleMuted = lipgloss.NewStyle().Foreground(ColorMuted)
