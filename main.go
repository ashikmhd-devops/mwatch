package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ashikmhd/mwatch/metrics"
	"github.com/ashikmhd/mwatch/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const tickInterval = 1500 * time.Millisecond

// --- App Model ---

type Model struct {
	cpu    metrics.CPUMetrics
	mem    metrics.MemMetrics
	net    metrics.NetMetrics
	disk   metrics.DiskMetrics
	procs  []metrics.ProcessInfo
	sortBy metrics.SortBy
	width  int
	height int
	tick   int
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func collectMetrics() tea.Cmd {
	return func() tea.Msg {
		return snapshotMsg{
			cpu:   metrics.GetCPU(),
			mem:   metrics.GetMemory(),
			net:   metrics.GetNetwork(),
			disk:  metrics.GetDisk(),
			procs: metrics.GetProcesses(metrics.SortByCPU, 14), // cached, fast
		}
	}
}

func collectAll(sortBy metrics.SortBy) tea.Cmd {
	return func() tea.Msg {
		return snapshotMsg{
			cpu:   metrics.GetCPU(),
			mem:   metrics.GetMemory(),
			net:   metrics.GetNetwork(),
			disk:  metrics.GetDisk(),
			procs: metrics.GetProcesses(sortBy, 14),
		}
	}
}

type snapshotMsg struct {
	cpu   metrics.CPUMetrics
	mem   metrics.MemMetrics
	net   metrics.NetMetrics
	disk  metrics.DiskMetrics
	procs []metrics.ProcessInfo
}

func initialModel() Model {
	return Model{
		sortBy: metrics.SortByCPU,
		width:  120,
		height: 40,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(collectAll(m.sortBy), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.sortBy = metrics.SortByCPU
			return m, collectAll(m.sortBy)
		case "m":
			m.sortBy = metrics.SortByMem
			return m, collectAll(m.sortBy)
		case "p":
			m.sortBy = metrics.SortByPID
			return m, collectAll(m.sortBy)
		case "n":
			m.sortBy = metrics.SortByName
			return m, collectAll(m.sortBy)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.tick++
		return m, tea.Batch(
			collectMetrics(), // fast: cpu/mem/net/disk only
			tickCmd(),
		)

	case snapshotMsg:
		m.cpu = msg.cpu
		m.mem = msg.mem
		m.net = msg.net
		m.disk = msg.disk
		m.procs = msg.procs
	}

	return m, nil
}

func (m Model) View() string {
	w := m.width
	if w < 80 {
		w = 80
	}

	half := w/2 - 1

	// Header
	hostname, _ := os.Hostname()
	uptime := getUptime()

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(ui.ColorCyan).
			Bold(true).
			Render("  ◆ mwatch"),
		lipgloss.NewStyle().
			Foreground(ui.ColorDim).
			Render(fmt.Sprintf("  %s  %s  %s  ●%s",
				hostname,
				runtime.GOOS+"/"+runtime.GOARCH,
				uptime,
				dots(m.tick),
			)),
	)

	headerLine := lipgloss.NewStyle().
		Width(w).
		Background(ui.ColorSurface).
		PaddingTop(0).
		PaddingBottom(0).
		Render(header)

	divider := lipgloss.NewStyle().
		Foreground(ui.ColorBorder).
		Render(strings.Repeat("─", w))

	// Top row: CPU + Memory
	cpuPanel := ui.RenderCPU(m.cpu, half)
	memPanel := ui.RenderMemory(m.mem, half)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, cpuPanel, " ", memPanel)

	// Mid row: Network + Disk
	netPanel := ui.RenderNetwork(m.net, half)
	diskPanel := ui.RenderDisk(m.disk, half)
	midRow := lipgloss.JoinHorizontal(lipgloss.Top, netPanel, " ", diskPanel)

	// Process table (full width)
	procPanel := ui.RenderProcesses(m.procs, m.sortBy, w)

	footer := ui.RenderFooter(w)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerLine,
		divider,
		topRow,
		midRow,
		procPanel,
		footer,
	)
}

// --- Helpers ---

func dots(tick int) string {
	states := []string{"·", "··", "···", "··", "·"}
	return states[tick%len(states)]
}

func getUptime() string {
	out, err := exec.Command("uptime").Output()
	if err != nil {
		return "uptime N/A"
	}
	s := strings.TrimSpace(string(out))
	// macOS uptime: "14:32  up 2 days, 3:14, 2 users, load averages: ..."
	if idx := strings.Index(s, "up "); idx >= 0 {
		s = s[idx:]
		if end := strings.Index(s, ","); end > 0 {
			// get first 2 segments
			parts := strings.SplitN(s, ",", 4)
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[0]) + "," + strings.TrimSpace(parts[1])
			}
			return strings.TrimSpace(parts[0])
		}
	}
	return "running"
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
