package session

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tiagoboas/antigravity-operator/templates"
)

// Result resume o que foi criado durante a inicialização.
type Result struct {
	SessionDir string
	Created    []string
	Skipped    []string
}

// Init inicializa a pasta de memória operacional de sessão (.agents/session/) no diretório alvo.
func Init(targetDir string, force bool) (*Result, error) {
	sessionDir := filepath.Join(targetDir, ".agents", "session")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de sessão: %w", err)
	}

	result := &Result{
		SessionDir: sessionDir,
	}

	// 1. Arquivos canônicos de memória operacional
	files := []struct {
		destName     string
		templatePath string
	}{
		{"state.md", "session/state.md"},
		{"decisions.md", "session/decisions.md"},
		{"todo.md", "session/todo.md"},
	}

	for _, f := range files {
		destPath := filepath.Join(sessionDir, f.destName)
		if fileExists(destPath) && !force {
			result.Skipped = append(result.Skipped, destPath)
			continue
		}

		content, err := templates.FS.ReadFile(f.templatePath)
		if err != nil {
			return nil, fmt.Errorf("falha ao ler template %s: %w", f.templatePath, err)
		}

		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return nil, fmt.Errorf("falha ao escrever %s: %w", destPath, err)
		}
		result.Created = append(result.Created, destPath)
	}

	// 2. Criar .gitignore dentro de .agents para proteger logs e dados sensíveis
	agentsGitignore := filepath.Join(targetDir, ".agents", ".gitignore")
	if !fileExists(agentsGitignore) {
		gitignoreContent := []byte("# Proteção de dados operacionais e temporários do agente\n*.log\ntmp/\ncache/\n*.dump\ncredentials*\n")
		_ = os.WriteFile(agentsGitignore, gitignoreContent, 0644)
		result.Created = append(result.Created, agentsGitignore)
	}

	return result, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
