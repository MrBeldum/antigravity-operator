package platform

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Info agrega informações sobre a plataforma e ambiente de execução.
type Info struct {
	OS            string
	Arch          string
	HomeDir       string
	HasDisplay    bool
	ChromeBin     string
	GeminiDir     string
	BrowserProfile string
	HarnessCore   string
}

// Detect inspeciona o sistema operacional e resolve os caminhos canônicos.
func Detect() (*Info, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	info := &Info{
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		HomeDir:        home,
		GeminiDir:      filepath.Join(home, ".gemini", "antigravity"),
		BrowserProfile: filepath.Join(home, ".gemini", "antigravity-browser-profile"),
		HarnessCore:    filepath.Join(home, "Github", "harness-core"),
	}

	info.HasDisplay = checkDisplay(info.OS)
	info.ChromeBin = findChromeBinary(info.OS)

	return info, nil
}

func checkDisplay(osName string) bool {
	if osName == "darwin" {
		return true // macOS desktop assume display ativo por padrão
	}
	// No Linux / BSD, verifica variáveis de ambiente de servidor gráfico
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}

func findChromeBinary(osName string) string {
	if osName == "darwin" {
		standardMacPath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		if _, err := os.Stat(standardMacPath); err == nil {
			return standardMacPath
		}
	}

	if osName == "windows" {
		winCandidates := []string{
			os.Getenv("ProgramFiles") + `\Google\Chrome\Application\chrome.exe`,
			os.Getenv("ProgramFiles(x86)") + `\Google\Chrome\Application\chrome.exe`,
			os.Getenv("LocalAppData") + `\Google\Chrome\Application\chrome.exe`,
		}
		for _, p := range winCandidates {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}

	// Tenta binários disponíveis no PATH (comum em Linux, macOS com brew/symlink e Windows)
	candidates := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"chrome.exe",
	}

	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			return path
		}
	}

	return ""
}
