package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type uiStyles struct {
	app         lipgloss.Style
	header      lipgloss.Style
	divider     lipgloss.Style
	row         lipgloss.Style
	rowSelected lipgloss.Style
	todoText    lipgloss.Style
	todoDone    lipgloss.Style
	cursor      lipgloss.Style
	footer      lipgloss.Style
	hint        lipgloss.Style
	status      lipgloss.Style
	error       lipgloss.Style
}

func newUIStyles() uiStyles {
	return uiStyles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")),
		divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		row: lipgloss.NewStyle().
			PaddingLeft(0),
		rowSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("230")),
		todoText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		todoDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")),
		cursor: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")),
		footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")),
		hint: lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")),
		status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")),
		error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")),
	}
}

func dividerLine(width int) string {
	if width < 24 {
		width = 24
	}
	return strings.Repeat("─", width)
}
