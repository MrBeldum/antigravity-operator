package doctor

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/tiagoboas/antigravity-operator/internal/platform"
	"github.com/tiagoboas/antigravity-operator/internal/profile"
)

// CheckItem representa o resultado de uma verificação individual.
type CheckItem struct {
	Name    string
	Status  string // "OK", "WARN", "FAIL"
	Details string
}

// Report agrega todas as verificações do sistema.
type Report struct {
	Platform *platform.Info
	Checks   []CheckItem
}

// Run executa todas as verificações de diagnóstico.
func Run(info *platform.Info) *Report {
	rep := &Report{Platform: info}

	// 1. Git e Configuração de Usuário
	rep.add(checkGit())

	// 2. Google Chrome / Chromium
	rep.add(checkChrome(info))

	// 3. Porta DevTools & Perfil Isolado
	rep.add(checkDevTools(info))

	// 4. Node e NPX (para MCPs)
	rep.add(checkNode())

	// 5. Integração com Harness Core
	rep.add(checkHarnessCore(info))

	return rep
}

func (r *Report) add(items ...CheckItem) {
	r.Checks = append(r.Checks, items...)
}

func checkGit() CheckItem {
	out, err := exec.Command("git", "--version").Output()
	if err != nil {
		return CheckItem{Name: "Git", Status: "FAIL", Details: "git não encontrado no PATH"}
	}
	version := strings.TrimSpace(string(out))

	nameOut, _ := exec.Command("git", "config", "user.name").Output()
	emailOut, _ := exec.Command("git", "config", "user.email").Output()

	name := strings.TrimSpace(string(nameOut))
	email := strings.TrimSpace(string(emailOut))

	if name == "" || email == "" {
		return CheckItem{
			Name:    "Git Identity",
			Status:  "WARN",
			Details: fmt.Sprintf("%s (user.name ou user.email não configurados)", version),
		}
	}

	return CheckItem{
		Name:    "Git",
		Status:  "OK",
		Details: fmt.Sprintf("%s (%s <%s>)", version, name, email),
	}
}

func checkChrome(info *platform.Info) CheckItem {
	if info.ChromeBin == "" {
		return CheckItem{
			Name:    "Google Chrome",
			Status:  "FAIL",
			Details: "Nenhum binário de Chrome/Chromium encontrado no sistema",
		}
	}
	return CheckItem{
		Name:    "Google Chrome",
		Status:  "OK",
		Details: fmt.Sprintf("Localizado em: %s", info.ChromeBin),
	}
}

func checkDevTools(info *platform.Info) CheckItem {
	st := profile.CheckStatus(info.BrowserProfile)
	if st.IsRunning {
		return CheckItem{
			Name:    "Chrome DevTools (Port 9222)",
			Status:  "OK",
			Details: fmt.Sprintf("Ativo (%s) no perfil isolado", st.Version),
		}
	}
	return CheckItem{
		Name:    "Chrome DevTools (Port 9222)",
		Status:  "WARN",
		Details: "Inativo (execute 'agyo browser start' para iniciar)",
	}
}

func checkNode() CheckItem {
	out, err := exec.Command("npx", "--version").Output()
	if err != nil {
		return CheckItem{
			Name:    "NPX (MCP Runtime)",
			Status:  "WARN",
			Details: "npx não encontrado. MCP servers baseados em Node podem falhar",
		}
	}
	return CheckItem{
		Name:    "NPX (MCP Runtime)",
		Status:  "OK",
		Details: fmt.Sprintf("Versão %s disponível", strings.TrimSpace(string(out))),
	}
}

func checkHarnessCore(info *platform.Info) CheckItem {
	if info.HarnessCore == "" {
		return CheckItem{
			Name:    "Harness Core",
			Status:  "INFO",
			Details: "Modo Standalone (regras e templates embutidos)",
		}
	}
	if out, err := exec.Command("git", "-C", info.HarnessCore, "rev-parse", "--is-inside-work-tree").Output(); err == nil && strings.TrimSpace(string(out)) == "true" {
		return CheckItem{
			Name:    "Harness Core",
			Status:  "OK",
			Details: fmt.Sprintf("Conectado em %s", info.HarnessCore),
		}
	}
	return CheckItem{
		Name:    "Harness Core",
		Status:  "INFO",
		Details: "Modo Standalone (fonte central não detectada)",
	}
}
