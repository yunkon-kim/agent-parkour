#!/usr/bin/env bash
set -e

# ==============================================================================
# Token-Hop (thop) Installer Script
# ==============================================================================

OWNER="yunkon-kim"
REPO="token-hop"
BINARY="token-hop"
ALIAS="thop"
INSTALL_DIR="/usr/local/bin"

echo "🦘 Installing token-hop (thop)..."

# Check OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "   Detected OS: $OS, Arch: $ARCH"

# Determine target binary location
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

# If go is installed, try go install first
if command -v go >/dev/null 2>&1; then
  echo "📦 Installing via Go..."
  GOBIN="$INSTALL_DIR" go install "github.com/${OWNER}/${REPO}/cmd/token-hop@latest"
else
  # Fetch latest release binary from GitHub
  echo "📥 Downloading latest binary from GitHub Releases..."
  DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${BINARY}_${OS}_${ARCH}"
  curl -fsSL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/${BINARY}"
  chmod +x "${INSTALL_DIR}/${BINARY}"
fi

# Create symlink for 'thop'
ln -sf "${INSTALL_DIR}/${BINARY}" "${INSTALL_DIR}/${ALIAS}"

echo ""
echo "✅ token-hop and 'thop' alias installed successfully to ${INSTALL_DIR}!"
echo ""
echo "👉 Try running:"
echo "   thop --help"
