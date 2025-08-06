# CRUSH.md - Development Guide

## Build/Test Commands
```bash
# Build and install locally
./build.sh

# Build for all platforms
./build-all.sh [version]

# Build the dashboard manually
go build -o dashboard main.go

# Run the dashboard
./dashboard

# Run with custom config
./dashboard my-config.json

# Run directly with go
go run main.go

# Run all tests
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

## Code Style Guidelines
- Use `gofmt` for consistent formatting
- Follow Go naming conventions: PascalCase for exported, camelCase for unexported
- Use meaningful variable names, avoid abbreviations
- Keep functions small and focused
- Use early returns to reduce nesting
- Handle errors explicitly, don't ignore them
- Use `context.Context` for cancellation and timeouts
- Prefer composition over inheritance
- Write tests for all public functions
- Use table-driven tests when appropriate