#!/usr/bin/env bash
set -e

# ==============================================================================
# Agent-Parkour (parkour / pk) Installer Script
# ==============================================================================

OWNER="yunkon-kim"
REPO="agent-parkour"
BINARY="parkour"
ALIAS="pk"
INSTALL_DIR="/usr/local/bin"

echo "🏃 Installing agent-parkour (parkour / pk)..."

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
  GOBIN="$INSTALL_DIR" go install "github.com/${OWNER}/${REPO}/cmd/parkour@latest"
else
  # Fetch latest release binary from GitHub
  echo "📥 Downloading latest binary from GitHub Releases..."
  DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${BINARY}_${OS}_${ARCH}"
  curl -fsSL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/${BINARY}"
  chmod +x "${INSTALL_DIR}/${BINARY}"
fi

# Create symlinks for 'pk', 'agent-parkour', and legacy 'thop'
ln -sf "${INSTALL_DIR}/${BINARY}" "${INSTALL_DIR}/${ALIAS}"
ln -sf "${INSTALL_DIR}/${BINARY}" "${INSTALL_DIR}/agent-parkour"
ln -sf "${INSTALL_DIR}/${BINARY}" "${INSTALL_DIR}/thop"

echo ""
echo "✅ agent-parkour installed successfully to ${INSTALL_DIR}!"
echo "   • Main binary  : ${INSTALL_DIR}/parkour"
echo "   • Fast alias   : ${INSTALL_DIR}/pk"
echo "   • Legacy alias : ${INSTALL_DIR}/thop"
echo ""
echo "👉 Try running:"
echo "   parkour --help"
echo "   pk --help"
