package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tiagoboas/antigravity-operator/internal/platform"
	"github.com/tiagoboas/antigravity-operator/templates"
)

// SyncResult reporta as ações executadas pelo instalador.
type SyncResult struct {
	RulesPath string
	MCPPath   string
	Updated   []string
}

// Sync instala ou atualiza as configurações do Antigravity na máquina local.
func Sync(info *platform.Info) (*SyncResult, error) {
	result := &SyncResult{}

	// 1. Assegurar diretórios base
	rulesDir := filepath.Join(info.GeminiDir, "rules")
	mcpDir := filepath.Join(info.GeminiDir, "mcp")

	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar pasta de rules: %w", err)
	}
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar pasta de mcp: %w", err)
	}

	// 2. Instalar regra canônica do Session Agent
	ruleContent, err := templates.FS.ReadFile("rules/session-agent.md")
	if err != nil {
		return nil, fmt.Errorf("falha ao ler regra embutida: %w", err)
	}

	sessionRuleTarget := filepath.Join(rulesDir, "session-agent.md")
	if err := os.WriteFile(sessionRuleTarget, ruleContent, 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar regra do session agent: %w", err)
	}
	result.RulesPath = sessionRuleTarget
	result.Updated = append(result.Updated, sessionRuleTarget)

	// 3. MCPs padrão
	mcpDefaultTarget := filepath.Join(mcpDir, "default-servers.json")
	if _, err := os.Stat(mcpDefaultTarget); os.IsNotExist(err) {
		mcpContent, err := templates.FS.ReadFile("mcps/default-servers.json")
		if err == nil {
			_ = os.WriteFile(mcpDefaultTarget, mcpContent, 0644)
			result.MCPPath = mcpDefaultTarget
			result.Updated = append(result.Updated, mcpDefaultTarget)
		}
	} else {
		result.MCPPath = mcpDefaultTarget
	}

	return result, nil
}
