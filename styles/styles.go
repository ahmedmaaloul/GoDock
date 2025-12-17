package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// App-wide padding
	AppStyle = lipgloss.NewStyle().Margin(1, 2)

	// Table Header
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)

	// Rows
	RowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Selected Row
	SelectedRowStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("63")).
				Bold(true)

	// Status Colors
	StatusRunning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			Bold(true)

	StatusExited = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red
			Bold(true)

	StatusPaused = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")). // Yellow
			Bold(true)

	// Footer help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)
