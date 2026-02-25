#!/bin/bash
# Cross-platform build script

set -e

VERSION="1.0.0"
BINARY="auto-audio-convert"
DIST_DIR="dist"

echo "🏗️  Building $BINARY v$VERSION for multiple platforms..."

rm -rf $DIST_DIR
mkdir -p $DIST_DIR

# Linux AMD64
echo "📦 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $DIST_DIR/$BINARY-linux-amd64

# Linux ARM64 (Raspberry Pi, etc.)
echo "📦 Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $DIST_DIR/$BINARY-linux-arm64

# macOS AMD64 (Intel Macs)
echo "📦 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $DIST_DIR/$BINARY-darwin-amd64

# macOS ARM64 (Apple Silicon)
echo "📦 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $DIST_DIR/$BINARY-darwin-arm64

# Windows AMD64
echo "📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $DIST_DIR/$BINARY-windows-amd64.exe

echo ""
echo "✅ Build complete! Binaries in $DIST_DIR/:"
ls -lh $DIST_DIR/

echo ""
echo "📤 To create release archives:"
echo "  cd $DIST_DIR && tar czf $BINARY-linux-amd64.tar.gz $BINARY-linux-amd64"
echo "  cd $DIST_DIR && zip $BINARY-windows-amd64.zip $BINARY-windows-amd64.exe"
