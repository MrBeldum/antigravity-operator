package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CompactOptions define os parâmetros para a compactação da sessão.
type CompactOptions struct {
	Threshold int  // Quantidade mínima de tarefas concluídas para disparar a compactação (default: 5)
	KeepLast  int  // Quantidade de tarefas concluídas recentes a preservar no todo.md ativo (default: 3)
	DryRun    bool // Se verdadeiro, calcula métricas sem gravar em disco
}

// DefaultCompactOptions retorna as opções padrão recomendadas.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		Threshold: 5,
		KeepLast:  3,
		DryRun:    false,
	}
}

// CompactResult sumariza o impacto da compactação de contexto.
type CompactResult struct {
	CompactedTasks int
	RetainedTasks  int
	OriginalBytes  int
	CompactedBytes int
	TokensSavedEst int
	ArchiveFile    string
	AlreadyCompact bool
}

// Compact analisa .agents/session/todo.md, arquiva tarefas concluídas antigas
// em .agents/session/archive/ e reescreve o todo.md com rollup executivo limpo.
func Compact(targetDir string, opts CompactOptions) (*CompactResult, error) {
	if opts.Threshold <= 0 {
		opts.Threshold = 5
	}
	if opts.KeepLast < 0 {
		opts.KeepLast = 0
	}

	sessionDir := filepath.Join(targetDir, ".agents", "session")
	todoFile := filepath.Join(sessionDir, "todo.md")

	if !fileExists(todoFile) {
		return nil, fmt.Errorf("arquivo de tarefas não encontrado em %s", todoFile)
	}

	originalContent, err := os.ReadFile(todoFile)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler %s: %w", todoFile, err)
	}

	lines := strings.Split(string(originalContent), "\n")
	type taskEntry struct {
		index int
		text  string
	}

	var completedTasks []taskEntry
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [x]") || strings.HasPrefix(trimmed, "- [X]") {
			completedTasks = append(completedTasks, taskEntry{index: i, text: line})
		}
	}

	totalCompleted := len(completedTasks)
	if totalCompleted <= opts.Threshold {
		return &CompactResult{
			CompactedTasks: 0,
			RetainedTasks:  totalCompleted,
			OriginalBytes:  len(originalContent),
			CompactedBytes: len(originalContent),
			TokensSavedEst: 0,
			AlreadyCompact: true,
		}, nil
	}

	toArchiveCount := totalCompleted - opts.KeepLast
	if toArchiveCount <= 0 {
		return &CompactResult{
			CompactedTasks: 0,
			RetainedTasks:  totalCompleted,
			OriginalBytes:  len(originalContent),
			CompactedBytes: len(originalContent),
			TokensSavedEst: 0,
			AlreadyCompact: true,
		}, nil
	}

	archiveEntries := completedTasks[:toArchiveCount]
	retainedEntries := completedTasks[toArchiveCount:]

	archiveDir := filepath.Join(sessionDir, "archive")
	timestamp := time.Now().Format("2006-01-02")
	archiveFileName := fmt.Sprintf("completed-tasks-%s.md", timestamp)
	archivePath := filepath.Join(archiveDir, archiveFileName)
	relativeArchiveRef := filepath.Join("archive", archiveFileName)

	var archiveLines []string
	for _, entry := range archiveEntries {
		archiveLines = append(archiveLines, entry.text)
	}

	// Monta novo conteúdo do todo.md
	// Remove os índices das tarefas arquivadas e insere nota de rollup
	skipIndices := make(map[int]bool)
	for _, entry := range archiveEntries {
		skipIndices[entry.index] = true
	}

	var newLines []string
	rollupInserted := false

	for i, line := range lines {
		if skipIndices[i] {
			if !rollupInserted {
				rollupNote := fmt.Sprintf("> 📦 *Histórico compactado: %d tarefas concluídas arquivadas em [%s](%s)*",
					toArchiveCount, archiveFileName, relativeArchiveRef)
				newLines = append(newLines, rollupNote)
				rollupInserted = true
			}
			continue
		}
		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	tokensSaved := (len(originalContent) - len(newContent)) / 4
	if tokensSaved < 0 {
		tokensSaved = 0
	}

	res := &CompactResult{
		CompactedTasks: toArchiveCount,
		RetainedTasks:  len(retainedEntries),
		OriginalBytes:  len(originalContent),
		CompactedBytes: len(newContent),
		TokensSavedEst: tokensSaved,
		ArchiveFile:    archivePath,
		AlreadyCompact: false,
	}

	if opts.DryRun {
		return res, nil
	}

	// 1. Grava no archive
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar pasta de arquivo %s: %w", archiveDir, err)
	}

	archiveHeader := fmt.Sprintf("\n### 📦 Tarefas Compactadas em %s (%d itens)\n",
		time.Now().Format("2006-01-02 15:04:05"), len(archiveLines))
	archiveBlock := archiveHeader + strings.Join(archiveLines, "\n") + "\n"

	f, err := os.OpenFile(archivePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo de arquivo %s: %w", archivePath, err)
	}
	if _, err := f.WriteString(archiveBlock); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("falha ao escrever no arquivo %s: %w", archivePath, err)
	}
	_ = f.Close()

	// 2. Grava novo todo.md
	if err := os.WriteFile(todoFile, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("falha ao atualizar %s: %w", todoFile, err)
	}

	return res, nil
}
