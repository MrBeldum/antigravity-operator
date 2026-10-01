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
	CompactedTasks   int
	RetainedTasks    int
	OriginalBytes    int
	CompactedBytes   int
	TokensSavedEst   int
	ArchiveFile      string
	AlreadyCompact   bool
	StateCompacted   bool
	StateBytesSaved  int
	StateArchiveFile string
}

// Compact analisa .agents/session/todo.md e state.md, arquiva tarefas concluídas
// antigas e seções de histórico verboso em .agents/session/archive/ e reescreve
// os arquivos com rollup executivo limpo.
func Compact(targetDir string, opts CompactOptions) (*CompactResult, error) {
	if opts.Threshold <= 0 {
		opts.Threshold = 5
	}
	if opts.KeepLast < 0 {
		opts.KeepLast = 0
	}

	sessionDir := filepath.Join(targetDir, ".agents", "session")
	todoFile := filepath.Join(sessionDir, "todo.md")
	stateFile := filepath.Join(sessionDir, "state.md")

	if !fileExists(todoFile) {
		return nil, fmt.Errorf("arquivo de tarefas não encontrado em %s", todoFile)
	}

	originalContent, err := os.ReadFile(todoFile)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler %s: %w", todoFile, err)
	}

	res := &CompactResult{
		OriginalBytes:  len(originalContent),
		CompactedBytes: len(originalContent),
		AlreadyCompact: true,
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
	res.RetainedTasks = totalCompleted

	archiveDir := filepath.Join(sessionDir, "archive")
	timestamp := time.Now().Format("2006-01-02")

	// 1. Compactação de todo.md se exceder o threshold
	if totalCompleted > opts.Threshold {
		toArchiveCount := totalCompleted - opts.KeepLast
		if toArchiveCount > 0 {
			archiveEntries := completedTasks[:toArchiveCount]
			retainedEntries := completedTasks[toArchiveCount:]

			archiveFileName := fmt.Sprintf("completed-tasks-%s.md", timestamp)
			archivePath := filepath.Join(archiveDir, archiveFileName)
			relativeArchiveRef := filepath.Join("archive", archiveFileName)

			var archiveLines []string
			for _, entry := range archiveEntries {
				archiveLines = append(archiveLines, entry.text)
			}

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

			res.CompactedTasks = toArchiveCount
			res.RetainedTasks = len(retainedEntries)
			res.CompactedBytes = len(newContent)
			res.TokensSavedEst += tokensSaved
			res.ArchiveFile = archivePath
			res.AlreadyCompact = false

			if !opts.DryRun {
				if err := os.MkdirAll(archiveDir, 0755); err != nil {
					return nil, fmt.Errorf("falha ao criar pasta de arquivo %s: %w", archiveDir, err)
				}

				archiveHeader := fmt.Sprintf("\n### 📦 Tarefas Compactadas em %s (%d itens)\n",
					time.Now().Format("2006-01-02 15:04:05"), len(archiveLines))
				archiveBlock := archiveHeader + strings.Join(archiveLines, "\n") + "\n"

				if err := appendToFile(archivePath, archiveBlock); err != nil {
					return nil, fmt.Errorf("falha ao escrever no arquivo %s: %w", archivePath, err)
				}

				if err := os.WriteFile(todoFile, []byte(newContent), 0644); err != nil {
					return nil, fmt.Errorf("falha ao atualizar %s: %w", todoFile, err)
				}
			}
		}
	}

	// 2. Compactação de seções de histórico verboso em state.md
	if fileExists(stateFile) {
		stateContent, err := os.ReadFile(stateFile)
		if err == nil {
			stateLines := strings.Split(string(stateContent), "\n")
			var cleanStateLines []string
			var historyLines []string
			inHistorySection := false

			for _, sLine := range stateLines {
				trimmed := strings.TrimSpace(sLine)
				if strings.HasPrefix(trimmed, "## Histórico") ||
					strings.HasPrefix(trimmed, "## Execuções Anteriores") ||
					strings.HasPrefix(trimmed, "## Resultados Anteriores") ||
					strings.HasPrefix(trimmed, "## Tentativas Anteriores") ||
					strings.HasPrefix(trimmed, "## Logs") {
					inHistorySection = true
				} else if strings.HasPrefix(trimmed, "## ") && inHistorySection {
					inHistorySection = false
				}

				if inHistorySection {
					historyLines = append(historyLines, sLine)
				} else {
					cleanStateLines = append(cleanStateLines, sLine)
				}
			}

			if len(historyLines) > 0 {
				newCleanState := strings.Join(cleanStateLines, "\n")
				diff := len(stateContent) - len(newCleanState)
				if diff > 0 {
					res.StateCompacted = true
					res.StateBytesSaved = diff
					res.TokensSavedEst += diff / 4
					res.AlreadyCompact = false

					stateArchiveFileName := fmt.Sprintf("state-history-%s.md", timestamp)
					stateArchivePath := filepath.Join(archiveDir, stateArchiveFileName)
					res.StateArchiveFile = stateArchivePath

					if !opts.DryRun {
						if err := os.MkdirAll(archiveDir, 0755); err != nil {
							return nil, fmt.Errorf("falha ao criar pasta de arquivo %s: %w", archiveDir, err)
						}

						historyHeader := fmt.Sprintf("\n### 📜 Histórico de Estado Arquivado em %s\n",
							time.Now().Format("2006-01-02 15:04:05"))
						historyBlock := historyHeader + strings.Join(historyLines, "\n") + "\n"

						if err := appendToFile(stateArchivePath, historyBlock); err != nil {
							return nil, fmt.Errorf("falha ao escrever no arquivo de histórico de estado %s: %w", stateArchivePath, err)
						}

						if err := os.WriteFile(stateFile, []byte(newCleanState), 0644); err != nil {
							return nil, fmt.Errorf("falha ao atualizar %s: %w", stateFile, err)
						}
					}
				}
			}
		}
	}

	return res, nil
}

func appendToFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
