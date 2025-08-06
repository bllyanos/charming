#!/bin/bash

# Charming Dashboard Build Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="charming"
BUILD_DIR="./build"
INSTALL_DIR="$HOME/.local/bin"

echo -e "${BLUE}🫧 Building Charming Dashboard...${NC}"

# Create build directory
mkdir -p "$BUILD_DIR"

# Build for current platform
echo -e "${YELLOW}Building for $(go env GOOS)/$(go env GOARCH)...${NC}"
go build -ldflags="-s -w" -o "$BUILD_DIR/$BINARY_NAME" main.go

# Check if build was successful
if [ ! -f "$BUILD_DIR/$BINARY_NAME" ]; then
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful!${NC}"

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Copy binary to install directory
cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo -e "${GREEN}✅ Installed to $INSTALL_DIR/$BINARY_NAME${NC}"

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}⚠️  Warning: $INSTALL_DIR is not in your PATH${NC}"
    echo -e "${YELLOW}   Add this to your shell profile (.bashrc, .zshrc, etc.):${NC}"
    echo -e "${BLUE}   export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
    echo
fi

# Show usage
echo -e "${BLUE}🚀 Usage:${NC}"
if [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
    echo -e "   ${GREEN}charming${NC}                    # Run with config.json"
    echo -e "   ${GREEN}charming my-config.json${NC}     # Run with custom config"
else
    echo -e "   ${GREEN}$INSTALL_DIR/charming${NC}                    # Run with config.json"
    echo -e "   ${GREEN}$INSTALL_DIR/charming my-config.json${NC}     # Run with custom config"
fi

echo
echo -e "${GREEN}🎉 Charming Dashboard is ready!${NC}"