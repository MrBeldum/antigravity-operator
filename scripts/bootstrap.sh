#!/usr/bin/env bash
set -euo pipefail

# Script de bootstrap do Antigravity Operator (agyo)

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "❌ Arquitetura não suportada: $ARCH" && exit 1 ;;
esac

echo "🚀 Configurando Antigravity Operator para $OS/$ARCH..."

# Compila ou instala se tiver Go local, senão orienta o download
if command -v go >/dev/null 2>&1; then
  echo "🔨 Compilando agyo a partir do fonte local..."
  go build -ldflags="-s -w" -o agyo ./cmd/agyo
  
  TARGET_BIN="/usr/local/bin/agyo"
  if [ -w "/usr/local/bin" ]; then
    cp agyo "$TARGET_BIN"
  else
    mkdir -p "$HOME/.local/bin"
    cp agyo "$HOME/.local/bin/agyo"
    TARGET_BIN="$HOME/.local/bin/agyo"
  fi
  echo "✅ Instalado com sucesso em: $TARGET_BIN"
else
  echo "⚠️  Go não encontrado. Copie o binário pré-compilado correspondente em bin/ para o seu PATH."
fi
