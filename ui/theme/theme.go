package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Purple    = lipgloss.Color("99")
	Gray      = lipgloss.Color("245")
	LightGray = lipgloss.Color("241")
	Red       = lipgloss.Color("196")
	White     = lipgloss.Color("255")

	StatusLineStyle = lipgloss.NewStyle().Background(Purple).Foreground(White).Padding(0, 1)
	HeaderStyle     = lipgloss.NewStyle().Foreground(Purple).Bold(true).Align(lipgloss.Center)
	CellStyle       = lipgloss.NewStyle().Padding(0, 1)
	OddRowStyle     = CellStyle

	EventTypeStyle     = lipgloss.NewStyle().Faint(true)
	SQLKeywordStyle    = lipgloss.NewStyle().Foreground(Purple).Bold(true)
	SQLWordsStyle      = lipgloss.NewStyle().Italic(true)
	ErrorResponseStyle = lipgloss.NewStyle().Foreground(Red).Bold(true)
)
