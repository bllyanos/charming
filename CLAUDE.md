# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Charming is a CLI dashboard built with Bubble Tea (Go TUI library) that displays data from JSON APIs with automatic loading and manual refresh. It's a single-file Go application that fetches data from multiple services and displays them in a scrollable, styled terminal interface.

## Build Commands

```bash
# Build and install locally (recommended for development)
./build.sh

# Build for all platforms (Linux, macOS, Windows)
./build-all.sh [version]

# Manual build
go build -o charming main.go

# Development run
go run main.go
go run main.go custom-config.json
```

## Test Commands

```bash
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -run TestFunctionName ./path/to/package

# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run

# Vet code for issues
go vet ./...

# Tidy dependencies
go mod tidy
```

## Architecture

### Single-File Design
The entire application is contained in `main.go` with these key components:

- **Config structs**: `Config`, `Service`, `Selector` for JSON configuration parsing
- **Model structs**: `ServiceData` for runtime state, `model` for Bubble Tea state management
- **Message types**: `fetchMsg`, `spinnerMsg` for async operations via Bubble Tea's command pattern
- **View rendering**: Uses Lipgloss for styling with a scrollable viewport from Bubbles

### Bubble Tea Pattern
- **Model**: Holds application state (services data, viewport, spinner state)
- **Update**: Handles messages (keyboard input, API responses, spinner ticks)
- **View**: Renders UI with header, scrollable content, and footer
- **Commands**: Async operations return messages (API fetches, timer ticks)

### Data Flow
1. Load JSON config defining services and selectors
2. Initialize all services in loading state
3. Fetch all services concurrently on startup
4. Update UI as responses arrive
5. Manual refresh refetches all services

### Key Libraries
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles/viewport` - Scrollable content
- `github.com/tidwall/gjson` - JSON path selectors

## Configuration

The application uses `charming_config.json` (or custom file) with services array containing:
- `url`: API endpoint
- `title`: Display name (supports emojis)
- `selectors`: Array of name/value pairs where value is a gjson path

## Development Notes

- All styling is done through Lipgloss styles, not raw ANSI codes
- Spinner animation runs on 100ms ticks
- HTTP requests have 10-second timeout
- Viewport handles scrolling automatically
- Status indicators show HTTP codes and response times
- Skeleton loading shows placeholders during API calls