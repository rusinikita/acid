package ui

import "github.com/charmbracelet/lipgloss"

var (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	headerStyle = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
	oddRowStyle = cellStyle.Foreground(gray)
)
