package tui

import (
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

	// Fetch initial data for all services
	for i := range m.services {
		cmds = append(cmds, service.FetchService(i, m.services[i].Service))
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}

	return m, nil
}
