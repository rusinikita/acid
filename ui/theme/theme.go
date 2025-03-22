package theme

import "github.com/charmbracelet/lipgloss"

var (
	Purple    = lipgloss.Color("99")
	Gray      = lipgloss.Color("245")
	LightGray = lipgloss.Color("241")

	HeaderStyle = lipgloss.NewStyle().Foreground(Purple).Bold(true).Align(lipgloss.Center)
	CellStyle   = lipgloss.NewStyle().Padding(0, 1)
	OddRowStyle = CellStyle.Foreground(Gray)
)
