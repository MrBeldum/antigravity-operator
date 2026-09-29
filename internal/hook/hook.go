package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const preCommitScript = `#!/bin/sh
# Antigravity Operator (agyo) - Session Memory Continuity Hook
# Validates session state before committing code

if [ -f ".agents/session/state.md" ]; then
    echo "📋 [agyo] Auditando memória operacional de sessão..."
    if command -v agyo >/dev/null 2>&1; then
        agyo session status
    fi
fi
exit 0
`

// Install instala o pre-commit hook de sessão no repositório Git local.
func Install(targetDir string) (string, error) {
	gitHooksDir := filepath.Join(targetDir, ".git", "hooks")
	if _, err := os.Stat(gitHooksDir); os.IsNotExist(err) {
		return "", fmt.Errorf("diretório .git/hooks não encontrado em %s (certifique-se de que é a raiz de um repositório git)", targetDir)
	}

	hookPath := filepath.Join(gitHooksDir, "pre-commit")
	if err := os.WriteFile(hookPath, []byte(preCommitScript), 0755); err != nil {
		return "", fmt.Errorf("falha ao gravar hook pre-commit: %w", err)
	}

	return hookPath, nil
}

// Uninstall remove o hook do repositório local.
func Uninstall(targetDir string) error {
	hookPath := filepath.Join(targetDir, ".git", "hooks", "pre-commit")
	content, err := os.ReadFile(hookPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// Remove apenas se foi criado pelo agyo
	if strings.Contains(string(content), "Antigravity Operator") {
		return os.Remove(hookPath)
	}
	return nil
}
