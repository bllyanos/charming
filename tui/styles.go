package tui

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("#9932CC")).
		Bold(true).
		Padding(0, 2).
		Margin(1, 0).
		Align(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		Padding(1, 1)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	statusStyle = lipgloss.NewStyle().
		Bold(true)

	skeletonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true)

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	valueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Italic(true)

	separatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
)
