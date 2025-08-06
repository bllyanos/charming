# Charming Dashboard

A simple CLI dashboard built with Bubble Tea that displays data from JSON APIs with automatic initial load and manual refresh.

## Requirements

- **Go 1.24.5+** (required for building)
- **Terminal with color support** (recommended)
- **Internet connection** (for fetching API data)

## Quick Start

### 1. Build & Install (Recommended)
```bash
# Clone the repository
git clone https://github.com/bllyanos/charming.git
cd charming

# Build and install to ~/.local/bin
./build.sh

# Run from anywhere
charming
```

### 2. Manual Build
```bash
# Build the binary
go build -o charming main.go

# Run the dashboard
./charming

# Run with custom config
./charming my-config.json
```

### 3. Development Mode
```bash
# Run directly with Go (for development only)
go run main.go
```

## Features

- Automatic data fetching on startup
- Manual refresh for subsequent updates
- **Scrollable viewport** for handling many services
- List-based layout with individual service spinners
- HTTP status codes and response time tracking
- Stable UI layout with skeleton placeholders (no jumping)
- Skeleton loading animation based on selectors
- JSON path selectors using gjson syntax
- Beautiful gradient header and styling

## Usage

After building, you can run the dashboard with:

```bash
# Use default charming_config.json
charming

# Use custom configuration file
charming my-config.json

# Examples with different configs
charming production.json
charming staging.json
```

## Controls

- `q` or `Ctrl+C`: Quit
- `r` or `Space`: Refresh all services
- `↑` / `↓` or `j` / `k`: Scroll up/down
- `Page Up` / `Page Down`: Scroll by page
- `Home` / `End`: Go to top/bottom

## Behavior

The dashboard automatically fetches data from all configured services when it starts up. After the initial load, you can manually refresh the data using `r` or `Space` key.

## Configuration

Create a `charming_config.json` file with the following structure:

```json
{
  "services": [
    {
      "url": "https://api.example.com/status",
      "title": "🚀 Service Status",
      "selectors": [
        {
          "name": "Version",
          "value": "meta.version"
        },
        {
          "name": "Status",
          "value": "status"
        }
      ]
    }
  ]
}
```

### Configuration Fields

- `services`: Array of service configurations
  - `url`: API endpoint URL
  - `title`: Display title for the service (emojis encouraged!)
  - `selectors`: Array of data selectors
    - `name`: Display name for the field (emojis encouraged!)
    - `value`: JSON path selector (uses gjson syntax)

### JSON Path Examples

- `"status"` - Simple field
- `"meta.version"` - Nested field
- `"items.0.name"` - Array element
- `"users.#.name"` - Array of names
- `"data.@reverse"` - Array modifier

### Skeleton Loading

The dashboard shows skeleton loading animations that match the number of selectors for each service, giving users a preview of what data will be displayed.

## Build & Install

### Quick Install (Recommended)
```bash
# Build and install to ~/.local/bin
./build.sh
```

This will:
- Build an optimized binary
- Install it to `~/.local/bin/charming`
- Make it available system-wide (if `~/.local/bin` is in your PATH)

### Cross-Platform Build
```bash
# Build for all platforms
./build-all.sh v1.0.0
```

Creates binaries and archives for:
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)

### Manual Build
```bash
go build -o charming main.go
./charming
```

## Installation Notes

- The build script installs to `~/.local/bin/charming`
- Make sure `~/.local/bin` is in your PATH for system-wide access
- If PATH is not configured, the build script will show instructions
- Use `which charming` to verify installation

## Troubleshooting

### Command not found
If you get "command not found" after building:
```bash
# Check if ~/.local/bin is in PATH
echo $PATH | grep -q "$HOME/.local/bin" && echo "✓ PATH configured" || echo "✗ PATH not configured"

# Add to your shell profile (.bashrc, .zshrc, etc.)
export PATH="$PATH:$HOME/.local/bin"
```

### Build Issues
- Ensure Go 1.24.5+ is installed: `go version`
- Run `go mod tidy` to resolve dependencies
- Check internet connection for module downloads