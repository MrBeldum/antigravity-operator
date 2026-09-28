package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tiagoboas/antigravity-operator/internal/doctor"
	"github.com/tiagoboas/antigravity-operator/internal/installer"
	"github.com/tiagoboas/antigravity-operator/internal/platform"
	"github.com/tiagoboas/antigravity-operator/internal/profile"
	"github.com/tiagoboas/antigravity-operator/internal/session"
)

const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	info, err := platform.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao detectar plataforma: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		runInit(os.Args[2:])
	case "doctor":
		runDoctor(info)
	case "browser":
		runBrowser(info, os.Args[2:])
	case "sync":
		runSync(info)
	case "version", "-v", "--version":
		fmt.Printf("agyo (Antigravity Operator) v%s [%s/%s]\n", Version, info.OS, info.Arch)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`agyo - Antigravity Operator CLI (Autonomous Session Agent Engine)

Uso:
  agyo <comando> [argumentos]

Comandos disponíveis:
  init [dir]          Inicializa a memória operacional (.agents/session/) no projeto
  doctor              Audita a saúde da máquina (OS, Chrome, DevTools 9222, Git, Node)
  browser start       Inicia o Chrome isolado com remote debugging e flags corretas
  browser status      Verifica a integridade da porta DevTools (9222)
  sync                Sincroniza regras globais e MCPs canônicos no Antigravity
  version             Exibe a versão do operador e arquitetura do sistema`)
}

func runInit(args []string) {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	force := initCmd.Bool("force", false, "Sobrescrever arquivos existentes de sessão")
	_ = initCmd.Parse(args)

	targetDir := "."
	if initCmd.NArg() > 0 {
		targetDir = initCmd.Arg(0)
	}

	res, err := session.Init(targetDir, *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar sessão: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Memória operacional inicializada em: %s\n", res.SessionDir)
	for _, c := range res.Created {
		fmt.Printf("   + Criado: %s\n", c)
	}
	for _, s := range res.Skipped {
		fmt.Printf("   - Mantido (existente): %s\n", s)
	}
}

func runDoctor(info *platform.Info) {
	fmt.Printf("🔍 Antigravity Operator Doctor [SO: %s | Arch: %s]\n", info.OS, info.Arch)
	if info.HasDisplay {
		fmt.Println("🖥️  Ambiente Gráfico: Detectado (Desktop GUI)")
	} else {
		fmt.Println("🖥️  Ambiente Gráfico: Não detectado (Headless Mode ativo)")
	}
	fmt.Println("-----------------------------------------------------------------")

	report := doctor.Run(info)
	for _, chk := range report.Checks {
		var icon string
		switch chk.Status {
		case "OK":
			icon = "✅"
		case "WARN":
			icon = "⚠️ "
		case "FAIL":
			icon = "❌"
		default:
			icon = "ℹ️ "
		}
		fmt.Printf("%s %-28s : %s\n", icon, chk.Name, chk.Details)
	}
	fmt.Println("-----------------------------------------------------------------")
}

func runBrowser(info *platform.Info, args []string) {
	if len(args) == 0 {
		fmt.Println("Uso: agyo browser [start|status]")
		os.Exit(1)
	}

	sub := args[0]
	switch sub {
	case "status":
		st := profile.CheckStatus(info.BrowserProfile)
		if st.IsRunning {
			fmt.Printf("✅ Chrome DevTools ATIVO na porta %d (%s)\n", st.Port, st.Version)
			fmt.Printf("   Perfil isolado: %s\n", st.ProfileDir)
		} else {
			fmt.Printf("⚠️  Chrome DevTools INATIVO na porta %d\n", st.Port)
			fmt.Printf("   Para iniciar, execute: agyo browser start\n")
		}
	case "start":
		browserCmd := flag.NewFlagSet("browser start", flag.ExitOnError)
		headless := browserCmd.Bool("headless", false, "Forçar modo headless mesmo com display")
		_ = browserCmd.Parse(args[1:])

		fmt.Println("🚀 Iniciando Chrome isolado para o Session Agent...")
		err := profile.Start(info, profile.StartOptions{
			ForceHeadless: *headless,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao iniciar browser: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Chrome isolado ativo com sucesso na porta 9222!")
	default:
		fmt.Fprintf(os.Stderr, "Subcomando de browser desconhecido: %s\n", sub)
		os.Exit(1)
	}
}

func runSync(info *platform.Info) {
	fmt.Println("🔄 Sincronizando regras e manifestos do Session Agent...")
	res, err := installer.Sync(info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro na sincronização: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Regras sincronizadas em: %s\n", res.RulesPath)
	if res.MCPPath != "" {
		fmt.Printf("✅ Manifesto MCP preparado em: %s\n", res.MCPPath)
	}
}
