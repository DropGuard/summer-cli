#!/bin/bash
set -e

# Summer CLI Installer

echo "Installing Summer CLI..."

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo "Unsupported OS: $OS"
    exit 1
fi

# Define binary URL (pointing to the latest GitHub Release)
BINARY_URL="https://github.com/DropGuard/summer-cli/releases/latest/download/summer-${OS}-${ARCH}"
INSTALL_DIR="/usr/local/bin"
BINARY_PATH="${INSTALL_DIR}/summer"

echo "Downloading $BINARY_URL..."

# Download binary to a temporary file
TMP_FILE=$(mktemp)
if ! curl -fsSL "$BINARY_URL" -o "$TMP_FILE"; then
    echo "Error: Failed to download the binary. Please ensure the repository has published a release for your OS/Arch."
    rm -f "$TMP_FILE"
    exit 1
fi

# Move and make executable
echo "Installing to $BINARY_PATH (may require sudo)..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$BINARY_PATH"
    chmod +x "$BINARY_PATH"
else
    sudo mv "$TMP_FILE" "$BINARY_PATH"
    sudo chmod +x "$BINARY_PATH"
fi

echo "✅ Summer CLI installed successfully!"
echo "Run 'summer --help' to get started."
