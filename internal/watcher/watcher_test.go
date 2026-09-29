package watcher

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFindActiveTranscript(t *testing.T) {
	tempGemini := t.TempDir()
	brainDir := filepath.Join(tempGemini, "brain")

	// 1. Empty brain dir
	_, err := FindActiveTranscript(tempGemini)
	if err == nil {
		t.Error("expected error for empty brain dir")
	}

	// 2. Multiple conversations with different timestamps
	conv1 := filepath.Join(brainDir, "conv-1", ".system_generated", "logs")
	conv2 := filepath.Join(brainDir, "conv-2", ".system_generated", "logs")
	_ = os.MkdirAll(conv1, 0755)
	_ = os.MkdirAll(conv2, 0755)

	file1 := filepath.Join(conv1, "transcript.jsonl")
	file2 := filepath.Join(conv2, "transcript.jsonl")

	_ = os.WriteFile(file1, []byte(`{"step_index":1}`), 0644)
	time.Sleep(20 * time.Millisecond)
	_ = os.WriteFile(file2, []byte(`{"step_index":2}`), 0644)

	latest, err := FindActiveTranscript(tempGemini)
	if err != nil {
		t.Fatalf("unexpected error finding transcript: %v", err)
	}

	if latest.ConversationID != "conv-2" {
		t.Errorf("expected conv-2 as latest, got %s", latest.ConversationID)
	}
	if latest.Path != file2 {
		t.Errorf("expected path %s, got %s", file2, latest.Path)
	}
}

func TestParseLine(t *testing.T) {
	// 1. Empty line
	if _, err := ParseLine([]byte("   \n")); err == nil {
		t.Error("expected error on empty line")
	}

	// 2. Invalid JSON
	if _, err := ParseLine([]byte("{invalid")); err == nil {
		t.Error("expected error on invalid json")
	}

	// 3. User input
	userJSON := `{"step_index":1,"source":"USER_EXPLICIT","type":"USER_INPUT","status":"DONE","content":"Hello Agent"}`
	evt, err := ParseLine([]byte(userJSON))
	if err != nil {
		t.Fatalf("unexpected error parsing user input: %v", err)
	}
	if evt.StepIndex != 1 || evt.Type != "USER_INPUT" || evt.Content != "Hello Agent" {
		t.Errorf("unexpected event content: %+v", evt)
	}

	// 4. Planner response with thinking and tool calls
	plannerJSON := `{
		"step_index": 2,
		"source": "MODEL",
		"type": "PLANNER_RESPONSE",
		"status": "DONE",
		"thinking": "Need to run test command",
		"tool_calls": [
			{"name": "run_command", "args": {"CommandLine": "go test ./..."}},
			{"name": "ask_question", "args": {"question": "Should I proceed?"}}
		]
	}`
	evtPlanner, err := ParseLine([]byte(plannerJSON))
	if err != nil {
		t.Fatalf("unexpected error parsing planner: %v", err)
	}
	if !evtPlanner.IsQuestion {
		t.Error("expected IsQuestion to be true when ask_question tool is called")
	}
	if !evtPlanner.IsDone {
		t.Error("expected IsDone to be true when status is DONE")
	}
	if len(evtPlanner.ToolCalls) != 2 {
		t.Errorf("expected 2 tool calls, got %d", len(evtPlanner.ToolCalls))
	}
}

func TestEventSummary(t *testing.T) {
	// 1. User input summary
	uEvt := &Event{StepIndex: 1, Type: "USER_INPUT", Content: "Run doctor"}
	if !strings.Contains(uEvt.Summary(), "Run doctor") {
		t.Errorf("expected summary to contain user content: %s", uEvt.Summary())
	}

	// 2. Planner with thinking and tools
	pEvt := &Event{
		StepIndex: 2,
		Type:      "PLANNER_RESPONSE",
		Thinking:  "Thinking deep...",
		ToolCalls: []ToolCall{
			{Name: "run_command", Args: []byte(`{"CommandLine":"git status"}`)},
			{Name: "view_file", Args: []byte(`{"AbsolutePath":"/path/to/main.go"}`)},
			{Name: "call_mcp_tool", Args: []byte(`{"ServerName":"mysql","ToolName":"query"}`)},
			{Name: "ask_question", Args: []byte(`{}`)},
		},
	}
	pSummary := pEvt.Summary()
	if !strings.Contains(pSummary, "Thinking deep...") || !strings.Contains(pSummary, "run_command(git status)") {
		t.Errorf("unexpected planner summary: %s", pSummary)
	}
	if !strings.Contains(pSummary, "main.go") || !strings.Contains(pSummary, "mysql/query") {
		t.Errorf("unexpected args formatting in summary: %s", pSummary)
	}

	// 3. Generic step
	gEvt := &Event{StepIndex: 3, Type: "GENERIC", Content: "Success output"}
	if !strings.Contains(gEvt.Summary(), "Success output") {
		t.Errorf("unexpected generic summary: %s", gEvt.Summary())
	}

	// 4. Default / unknown type
	dEvt := &Event{StepIndex: 4, Type: "CHECKPOINT", Status: "DONE"}
	if !strings.Contains(dEvt.Summary(), "CHECKPOINT") {
		t.Errorf("unexpected checkpoint summary: %s", dEvt.Summary())
	}

	// 5. extractArgsSummary test coverage
	if s := extractArgsSummary("write_to_file", []byte(`{"TargetFile":"/path/to/script.go"}`)); s != "script.go" {
		t.Errorf("expected script.go, got %s", s)
	}
	if s := extractArgsSummary("replace_file_content", []byte(`{"TargetFile":"/path/to/mod.go"}`)); s != "mod.go" {
		t.Errorf("expected mod.go, got %s", s)
	}
	if s := extractArgsSummary("unknown", []byte(`{"name":"foobar"}`)); s != "foobar" {
		t.Errorf("expected foobar, got %s", s)
	}
	if s := extractArgsSummary("empty", nil); s != "" {
		t.Errorf("expected empty string for nil args, got %s", s)
	}
	if s := extractArgsSummary("invalid", []byte(`invalid-json`)); s != "" {
		t.Errorf("expected empty string for invalid json, got %s", s)
	}
}

func TestNotify(t *testing.T) {
	// Should not crash or hang on any platform
	Notify("mock-os", "Test Title", "Test Message")
	Notify("darwin", "Test Title", "Test Message")
	Notify("linux", "Test Title", "Test Message")
}

func TestStream_InvalidFile(t *testing.T) {
	err := Stream(context.Background(), "/nonexistent/path/transcript.jsonl", WatchOptions{}, func(*Event) {})
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestStream(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "transcript.jsonl")

	initialLines := `{"step_index":1,"type":"USER_INPUT","content":"Start"}
{"step_index":2,"type":"PLANNER_RESPONSE","thinking":"Planning"}
`
	_ = os.WriteFile(tempFile, []byte(initialLines), 0644)

	var received []*Event
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := WatchOptions{
		Follow:       true,
		PollInterval: 50 * time.Millisecond,
		InitialSteps: 2,
		OSName:       "mock",
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		// Append a new line dynamically
		f, _ := os.OpenFile(tempFile, os.O_APPEND|os.O_WRONLY, 0644)
		_, _ = f.WriteString("{\"step_index\":3,\"type\":\"GENERIC\",\"content\":\"Appended\"}\n")
		_ = f.Close()

		time.Sleep(150 * time.Millisecond)
		cancel()
	}()

	err := Stream(ctx, tempFile, opts, func(evt *Event) {
		received = append(received, evt)
	})

	if err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}

	if len(received) < 3 {
		t.Errorf("expected at least 3 events, got %d", len(received))
	}
}
