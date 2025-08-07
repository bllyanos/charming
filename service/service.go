package service

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bllyanos/charming/config"
	"github.com/charmbracelet/bubbletea"
	"github.com/tidwall/gjson"
)

type ServiceData struct {
	Service      config.Service
	Data         map[string]string
	Error        string
	Loading      bool
	HttpStatus   int
	ResponseTime time.Duration
}

func FetchService(index int, service config.Service) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		client := &http.Client{Timeout: 10 * time.Second}

		req, err := http.NewRequest("GET", service.URL, nil)
		if err != nil {
			return FetchMsg{Index: index, Err: err, ResponseTime: 0} // No response time for request creation error
		}

		// Add custom headers
		for key, values := range prepareHeaders(service.Headers) {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := client.Do(req)
		responseTime := time.Since(start)

		if err != nil {
			return FetchMsg{Index: index, Err: err, ResponseTime: responseTime}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return FetchMsg{Index: index, Err: err, HttpStatus: resp.StatusCode, ResponseTime: responseTime}
		}

		data := make(map[string]string)
		for _, selector := range service.Selectors {
			result := gjson.Get(string(body), selector.Value)
			data[selector.Name] = result.String()
		}

		return FetchMsg{Index: index, Data: data, HttpStatus: resp.StatusCode, ResponseTime: responseTime}
	}
}

func prepareHeaders(headerLines []string) http.Header {
	headers := make(http.Header)
	for _, line := range headerLines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			log.Printf("Warning: Invalid header format: %s. Expected 'Key: Value'.", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Environment variable substitution
		if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			envVarName := strings.TrimSuffix(strings.TrimPrefix(value, "{"), "}")
			envValue := os.Getenv(envVarName)
			// If environment variable is not set, value will become empty.
			value = envValue
		}

		// Only add the header if both key and value are not empty after processing.
		if key != "" && value != "" {
			headers.Add(key, value)
		}
	}
	return headers
}

type FetchMsg struct {
	Index        int
	Data         map[string]string
	Err          error
	HttpStatus   int
	ResponseTime time.Duration
}
