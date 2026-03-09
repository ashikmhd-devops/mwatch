package ui

import "github.com/charmbracelet/lipgloss"

// Palette — industrial precision: deep black + electric cyan + amber + crimson
var (
	ColorBg      = lipgloss.Color("#0d0d0d")
	ColorSurface = lipgloss.Color("#141414")
	ColorBorder  = lipgloss.Color("#2a2a2a")
	ColorDim     = lipgloss.Color("#444444")
	ColorMuted   = lipgloss.Color("#666666")
	ColorText    = lipgloss.Color("#c8c8c8")
	ColorBright  = lipgloss.Color("#f0f0f0")

	ColorCyan   = lipgloss.Color("#00d4ff")
	ColorGreen  = lipgloss.Color("#00e676")
	ColorAmber  = lipgloss.Color("#ffb300")
	ColorRed    = lipgloss.Color("#ff3d3d")
	ColorPurple = lipgloss.Color("#bb86fc")

	// Styles
	StyleBase = lipgloss.NewStyle().
			Background(ColorBg).
			Foreground(ColorText)

	StyleHeader = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleValue = lipgloss.NewStyle().
			Foreground(ColorBright).
			Bold(true)

	StyleDim = lipgloss.NewStyle().
			Foreground(ColorDim)

	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	StylePanelTitle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true).
			PaddingLeft(1)

	StyleTableHeader = lipgloss.NewStyle().
				Foreground(ColorDim).
				Bold(true)

	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)

	StyleKeyHint = lipgloss.NewStyle().
			Foreground(ColorDim)

	StyleKeyBind = lipgloss.NewStyle().
			Foreground(ColorAmber).
			Bold(true)

	StyleCritical = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorAmber)

	StyleOK = lipgloss.NewStyle().
		Foreground(ColorGreen)
)

// BarColor returns a color based on percentage value
func BarColor(pct float64) lipgloss.Color {
	switch {
	case pct >= 90:
		return ColorRed
	case pct >= 70:
		return ColorAmber
	default:
		return ColorCyan
	}
}

// Bar renders a filled progress bar using block characters
func Bar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	color := BarColor(pct)
	bar := lipgloss.NewStyle().Foreground(color).Render(repeat("▪", filled))
	rest := lipgloss.NewStyle().Foreground(ColorDim).Render(repeat("╌", empty))
	return bar + rest
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
