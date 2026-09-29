package main

import (
	"path/filepath"
	"testing"

	"github.com/tiagoboas/antigravity-operator/internal/platform"
)

func TestCLI_PrintUsage(t *testing.T) {
	printUsage()
}

func TestCLI_PrintAbout(t *testing.T) {
	printAbout()
}

func TestCLI_RunInitAndSession(t *testing.T) {
	targetDir := t.TempDir()

	// 1. Initial init
	runInit([]string{targetDir})

	// 2. Force init
	runInit([]string{"--force", targetDir})

	// 3. Session status
	runSession([]string{"status", targetDir})

	// 4. Session archive
	runSession([]string{"archive", targetDir})
}

func TestCLI_RunDoctor(t *testing.T) {
	info, err := platform.Detect()
	if err != nil {
		t.Fatalf("failed to detect platform: %v", err)
	}
	runDoctor(info)
}

func TestCLI_RunSync(t *testing.T) {
	tempGemini := filepath.Join(t.TempDir(), ".gemini", "antigravity")
	info := &platform.Info{
		GeminiDir: tempGemini,
	}
	runSync(info)
}

func TestCLI_RunBrowserStatus(t *testing.T) {
	info := &platform.Info{
		BrowserProfile: t.TempDir(),
	}
	runBrowser(info, []string{"status"})
}
