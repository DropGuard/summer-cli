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

# Download binary and its checksum; verify before installing. A release without
# SHA256SUMS is treated as untrustworthy rather than silently accepted.
SUMS_URL="https://github.com/DropGuard/summer-cli/releases/latest/download/SHA256SUMS"
TMP_FILE=$(mktemp)
TMP_SUMS=$(mktemp)
cleanup() { rm -f "$TMP_FILE" "$TMP_SUMS"; }
trap cleanup EXIT

if ! curl -fsSL "$BINARY_URL" -o "$TMP_FILE"; then
    echo "Error: Failed to download the binary. Please ensure the repository has published a release for your OS/Arch."
    exit 1
fi

if ! curl -fsSL "$SUMS_URL" -o "$TMP_SUMS"; then
    echo "Error: Failed to download SHA256SUMS — refusing to install an unverified binary."
    exit 1
fi

EXPECTED_NAME="summer-${OS}-${ARCH}"
WANT=$(grep " ${EXPECTED_NAME}\$" "$TMP_SUMS" | awk '{print $1}')
if [ -z "$WANT" ]; then
    echo "Error: no checksum entry for $EXPECTED_NAME in SHA256SUMS."
    exit 1
fi
GOT=$(sha256sum "$TMP_FILE" | awk '{print $1}')
if [ "$GOT" != "$WANT" ]; then
    echo "Error: checksum mismatch (want $WANT, got $GOT). Aborting install."
    exit 1
fi
echo "Checksum verified."

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
