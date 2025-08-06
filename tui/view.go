package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bllyanos/charming/service"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) updateViewportContent() {
	// Generate content for viewport
	var serviceViews []string
	for i, s := range m.services {
		serviceViews = append(serviceViews, renderService(s, i, m.viewport.Width, m.spinner))
	}
	content := strings.Join(serviceViews, "\n")
	m.viewport.SetContent(content)
}

func (m *Model) View() string {
	if m.viewport.Width == 0 {
		return "Loading..."
	}

	// Update viewport content
	m.updateViewportContent()

	// Header with single color background
	header := headerStyle.Width(m.viewport.Width).Render("/////////////// CHARMING ///////////////")

	// Footer with status and scroll info
	lastUpdateStr := "Loading initial data..."
	if !m.lastUpdate.IsZero() {
		lastUpdateStr = "Last update: " + m.lastUpdate.Format("15:04:05")
	}

	// Add scroll info
	scrollInfo := ""
	if m.viewport.TotalLineCount() > m.viewport.Height {
		scrollInfo = fmt.Sprintf(" • %d%%", int(float64(m.viewport.YOffset)/float64(m.viewport.TotalLineCount()-m.viewport.Height)*100))
	}

	footer := footerStyle.Render(
		fmt.Sprintf("%s%s • Press 'q' to quit, 'r' or space to refresh • ↑↓ to scroll", lastUpdateStr, scrollInfo))

	return lipgloss.JoinVertical(lipgloss.Left, header, m.viewport.View(), footer)
}

func renderService(service service.ServiceData, index int, width int, spinner int) string {
	// Spinner characters for loading state
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// Status indicator
	var statusIcon, statusColor string
	if service.Error != "" {
		statusIcon = "✗"
		statusColor = "196" // Red
	} else if service.Loading {
		statusIcon = spinnerChars[spinner%len(spinnerChars)]
		statusColor = "214" // Orange
	} else if len(service.Data) > 0 {
		statusIcon = "✓"
		statusColor = "42" // Green
	} else {
		statusIcon = "○"
		statusColor = "240" // Gray
	}

	// Service title with status
	titleLine := fmt.Sprintf("%s %s %s",
		statusStyle.Foreground(lipgloss.Color(statusColor)).Render(statusIcon),
		titleStyle.Render(service.Service.Title),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(service.Service.URL))

	// Status info line (HTTP status and response time)
	var statusInfo string
	if service.Loading {
		// Placeholder during loading
		statusInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("  " + strings.Repeat("▒", 12) + " • " + strings.Repeat("▒", 8))
	} else if service.HttpStatus > 0 {
		statusColor := "42" // Green
		if service.HttpStatus >= 400 {
			statusColor = "196" // Red
		} else if service.HttpStatus >= 300 {
			statusColor = "214" // Orange
		}

		statusText := lipgloss.NewStyle().
			Foreground(lipgloss.Color(statusColor)).
			Render(fmt.Sprintf("HTTP %d", service.HttpStatus))

		responseTimeText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Render(fmt.Sprintf("%v", service.ResponseTime.Round(time.Millisecond)))

		statusInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("  %s • %s", statusText, responseTimeText))
	}

	// Content area
	var contentLines []string

	if service.Error != "" {
		contentLines = append(contentLines, "  "+errorStyle.Render("⚠ "+service.Error))
	} else if service.Loading {
		// Skeleton loading based on number of selectors
		for _, selector := range service.Service.Selectors {
			// Create skeleton line with varying lengths
			skeletonLength := 20 + (len(selector.Name) % 15) // Vary length based on selector name
			skeleton := strings.Repeat("▒", skeletonLength)

			line := fmt.Sprintf("  %s %s",
				labelStyle.Render(selector.Name+":"),
				skeletonStyle.Render(skeleton))
			contentLines = append(contentLines, line)
		}
	} else {
		for _, selector := range service.Service.Selectors {
			value, exists := service.Data[selector.Name]
			if !exists {
				value = "N/A"
			}

			// Truncate long values
			if len(value) > 50 {
				value = value[:47] + "..."
			}

			line := fmt.Sprintf("  %s %s",
				labelStyle.Foreground(lipgloss.Color("33")).Render(selector.Name+":"),
				valueStyle.Render(value))
			contentLines = append(contentLines, line)
		}
	}

	// Combine everything
	var result []string
	result = append(result, titleLine)
	result = append(result, statusInfo) // Always show status info (placeholder or actual)
	result = append(result, contentLines...)

	// Add separator except for last item
	separator := separatorStyle.Render(strings.Repeat("─", width-4))
	result = append(result, separator)

	return lipgloss.JoinVertical(lipgloss.Left, result...)
}