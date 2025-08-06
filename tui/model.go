package tui

import (
	"log"
	"time"

	"github.com/bllyanos/charming/config"
	"github.com/bllyanos/charming/service"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	config     config.Config
	services   []service.ServiceData
	viewport   viewport.Model
	lastUpdate time.Time
	spinner    int
}

type spinnerMsg struct{}

type refreshServiceMsg struct {
	Index int
}

func InitialModel(config config.Config) Model {
	services := make([]service.ServiceData, len(config.Services))
	for i, s := range config.Services {
		services[i] = service.ServiceData{
			Service: s,
			Data:    make(map[string]string),
			Loading: true,
		}
	}

	// Initialize viewport
	vp := viewport.New(80, 20) // Default size, will be updated on window resize
	vp.Style = lipgloss.NewStyle().Margin(1, 0)

	return Model{
		config:   config,
		services: services,
		viewport: vp,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Start spinner animation
	cmds = append(cmds, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerMsg{}
	}))

	// Fetch initial data and schedule periodic fetches for all services
	for i := range m.services {
		// Initial fetch
		cmds = append(cmds, service.FetchService(i, m.services[i].Service))

		// Determine refresh interval for the current service
		refreshIntervalStr := m.services[i].Service.RefreshInterval
		if refreshIntervalStr == "" {
			refreshIntervalStr = m.config.RefreshInterval
		}

		// Parse refresh interval, default to 30 seconds if invalid or not set
		refreshDuration, err := time.ParseDuration(refreshIntervalStr)
		if err != nil || refreshDuration <= 0 {
			log.Printf("Warning: Invalid refresh interval '%s' for service %s. Defaulting to 30s.", refreshIntervalStr, m.services[i].Service.Title)
			refreshDuration = 30 * time.Second
		}

		// Schedule periodic refresh
		cmds = append(cmds, tea.Tick(refreshDuration, func(t time.Time) tea.Msg {
			return refreshServiceMsg{Index: i}
		}))
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd // This will hold the command to return

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
				cmds = append(cmds, service.FetchService(i, m.services[i].Service))
			}
			return m, tea.Batch(cmds...)
		default:
			// Pass other keys to viewport for scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case service.FetchMsg:
		if msg.Index < len(m.services) {
			m.services[msg.Index].Loading = false
			m.services[msg.Index].HttpStatus = msg.HttpStatus
			m.services[msg.Index].ResponseTime = msg.ResponseTime
			if msg.Err != nil {
				m.services[msg.Index].Error = msg.Err.Error()
				m.services[msg.Index].Data = make(map[string]string)
			} else {
				m.services[msg.Index].Error = ""
				m.services[msg.Index].Data = msg.Data
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

	case refreshServiceMsg:
		if msg.Index < len(m.services) {
			// Mark service as loading and clear error before fetching
			m.services[msg.Index].Loading = true
			m.services[msg.Index].Error = ""

			// Fetch and reschedule
			refreshCmd := service.FetchService(msg.Index, m.services[msg.Index].Service)

			refreshIntervalStr := m.services[msg.Index].Service.RefreshInterval
			if refreshIntervalStr == "" {
				refreshIntervalStr = m.config.RefreshInterval
			}

			refreshDuration, err := time.ParseDuration(refreshIntervalStr)
			if err != nil || refreshDuration <= 0 {
				log.Printf("Warning: Invalid refresh interval '%s' for service %s. Defaulting to 30s.", refreshIntervalStr, m.services[msg.Index].Service.Title)
				refreshDuration = 30 * time.Second
			}
			rescheduleCmd := tea.Tick(refreshDuration, func(t time.Time) tea.Msg {
				return refreshServiceMsg{Index: msg.Index}
			})

			return m, tea.Batch(refreshCmd, rescheduleCmd)
		}
		return m, nil // Should not happen if index is out of bounds
	}

	return m, cmd // Return the model and the command
}
