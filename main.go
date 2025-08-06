package main

import (
	"log"
	"os"

	"github.com/bllyanos/charming/config"
	"github.com/bllyanos/charming/tui"
	"github.com/charmbracelet/bubbletea"
)

func main() {
	configFile := "charming_config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	m := tui.InitialModel(config)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}