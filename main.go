package main

import (
	"log"
	"os"

	"github.com/bllyanos/charming/config"
	"github.com/bllyanos/charming/tui"
	"github.com/charmbracelet/bubbletea"
)

func main() {
	var configFile string

	if len(os.Args) > 1 {
		// If a config file is provided as a command-line argument, use it directly.
		configFile = os.Args[1]
	} else {
		// Otherwise, try to find charming_config.json in CWD or home directory.
		var err error
		configFile, err = config.GetConfigPath("charming_config.json")
		if err != nil {
			log.Fatalf("Error finding config file: %v", err)
		}
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	m := tui.InitialModel(cfg)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}