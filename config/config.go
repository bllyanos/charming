package config

import (
	"encoding/json"
	"os"
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

func LoadConfig(filename string) (Config, error) {
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
