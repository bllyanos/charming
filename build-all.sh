#!/bin/bash

# Charming Dashboard Cross-Platform Build Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="charming"
BUILD_DIR="./dist"
VERSION=${1:-"dev"}

echo -e "${BLUE}🫧 Building Charming Dashboard v$VERSION for multiple platforms...${NC}"

# Clean and create build directory
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Build targets
declare -a targets=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Build for each target
for target in "${targets[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$target"
    
    output_name="$BINARY_NAME-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        output_name="$output_name.exe"
    fi
    
    echo -e "${YELLOW}Building for $GOOS/$GOARCH...${NC}"
    
    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X main.version=$VERSION" \
        -o "$BUILD_DIR/$output_name" \
        main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ $output_name${NC}"
    else
        echo -e "${RED}❌ Failed to build $output_name${NC}"
        exit 1
    fi
done

# Create archives
echo -e "${BLUE}📦 Creating archives...${NC}"
cd "$BUILD_DIR"

for target in "${targets[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$target"
    
    binary_name="$BINARY_NAME-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        binary_name="$binary_name.exe"
    fi
    
    archive_name="$BINARY_NAME-$VERSION-$GOOS-$GOARCH"
    
    if [ "$GOOS" = "windows" ]; then
        zip -q "$archive_name.zip" "$binary_name" ../config.json ../README.md
        echo -e "${GREEN}✅ $archive_name.zip${NC}"
    else
        tar -czf "$archive_name.tar.gz" "$binary_name" ../config.json ../README.md
        echo -e "${GREEN}✅ $archive_name.tar.gz${NC}"
    fi
done

cd ..

echo
echo -e "${GREEN}🎉 Cross-platform build complete!${NC}"
echo -e "${BLUE}📁 Binaries and archives are in: $BUILD_DIR/${NC}"
echo
echo -e "${BLUE}📋 Built targets:${NC}"
for target in "${targets[@]}"; do
    echo -e "   • $target"
done