package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Soft lavender background with white text
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F6FC")). // off-white
			Background(lipgloss.Color("#9B6FC7")). // soft purple
			Bold(true).
			Padding(0, 2).
			Margin(1, 0).
			Align(lipgloss.Center)

	// Muted lavender foreground
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C2A4D8")). // soft lilac
			Italic(true).
			Padding(1, 1)

	// Professional pink-purple accent
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B96EBE")). // muted magenta
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9C8CB9")). // soft status
			Bold(true)

	skeletonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A89BB2")). // lavender grey
			Italic(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8F7AAB")). // soft muted purple
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EDE7F6")) // very light purple

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D7749C")). // muted rose
			Italic(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B39DDB")) // soft purple separator
)
