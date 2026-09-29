package watcher

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ToolCall representa uma ferramenta invocada pelo modelo.
type ToolCall struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

// Event representa um registro estruturado da transcrição (transcript.jsonl).
type Event struct {
	StepIndex  int        `json:"step_index"`
	Source     string     `json:"source"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	CreatedAt  string     `json:"created_at"`
	Content    string     `json:"content"`
	Thinking   string     `json:"thinking"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	IsQuestion bool       `json:"-"`
	IsDone     bool       `json:"-"`
}

// TranscriptInfo contém metadados de uma sessão descoberta no disco.
type TranscriptInfo struct {
	Path           string
	ConversationID string
	ModTime        time.Time
	Size           int64
}

// FindActiveTranscript localiza a sessão mais recente em ~/.gemini/antigravity/brain/*/transcript.jsonl.
func FindActiveTranscript(geminiDir string) (*TranscriptInfo, error) {
	brainDir := filepath.Join(geminiDir, "brain")
	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler pasta brain: %w", err)
	}

	var latest *TranscriptInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(brainDir, entry.Name(), ".system_generated", "logs", "transcript.jsonl")
		info, err := os.Stat(candidate)
		if err == nil {
			if latest == nil || info.ModTime().After(latest.ModTime) {
				latest = &TranscriptInfo{
					Path:           candidate,
					ConversationID: entry.Name(),
					ModTime:        info.ModTime(),
					Size:           info.Size(),
				}
			}
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("nenhum transcript ativo encontrado em %s", brainDir)
	}

	return latest, nil
}

// ParseLine decodifica uma linha JSONL em um Event estruturado.
func ParseLine(line []byte) (*Event, error) {
	clean := strings.TrimSpace(string(line))
	if len(clean) == 0 {
		return nil, fmt.Errorf("linha vazia")
	}

	var evt Event
	if err := json.Unmarshal([]byte(clean), &evt); err != nil {
		return nil, err
	}

	for _, call := range evt.ToolCalls {
		if call.Name == "ask_question" {
			evt.IsQuestion = true
		}
	}
	if evt.Status == "DONE" {
		evt.IsDone = true
	}

	return &evt, nil
}

// Summary extrai uma descrição compacta e legível do evento para exibição no terminal.
func (e *Event) Summary() string {
	var sb strings.Builder

	switch e.Type {
	case "USER_INPUT":
		content := strings.TrimSpace(e.Content)
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		content = strings.ReplaceAll(content, "\n", " ")
		sb.WriteString(fmt.Sprintf("👤 [User #%d] %s", e.StepIndex, content))

	case "PLANNER_RESPONSE":
		if e.Thinking != "" {
			thinking := strings.TrimSpace(e.Thinking)
			if len(thinking) > 120 {
				thinking = thinking[:117] + "..."
			}
			thinking = strings.ReplaceAll(thinking, "\n", " ")
			sb.WriteString(fmt.Sprintf("💭 [Think #%d] %s\n", e.StepIndex, thinking))
		}
		if len(e.ToolCalls) > 0 {
			for i, tc := range e.ToolCalls {
				if i > 0 {
					sb.WriteString("\n")
				}
				argsSummary := extractArgsSummary(tc.Name, tc.Args)
				if tc.Name == "ask_question" {
					sb.WriteString(fmt.Sprintf("🔔 [INTERAÇÃO #%d] O agente precisa da sua resposta!", e.StepIndex))
				} else {
					sb.WriteString(fmt.Sprintf("🛠️  [Tool #%d] %s(%s)", e.StepIndex, tc.Name, argsSummary))
				}
			}
		} else if e.Thinking == "" {
			sb.WriteString(fmt.Sprintf("🤖 [Agent #%d] Planejando próxima ação...", e.StepIndex))
		}

	case "GENERIC":
		content := strings.TrimSpace(e.Content)
		if len(content) > 80 {
			content = content[:77] + "..."
		}
		content = strings.ReplaceAll(content, "\n", " ")
		if content != "" {
			sb.WriteString(fmt.Sprintf("⚡ [Step #%d] %s", e.StepIndex, content))
		}

	default:
		sb.WriteString(fmt.Sprintf("ℹ️  [%s #%d] Status: %s", e.Type, e.StepIndex, e.Status))
	}

	return sb.String()
}

func extractArgsSummary(toolName string, rawArgs json.RawMessage) string {
	if len(rawArgs) == 0 {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawArgs, &m); err != nil {
		return ""
	}

	switch toolName {
	case "run_command":
		if cmd, ok := m["CommandLine"].(string); ok {
			cmd = strings.TrimSpace(cmd)
			cmd = strings.Trim(cmd, `"`)
			if len(cmd) > 60 {
				return cmd[:57] + "..."
			}
			return cmd
		}
	case "view_file", "replace_file_content":
		if file, ok := m["AbsolutePath"].(string); ok {
			return filepath.Base(strings.Trim(file, `"`))
		}
		if file, ok := m["TargetFile"].(string); ok {
			return filepath.Base(strings.Trim(file, `"`))
		}
	case "write_to_file":
		if file, ok := m["TargetFile"].(string); ok {
			return filepath.Base(strings.Trim(file, `"`))
		}
	case "call_mcp_tool":
		if server, ok := m["ServerName"].(string); ok {
			if tool, ok := m["ToolName"].(string); ok {
				return fmt.Sprintf("%s/%s", server, tool)
			}
		}
	}

	// Resumo padrão se não mapeado
	for _, v := range m {
		if s, ok := v.(string); ok && len(s) > 0 && len(s) < 40 {
			return s
		}
	}
	return ""
}

// Notify emite notificação visual no SO (macOS/Linux) e alerta sonoro no terminal.
func Notify(osName, title, message string) {
	// Alerta sonoro de terminal (bell)
	fmt.Print("\a")

	switch osName {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, message, title)
		_ = exec.Command("osascript", "-e", script).Run()
	case "linux":
		_ = exec.Command("notify-send", title, message).Run()
	}
}

// WatchOptions configura o comportamento do tailing.
type WatchOptions struct {
	Follow       bool
	NotifyOnWait bool
	PollInterval time.Duration
	InitialSteps int
	OSName       string
}

// Stream acompanha o arquivo transcript.jsonl e envia eventos para o handler.
func Stream(ctx context.Context, transcriptPath string, opts WatchOptions, handler func(*Event)) error {
	file, err := os.Open(transcriptPath)
	if err != nil {
		return fmt.Errorf("falha ao abrir transcript: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var allEvents []*Event
	var partial []byte

	// 1. Fase Inicial: lê tudo o que já existe
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			if len(partial) > 0 {
				line = append(partial, line...)
				partial = nil
			}
			if line[len(line)-1] == '\n' {
				if evt, pErr := ParseLine(line); pErr == nil {
					allEvents = append(allEvents, evt)
				}
			} else {
				partial = append(partial, line...)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// Emite os passos iniciais solicitados (ex: últimos N passos)
	startIdx := 0
	if opts.InitialSteps > 0 && len(allEvents) > opts.InitialSteps {
		startIdx = len(allEvents) - opts.InitialSteps
	}
	for i := startIdx; i < len(allEvents); i++ {
		handler(allEvents[i])
	}

	if !opts.Follow {
		return nil
	}

	poll := opts.PollInterval
	if poll <= 0 {
		poll = 400 * time.Millisecond
	}

	ticker := time.NewTicker(poll)
	defer ticker.Stop()

	// 2. Fase de Tailing: escuta novos dados continuamente
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			for {
				line, err := reader.ReadBytes('\n')
				if len(line) > 0 {
					if len(partial) > 0 {
						line = append(partial, line...)
						partial = nil
					}
					if line[len(line)-1] == '\n' {
						if evt, pErr := ParseLine(line); pErr == nil {
							handler(evt)
							if opts.NotifyOnWait && evt.IsQuestion {
								Notify(opts.OSName, "Antigravity Operator", "O agente precisa da sua resposta (pergunta interativa)!")
							}
						}
					} else {
						partial = append(partial, line...)
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}
		}
	}
}
