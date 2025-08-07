package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	RefreshInterval string    `json:"refresh_interval"`
	Services        []Service `json:"services"`
}

type Service struct {
	URL             string     `json:"url"`
	Title           string     `json:"title"`
	Headers         []string   `json:"headers"`
	RefreshInterval string     `json:"refresh_interval"`
	Selectors       []Selector `json:"selectors"`
}

type Selector struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetConfigPath determines the path to the charming_config.json file.
// It first checks the current working directory, then the user's home directory.
func GetConfigPath(filename string) (string, error) {
	// Check current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	currentPath := filepath.Join(cwd, filename)
	if _, err := os.Stat(currentPath); err == nil {
		return currentPath, nil
	}

	// Check home directory
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	homePath := filepath.Join(usr.HomeDir, filename)
	if _, err := os.Stat(homePath); err == nil {
		return homePath, nil
	}

	return "", os.ErrNotExist // File not found in either location
}

func LoadConfig(filepath string) (Config, error) {
	var config Config

	file, err := os.Open(filepath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
