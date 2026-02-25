#!/bin/bash
# Quick installer for auto-audio-convert

set -e

REPO="developertyrone/auto-audio-convert"
INSTALL_DIR="/usr/local/bin"
BINARY="auto-audio-convert"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

BINARY_NAME="$BINARY-$OS-$ARCH"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="$BINARY_NAME.exe"
fi

echo "🔍 Detected: $OS ($ARCH)"
echo "📥 Installing $BINARY_NAME..."

# Download from releases
DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME"

if command -v curl &> /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "/tmp/$BINARY"
elif command -v wget &> /dev/null; then
    wget "$DOWNLOAD_URL" -O "/tmp/$BINARY"
else
    echo "❌ Neither curl nor wget found. Install one first."
    exit 1
fi

# Make executable
chmod +x "/tmp/$BINARY"

# Install (requires sudo for /usr/local/bin)
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
else
    echo "🔐 Installing to $INSTALL_DIR (requires sudo)..."
    sudo mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
fi

echo "✅ Installed to $INSTALL_DIR/$BINARY"
echo ""
echo "🎵 Usage:"
echo "  $BINARY --source=/path/to/music --from=flac --to=mp3"
echo ""
echo "⚠️  Don't forget to install ffmpeg:"
case $OS in
    linux)
        echo "  sudo apt install ffmpeg  # Debian/Ubuntu"
        echo "  sudo dnf install ffmpeg  # Fedora"
        ;;
    darwin)
        echo "  brew install ffmpeg"
        ;;
esac
