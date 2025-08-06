package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tidwall/gjson"
)

type Config struct {
	Services []Service `json:"services"`
}

type Service struct {
	URL       string     `json:"url"`
	Title     string     `json:"title"`
	Selectors []Selector `json:"selectors"`
}

type Selector struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ServiceData struct {
	Service      Service
	Data         map[string]string
	Error        string
	Loading      bool
	HttpStatus   int
	ResponseTime time.Duration
}

type model struct {
	config     Config
	services   []ServiceData
	viewport   viewport.Model
	lastUpdate time.Time
	spinner    int
}

type spinnerMsg struct{}
type fetchMsg struct {
	index        int
	data         map[string]string
	err          error
	httpStatus   int
	responseTime time.Duration
}

func initialModel(config Config) model {
	services := make([]ServiceData, len(config.Services))
	for i, service := range config.Services {
		services[i] = ServiceData{
			Service: service,
			Data:    make(map[string]string),
			Loading: true,
		}
	}

	// Initialize viewport
	vp := viewport.New(80, 20) // Default size, will be updated on window resize
	vp.Style = lipgloss.NewStyle().Margin(1, 0)

	return model{
		config:   config,
		services: services,
		viewport: vp,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Start spinner animation
	cmds = append(cmds, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerMsg{}
	}))

	// Fetch initial data for all services
	for i := range m.services {
		cmds = append(cmds, fetchService(i, m.services[i].Service))
	}

	return tea.Batch(cmds...)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update viewport size, accounting for header and footer
		headerHeight := 3 // Header with margin
		footerHeight := 3 // Footer with margin
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight - footerHeight
		// Update content after resize
		m.updateViewportContent()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r", " ": // Space bar also triggers refresh
			// Manual refresh
			var cmds []tea.Cmd
			for i := range m.services {
				m.services[i].Loading = true
				m.services[i].Error = ""
				cmds = append(cmds, fetchService(i, m.services[i].Service))
			}
			return m, tea.Batch(cmds...)
		default:
			// Pass other keys to viewport for scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case fetchMsg:
		if msg.index < len(m.services) {
			m.services[msg.index].Loading = false
			m.services[msg.index].HttpStatus = msg.httpStatus
			m.services[msg.index].ResponseTime = msg.responseTime
			if msg.err != nil {
				m.services[msg.index].Error = msg.err.Error()
				m.services[msg.index].Data = make(map[string]string)
			} else {
				m.services[msg.index].Error = ""
				m.services[msg.index].Data = msg.data
				m.lastUpdate = time.Now()
			}
			// Update viewport content after data changes
			m.updateViewportContent()
		}
		return m, nil

	case spinnerMsg:
		m.spinner++
		// Update viewport content to refresh spinners
		m.updateViewportContent()
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return spinnerMsg{}
		})
	}

	return m, nil
}

func (m *model) updateViewportContent() {
	// Generate content for viewport
	var serviceViews []string
	for i, service := range m.services {
		serviceViews = append(serviceViews, renderService(service, i, m.viewport.Width, m.spinner))
	}
	content := strings.Join(serviceViews, "\n")
	m.viewport.SetContent(content)
}

func (m *model) View() string {
	if m.viewport.Width == 0 {
		return "Loading..."
	}

	// Update viewport content
	m.updateViewportContent()

	// Header with single color background
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).     // White text
		Background(lipgloss.Color("#9932CC")). // Purple background
		Bold(true).
		Padding(0, 2).
		Margin(1, 0).
		Width(m.viewport.Width).
		Align(lipgloss.Center)

	header := headerStyle.Render("/////////////// CHARMING ///////////////")

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

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		Padding(1, 1)

	footer := footerStyle.Render(
		fmt.Sprintf("%s%s • Press 'q' to quit, 'r' or space to refresh • ↑↓ to scroll", lastUpdateStr, scrollInfo))

	return lipgloss.JoinVertical(lipgloss.Left, header, m.viewport.View(), footer)
}

func renderService(service ServiceData, index int, width int, spinner int) string {
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
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(statusColor)).
		Bold(true)

	titleLine := fmt.Sprintf("%s %s %s",
		statusStyle.Render(statusIcon),
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
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Italic(true)
		contentLines = append(contentLines, "  "+errorStyle.Render("⚠ "+service.Error))
	} else if service.Loading {
		// Skeleton loading based on number of selectors
		skeletonStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

		for _, selector := range service.Service.Selectors {
			// Create skeleton line with varying lengths
			skeletonLength := 20 + (len(selector.Name) % 15) // Vary length based on selector name
			skeleton := strings.Repeat("▒", skeletonLength)

			labelStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Bold(true)

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

			labelStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Bold(true)

			valueStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

			line := fmt.Sprintf("  %s %s",
				labelStyle.Render(selector.Name+":"),
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
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(strings.Repeat("─", min(width-4, 80)))
	result = append(result, separator)

	return lipgloss.JoinVertical(lipgloss.Left, result...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fetchService(index int, service Service) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(service.URL)
		responseTime := time.Since(start)

		if err != nil {
			return fetchMsg{index: index, err: err, responseTime: responseTime}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fetchMsg{index: index, err: err, httpStatus: resp.StatusCode, responseTime: responseTime}
		}

		data := make(map[string]string)
		for _, selector := range service.Selectors {
			result := gjson.Get(string(body), selector.Value)
			data[selector.Name] = result.String()
		}

		return fetchMsg{index: index, data: data, httpStatus: resp.StatusCode, responseTime: responseTime}
	}
}

func loadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func main() {
	configFile := "charming_config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	m := initialModel(config)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
