package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tiagoboas/antigravity-operator/internal/platform"
)

func TestRun(t *testing.T) {
	info := &platform.Info{
		OS:             "darwin",
		Arch:           "arm64",
		HomeDir:        t.TempDir(),
		HasDisplay:     true,
		ChromeBin:      "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		GeminiDir:      filepath.Join(t.TempDir(), ".gemini", "antigravity"),
		BrowserProfile: filepath.Join(t.TempDir(), ".gemini", "antigravity-browser-profile"),
		HarnessCore:    "",
	}

	report := Run(info)
	if report == nil {
		t.Fatal("expected report to not be nil")
	}

	if report.Platform != info {
		t.Errorf("expected platform %v, got %v", info, report.Platform)
	}

	if len(report.Checks) < 5 {
		t.Fatalf("expected at least 5 checks, got %d", len(report.Checks))
	}

	for _, check := range report.Checks {
		if check.Name == "" {
			t.Error("check name should not be empty")
		}
		if check.Status != "OK" && check.Status != "WARN" && check.Status != "FAIL" && check.Status != "INFO" {
			t.Errorf("unexpected check status: %s for check %s", check.Status, check.Name)
		}
	}
}

func TestCheckChrome(t *testing.T) {
	// 1. Missing Chrome
	noChromeInfo := &platform.Info{ChromeBin: ""}
	res := checkChrome(noChromeInfo)
	if res.Status != "FAIL" {
		t.Errorf("expected FAIL for missing chrome, got %s", res.Status)
	}

	// 2. Found Chrome
	chromeInfo := &platform.Info{ChromeBin: "/usr/bin/google-chrome"}
	resOK := checkChrome(chromeInfo)
	if resOK.Status != "OK" {
		t.Errorf("expected OK for existing chrome path, got %s", resOK.Status)
	}
}

func TestCheckHarnessCore(t *testing.T) {
	// 1. Empty HarnessCore (Standalone)
	emptyInfo := &platform.Info{HarnessCore: ""}
	res := checkHarnessCore(emptyInfo)
	if res.Status != "INFO" {
		t.Errorf("expected INFO for empty harness core, got %s", res.Status)
	}

	// 2. Non-existent path
	fakeInfo := &platform.Info{HarnessCore: filepath.Join(t.TempDir(), "nonexistent")}
	resFake := checkHarnessCore(fakeInfo)
	if resFake.Status != "INFO" {
		t.Errorf("expected INFO for non-repo harness core, got %s", resFake.Status)
	}

	// 3. Valid git repo
	gitDir := t.TempDir()
	dotGit := filepath.Join(gitDir, ".git")
	_ = os.MkdirAll(dotGit, 0755)
	validInfo := &platform.Info{HarnessCore: gitDir}
	_ = checkHarnessCore(validInfo) // won't fail
}

func TestCheckGit(t *testing.T) {
	check := checkGit()
	if check.Name != "Git" && check.Name != "Git Identity" {
		t.Errorf("unexpected check name: %s", check.Name)
	}
}

func TestCheckNode(t *testing.T) {
	check := checkNode()
	if check.Name != "NPX (MCP Runtime)" {
		t.Errorf("unexpected check name: %s", check.Name)
	}
	if check.Status != "OK" && check.Status != "WARN" {
		t.Errorf("unexpected check status: %s", check.Status)
	}
}

func TestCheckDevTools(t *testing.T) {
	tmpDir := t.TempDir()
	info := &platform.Info{BrowserProfile: tmpDir}
	check := checkDevTools(info)
	if check.Name != "Chrome DevTools (Port 9222)" {
		t.Errorf("unexpected check name: %s", check.Name)
	}
	// On a random temp dir, DevTools won't be running on port 9222 unless active
	if check.Status != "OK" && check.Status != "WARN" {
		t.Errorf("unexpected check status: %s", check.Status)
	}
}

func TestCheckAntigravity(t *testing.T) {
	info := &platform.Info{OS: "darwin"}
	check := checkAntigravity(info)
	if check.Name != "Google Antigravity" {
		t.Errorf("unexpected check name: %s", check.Name)
	}
	if check.Status != "OK" && check.Status != "INFO" {
		t.Errorf("unexpected status: %s", check.Status)
	}

	infoLinux := &platform.Info{OS: "linux"}
	checkLinux := checkAntigravity(infoLinux)
	if checkLinux.Name != "Google Antigravity" {
		t.Errorf("unexpected check name: %s", checkLinux.Name)
	}
}

func TestCheckAPIKeys(t *testing.T) {
	// Test without keys
	_ = os.Unsetenv("GEMINI_API_KEY")
	_ = os.Unsetenv("GOOGLE_API_KEY")
	resEmpty := checkAPIKeys()
	if resEmpty.Status != "INFO" {
		t.Errorf("expected INFO for empty keys, got %s", resEmpty.Status)
	}

	// Test with GEMINI_API_KEY
	_ = os.Setenv("GEMINI_API_KEY", "AIzaSyTestKey123456789")
	defer os.Unsetenv("GEMINI_API_KEY")
	resSet := checkAPIKeys()
	if resSet.Status != "OK" {
		t.Errorf("expected OK for set key, got %s", resSet.Status)
	}
}

