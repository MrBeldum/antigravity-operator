#!/usr/bin/env bash
set -euo pipefail

echo "🚀 Configurando ambiente de desenvolvimento para o Antigravity Operator..."

# 1. Instalar Git Hook de Pre-Commit
HOOKS_DIR=".git/hooks"
if [ -d "$HOOKS_DIR" ]; then
  cat << 'EOF' > "$HOOKS_DIR/pre-commit"
#!/usr/bin/env bash
set -e
echo "🔍 [Hook] Executando go vet..."
go vet ./...

echo "🧪 [Hook] Executando testes unitários..."
go test ./...

echo "✅ Código validado com sucesso!"
EOF
  chmod +x "$HOOKS_DIR/pre-commit"
  echo "✅ Hook de pre-commit configurado em $HOOKS_DIR/pre-commit"
fi

# 2. Compilar binário local
echo "🔨 Compilando binário local..."
make build

# 3. Executar o doctor
echo "🩺 Executando diagnóstico inicial:"
./bin/agyo doctor

echo ""
echo "🎉 Ambiente de desenvolvimento pronto! Para instalar a CLI globalmente, rode: make install"
