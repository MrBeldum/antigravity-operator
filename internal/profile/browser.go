package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/tiagoboas/antigravity-operator/internal/platform"
)

const (
	DefaultDebugPort = 9222
	VersionEndpoint  = "http://127.0.0.1:9222/json/version"
)

// ChromeVersionResponse reflete a resposta do endpoint /json/version do Chrome.
type ChromeVersionResponse struct {
	Browser              string `json:"Browser"`
	ProtocolVersion      string `json:"Protocol-Version"`
	UserAgent            string `json:"User-Agent"`
	V8Version            string `json:"V8-Version"`
	WebKitVersion        string `json:"WebKit-Version"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

// Status indica a situação atual do navegador isolado.
type Status struct {
	IsRunning  bool
	Port       int
	Version    string
	ProfileDir string
	Headless   bool
}

// CheckStatus verifica se a porta 9222 está ativa e respondendo com a API DevTools.
func CheckStatus(profileDir string) Status {
	st := Status{
		Port:       DefaultDebugPort,
		ProfileDir: profileDir,
	}

	client := http.Client{Timeout: 800 * time.Millisecond}
	resp, err := client.Get(VersionEndpoint)
	if err != nil {
		st.IsRunning = false
		return st
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var ver ChromeVersionResponse
		if err := json.NewDecoder(resp.Body).Decode(&ver); err == nil {
			st.IsRunning = true
			st.Version = ver.Browser
		}
	}

	return st
}

// StartOptions configura a inicialização do Chrome isolado.
type StartOptions struct {
	ForceHeadless bool
}

// Start inicializa o Chrome com o perfil isolado e flags corretas para a plataforma.
func Start(info *platform.Info, opts StartOptions) error {
	if info.ChromeBin == "" {
		return fmt.Errorf("binário do Google Chrome / Chromium não encontrado no sistema")
	}

	// 1. Assegurar que o diretório de perfil existe
	if err := os.MkdirAll(info.BrowserProfile, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta de perfil isolado: %w", err)
	}

	// 2. Se já estiver rodando, nada a fazer
	current := CheckStatus(info.BrowserProfile)
	if current.IsRunning {
		return nil // já está ativo e operando
	}

	// 3. Montar flags
	args := []string{
		fmt.Sprintf("--user-data-dir=%s", info.BrowserProfile),
		fmt.Sprintf("--remote-debugging-port=%d", DefaultDebugPort),
		"--no-first-run",
		"--no-default-browser-check",
	}

	// Decisão de Headless: forçado por flag ou obrigatório por falta de display gráfico no Linux
	needHeadless := opts.ForceHeadless || (!info.HasDisplay && info.OS == "linux")
	if needHeadless {
		args = append(args,
			"--headless=new",
			"--disable-gpu",
			"--disable-dev-shm-usage", // Crítico para Linux / Docker / VPS
			"--no-sandbox",
		)
	}

	cmd := exec.Command(info.ChromeBin, args...)
	// Desacoplar processo para rodar em background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("falha ao iniciar processo do Chrome: %w", err)
	}

	// 4. Aguardar até 5 segundos para o endpoint responder
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("tempo limite excedido aguardando Chrome responder na porta %d", DefaultDebugPort)
		default:
			if CheckStatus(info.BrowserProfile).IsRunning {
				return nil
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}
