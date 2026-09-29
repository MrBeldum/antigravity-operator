package exporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExportMarkdown(t *testing.T) {
	tempDir := t.TempDir()
	sessionDir := filepath.Join(tempDir, ".agents", "session")
	err := os.MkdirAll(sessionDir, 0755)
	if err != nil {
		t.Fatalf("falha ao criar pasta de sessão: %v", err)
	}

	stateContent := "# Estado da Sessão\n## Objetivo Atual\n- Testar exportador\n## Status em Tempo Real\n- **Fase Atual:** Execução"
	os.WriteFile(filepath.Join(sessionDir, "state.md"), []byte(stateContent), 0644)

	todoContent := "# Lista de Tarefas\n- [x] Fazer setup\n- [ ] Implementar relatorio"
	os.WriteFile(filepath.Join(sessionDir, "todo.md"), []byte(todoContent), 0644)

	decisionsContent := "# Decisões\n### Usar Go puro\nDecidimos usar text/template e html/template para evitar libs externas.\nCom isso o código fica mais leve."
	os.WriteFile(filepath.Join(sessionDir, "decisions.md"), []byte(decisionsContent), 0644)

	opts := ExportOptions{
		Format:    "markdown",
		TargetDir: tempDir,
	}

	result, err := Export(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if !strings.Contains(result, "Testar exportador") {
		t.Errorf("esperava 'Testar exportador' no relatório, mas não foi encontrado")
	}
	if !strings.Contains(result, "[x] Fazer setup") {
		t.Errorf("esperava tarefa concluída, mas não encontrou")
	}
	if !strings.Contains(result, "**Usar Go puro**") {
		t.Errorf("esperava decisão, mas não encontrou")
	}
	if !strings.Contains(result, "50%") {
		t.Errorf("esperava 50%% de progresso, obteve: %s", result)
	}
}

func TestExportHTML(t *testing.T) {
	tempDir := t.TempDir()
	sessionDir := filepath.Join(tempDir, ".agents", "session")
	err := os.MkdirAll(sessionDir, 0755)
	if err != nil {
		t.Fatalf("falha ao criar pasta de sessão: %v", err)
	}

	os.WriteFile(filepath.Join(sessionDir, "state.md"), []byte("## Objetivo Atual\n- Test HTML\n## Status em Tempo Real\n- **Fase Atual:** Final"), 0644)
	os.WriteFile(filepath.Join(sessionDir, "todo.md"), []byte("- [x] Done\n- [ ] Pending"), 0644)

	transcriptPath := filepath.Join(tempDir, "transcript.jsonl")
	transcriptContent := "{\"type\": \"PLANNER_RESPONSE\", \"thinking\": \"Let's test\", \"created_at\": \"" + time.Now().Format(time.RFC3339) + "\"}\n" +
		"{\"type\": \"PLANNER_RESPONSE\", \"tool_calls\": [{\"function\": {\"name\": \"invoke_subagent\"}}], \"created_at\": \"" + time.Now().Add(5*time.Second).Format(time.RFC3339) + "\"}\n"
	os.WriteFile(transcriptPath, []byte(transcriptContent), 0644)

	opts := ExportOptions{
		Format:         "html",
		TargetDir:      tempDir,
		TranscriptPath: transcriptPath,
	}

	result, err := Export(opts)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if !strings.Contains(result, "Test HTML") {
		t.Errorf("esperava 'Test HTML' no html, obteve: %s", result)
	}
	if !strings.Contains(result, "Subagentes") || !strings.Contains(result, "1") {
		t.Errorf("esperava métricas de subagentes, não encontrou")
	}
}
