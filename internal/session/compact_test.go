package session_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tiagoboas/antigravity-operator/internal/session"
)

func TestCompact(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session-compact-test-*")
	if err != nil {
		t.Fatalf("falha ao criar temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = session.Init(tempDir, false)
	if err != nil {
		t.Fatalf("Init falhou: %v", err)
	}

	todoFile := filepath.Join(tempDir, ".agents", "session", "todo.md")

	// 1. Teste de arquivo sem tarefas suficientes para compactar
	opts := session.DefaultCompactOptions()
	res, err := session.Compact(tempDir, opts)
	if err != nil {
		t.Fatalf("Compact falhou em sessão com poucas tarefas: %v", err)
	}
	if !res.AlreadyCompact {
		t.Errorf("esperava AlreadyCompact=true, obteve false")
	}

	// 2. Preencher todo.md com 10 tarefas concluídas e 3 pendentes
	var lines []string
	lines = append(lines, "# Lista de Tarefas")
	lines = append(lines, "## Pendentes")
	lines = append(lines, "- [ ] Pendente 1")
	lines = append(lines, "- [ ] Pendente 2")
	lines = append(lines, "## Concluídas")
	for i := 1; i <= 10; i++ {
		lines = append(lines, fmt.Sprintf("- [x] Tarefa finalizada %d", i))
	}
	_ = os.WriteFile(todoFile, []byte(strings.Join(lines, "\n")), 0644)

	// 3. Teste com DryRun=true
	opts.Threshold = 5
	opts.KeepLast = 3
	opts.DryRun = true

	resDry, err := session.Compact(tempDir, opts)
	if err != nil {
		t.Fatalf("Compact DryRun falhou: %v", err)
	}
	if resDry.CompactedTasks != 7 {
		t.Errorf("esperava 7 tarefas compactadas no DryRun, obteve %d", resDry.CompactedTasks)
	}
	if resDry.RetainedTasks != 3 {
		t.Errorf("esperava 3 tarefas retidas, obteve %d", resDry.RetainedTasks)
	}

	// 4. Executar compactação real (DryRun=false)
	opts.DryRun = false
	resReal, err := session.Compact(tempDir, opts)
	if err != nil {
		t.Fatalf("Compact real falhou: %v", err)
	}
	if resReal.CompactedTasks != 7 {
		t.Errorf("esperava 7 tarefas compactadas, obteve %d", resReal.CompactedTasks)
	}

	// Verificar se arquivo de histórico foi criado
	if _, err := os.Stat(resReal.ArchiveFile); os.IsNotExist(err) {
		t.Errorf("arquivo de archive esperado não existe: %s", resReal.ArchiveFile)
	}

	// Verificar se todo.md foi enxugado
	newTodoBytes, err := os.ReadFile(todoFile)
	if err != nil {
		t.Fatalf("falha ao ler novo todo.md: %v", err)
	}
	newTodoStr := string(newTodoBytes)

	if !strings.Contains(newTodoStr, "Histórico compactado: 7 tarefas") {
		t.Errorf("esperava mensagem de histórico compactado no todo.md")
	}
	if strings.Contains(newTodoStr, "- [x] Tarefa finalizada 1\n") {
		t.Errorf("tarefa antiga 1 não deveria estar no todo.md ativo")
	}
	if !strings.Contains(newTodoStr, "- [x] Tarefa finalizada 10") {
		t.Errorf("tarefa recente 10 deveria ser mantida no todo.md")
	}
	if !strings.Contains(newTodoStr, "- [ ] Pendente 1") {
		t.Errorf("tarefas pendentes devem permanecer intactas")
	}
}
