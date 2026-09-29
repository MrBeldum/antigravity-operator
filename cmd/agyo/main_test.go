package main

import (
	"os"
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

func TestCLI_RunBrowserCommands(t *testing.T) {
	info := &platform.Info{
		BrowserProfile: t.TempDir(),
	}
	// Status
	runBrowser(info, []string{"status"})
}

func TestCLI_RunHook(t *testing.T) {
	tempDir := t.TempDir()
	gitHooksDir := filepath.Join(tempDir, ".git", "hooks")
	_ = os.MkdirAll(gitHooksDir, 0755)

	runHook([]string{"install", tempDir})
	runHook([]string{"uninstall", tempDir})
}



