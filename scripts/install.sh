#!/usr/bin/env bash
set -euo pipefail

# Antigravity Operator (agyo) — Universal One-Line Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/tiagovilasboas/antigravity-operator/main/scripts/install.sh | bash

REPO="tiagovilasboas/antigravity-operator"
BINARY="agyo"
ALIAS="antigravity-operator"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

case "$OS" in
  darwin|linux) ;;
  *) echo "❌ Unsupported OS: $OS" && exit 1 ;;
esac

echo "🚀 Installing Antigravity Operator (agyo) for $OS/$ARCH..."

# Determine installation directory
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

# Try to download the latest release binary from GitHub Releases
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)

DOWNLOADED=0
if [ -n "$LATEST_TAG" ]; then
  TARBALL="agyo_${OS}_${ARCH}.tar.gz"
  URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TARBALL"
  
  TMP_DIR="$(mktemp -d)"
  if curl -sL "$URL" -o "$TMP_DIR/$TARBALL" && tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR" 2>/dev/null; then
    EXTRACTED_BIN="$(find "$TMP_DIR" -type f \( -name "agyo" -o -name "agyo_*" \) ! -name "*.tar.gz" ! -name "*.zip" | head -n 1)"
    if [ -n "$EXTRACTED_BIN" ]; then
      mv "$EXTRACTED_BIN" "$INSTALL_DIR/$BINARY"
      chmod +x "$INSTALL_DIR/$BINARY"
      ln -sf "$INSTALL_DIR/$BINARY" "$INSTALL_DIR/$ALIAS"
      rm -rf "$TMP_DIR"
      DOWNLOADED=1
    fi
  fi
fi

# Fallback: If no release tarball or release is pending, compile from source if Go is installed
if [ "$DOWNLOADED" -eq 0 ]; then
  if command -v go >/dev/null 2>&1; then
    echo "📦 Release binary pending; compiling static binary from source via Go..."
    TMP_SRC="$(mktemp -d)"
    git clone --depth 1 "https://github.com/$REPO.git" "$TMP_SRC"
    (cd "$TMP_SRC" && CGO_ENABLED=0 go build -ldflags="-s -w" -o "$INSTALL_DIR/$BINARY" ./cmd/agyo)
    ln -sf "$INSTALL_DIR/$BINARY" "$INSTALL_DIR/$ALIAS"
    rm -rf "$TMP_SRC"
    DOWNLOADED=1
  else
    echo "❌ Pre-built binary not found for $OS/$ARCH and Go toolchain is missing."
    echo "   Please install Go (https://go.dev) or download from https://github.com/$REPO/releases"
    exit 1
  fi
fi

echo "✅ Successfully installed 'agyo' and 'antigravity-operator' into $INSTALL_DIR!"
echo ""
echo "Verify installation:"
echo "  agyo version"
echo "  agyo doctor"
