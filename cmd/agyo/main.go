package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tiagoboas/antigravity-operator/internal/doctor"
	"github.com/tiagoboas/antigravity-operator/internal/hook"
	"github.com/tiagoboas/antigravity-operator/internal/installer"
	"github.com/tiagoboas/antigravity-operator/internal/platform"
	"github.com/tiagoboas/antigravity-operator/internal/profile"
	"github.com/tiagoboas/antigravity-operator/internal/session"
	"github.com/tiagoboas/antigravity-operator/internal/watcher"
)

const Version = "0.3.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	info, err := platform.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		runInit(os.Args[2:])
	case "session":
		runSession(info, os.Args[2:])
	case "doctor":
		runDoctor(info)
	case "browser":
		runBrowser(info, os.Args[2:])
	case "sync":
		runSync(info)
	case "hook":
		runHook(os.Args[2:])
	case "about":

		printAbout()
	case "version", "-v", "--version":
		fmt.Printf("agyo (Antigravity Operator) v%s [%s/%s]\n", Version, info.OS, info.Arch)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`agyo - Antigravity Operator CLI (Autonomous Session Agent Engine)
The Official Community Outer Harness for Google Antigravity & Google AI Pro

Usage:
  agyo <command> [arguments]

Available commands:
  init [dir]          Scaffold operational memory (.agents/session/) in the target project
  session status      Display active session objective, status, and task completion metrics
  session archive     Archive completed session to historical log and reset templates
  session watch       Stream active agent reasoning and desktop notifications in real time
  doctor              Audit host readiness (OS, Chrome, DevTools 9222, Git, Node/NPX)
  browser start       Launch isolated Chrome instance with remote debugging flags
  browser status      Inspect DevTools port (9222) readiness and Chrome process PID
  browser stop        Gracefully terminate isolated Chrome process (SIGTERM)
  browser tabs        List all open tabs and target IDs in the isolated Chrome
  browser open <url>  Open a new tab at the given URL in the isolated Chrome
  browser close <id>  Close a specific tab by ID
  browser eval "<js>" Evaluate JavaScript expression in the active tab (pure Go CDP)
  browser shot [file] Capture PNG screenshot of active tab without external libraries
  sync                Synchronize canonical rules and MCP manifests to Google Antigravity
  hook install [dir]  Install git pre-commit hook to safeguard session continuity
  hook uninstall [dir] Remove agyo git pre-commit hook
  about               Display manifesto and tribute to the community & Google AI Pro
  version             Print version and system architecture`)
}

func printAbout() {
	fmt.Println(`================================================================================
  ANTIGRAVITY OPERATOR (agyo) — SESSION RUNTIME & OUTER HARNESS
================================================================================

This project is an open engineering tribute to the developer community, students,
and researchers worldwide, and a special thank you to Google for the transformative
student access program through Google AI Pro.

Mission:
Transform the raw atomic power of Google Antigravity into an autonomous, safe, and
persistent Session Operator with filesystem memory (.agents/session/) and seamless
parity across macOS and Linux — empowering every student and engineer to leverage
100% of their Gemini Pro quota without token waste or runtime friction.

Built with care, precision, and canonical software engineering (Martin Fowler Outer Harness).
================================================================================`)
}

func runInit(args []string) {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	force := initCmd.Bool("force", false, "Overwrite existing session files")
	_ = initCmd.Parse(args)

	targetDir := "."
	if initCmd.NArg() > 0 {
		targetDir = initCmd.Arg(0)
	}

	res, err := session.Init(targetDir, *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Session memory initialized at: %s\n", res.SessionDir)
	for _, c := range res.Created {
		fmt.Printf("   + Created: %s\n", c)
	}
	for _, s := range res.Skipped {
		fmt.Printf("   - Kept (already exists): %s\n", s)
	}
}

func runSession(info *platform.Info, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: agyo session [status|archive|watch] [dir]")
		os.Exit(1)
	}

	sub := args[0]
	targetDir := "."
	if len(args) > 1 {
		targetDir = args[1]
	}

	switch sub {
	case "status":
		sum, err := session.GetSummary(targetDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		pct := 0
		if sum.TotalTasks > 0 {
			pct = (sum.DoneTasks * 100) / sum.TotalTasks
		}

		fmt.Println("📋 Active Session Overview (.agents/session/)")
		fmt.Println("-----------------------------------------------------------------")
		fmt.Printf("🎯 Objective : %s\n", sum.Objective)
		fmt.Printf("⚡ Phase     : %s\n", sum.Status)
		fmt.Printf("📊 Progress  : %d/%d tasks completed (%d%%)\n", sum.DoneTasks, sum.TotalTasks, pct)
		if len(sum.Pending) > 0 {
			fmt.Println("\n⏳ Pending Next Steps:")
			for i, p := range sum.Pending {
				if i >= 5 {
					fmt.Printf("   ... and %d more\n", len(sum.Pending)-i)
					break
				}
				fmt.Printf("   - [ ] %s\n", p)
			}
		}
		fmt.Println("-----------------------------------------------------------------")

	case "archive":
		archiveFile, err := session.Archive(targetDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error archiving session: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("📦 Session archived successfully to:\n   %s\n", archiveFile)
		fmt.Println("✅ State and tasks reset with clean templates for the next task!")

	case "watch", "tail":
		runSessionWatch(info, args[1:])

	default:
		fmt.Fprintf(os.Stderr, "Unknown session subcommand: %s\n", sub)
		os.Exit(1)
	}
}

func runSessionWatch(info *platform.Info, args []string) {
	watchCmd := flag.NewFlagSet("session watch", flag.ExitOnError)
	once := watchCmd.Bool("once", false, "Exibe os passos recentes e encerra sem acompanhar em tempo real")
	notify := watchCmd.Bool("notify", true, "Emite notificação no SO quando o agente fizer uma pergunta")
	steps := watchCmd.Int("steps", 5, "Número de passos recentes para exibir inicialmente")
	_ = watchCmd.Parse(args)

	tInfo, err := watcher.FindActiveTranscript(info.GeminiDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro localizando transcrição: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📡 Streaming Antigravity Brain [%s]\n", tInfo.ConversationID)
	fmt.Printf("📄 Transcrição: %s\n", tInfo.Path)
	fmt.Println("-----------------------------------------------------------------")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	opts := watcher.WatchOptions{
		Follow:       !*once,
		NotifyOnWait: *notify,
		InitialSteps: *steps,
		OSName:       info.OS,
	}

	err = watcher.Stream(ctx, tInfo.Path, opts, func(evt *watcher.Event) {
		summary := evt.Summary()
		if summary != "" {
			fmt.Println(summary)
		}
	})

	if err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "Erro no stream da transcrição: %v\n", err)
		os.Exit(1)
	}

	if !*once {
		fmt.Println("\n🛑 Stream encerrado.")
	}
}

func runDoctor(info *platform.Info) {
	fmt.Printf("🔍 Antigravity Operator Doctor [OS: %s | Arch: %s]\n", info.OS, info.Arch)
	if info.HasDisplay {
		fmt.Println("🖥️  Display Server: Detected (Desktop GUI)")
	} else {
		fmt.Println("🖥️  Display Server: Not detected (Headless Mode Active)")
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
		fmt.Println("Usage: agyo browser [start|status|stop]")
		os.Exit(1)
	}

	sub := args[0]
	switch sub {
	case "status":
		st := profile.CheckStatus(info.BrowserProfile)
		if st.IsRunning {
			pidInfo := ""
			if st.PID > 0 {
				pidInfo = fmt.Sprintf(" [PID: %d]", st.PID)
			}
			fmt.Printf("✅ Chrome DevTools ACTIVE on port %d (%s)%s\n", st.Port, st.Version, pidInfo)
			fmt.Printf("   Isolated profile: %s\n", st.ProfileDir)
		} else {
			fmt.Printf("⚠️  Chrome DevTools INACTIVE on port %d\n", st.Port)
			fmt.Printf("   To launch, run: agyo browser start\n")
		}
	case "start":
		browserCmd := flag.NewFlagSet("browser start", flag.ExitOnError)
		headless := browserCmd.Bool("headless", false, "Force headless mode even with display available")
		port := browserCmd.Int("port", profile.DefaultDebugPort, "Chrome remote debugging port")
		_ = browserCmd.Parse(args[1:])

		fmt.Printf("🚀 Launching isolated Chrome for Google Antigravity on port %d...\n", *port)
		err := profile.Start(info, profile.StartOptions{
			Port:          *port,
			ForceHeadless: *headless,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error starting browser: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Isolated Chrome active on port %d!\n", *port)
	case "stop":
		fmt.Println("🛑 Terminating isolated Chrome instance...")
		err := profile.Stop(info)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning terminating browser: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Isolated Chrome terminated successfully.")

	case "tabs":
		tabs, err := profile.ListTabs(profile.DefaultDebugPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing browser tabs: %v\n", err)
			os.Exit(1)
		}
		if len(tabs) == 0 {
			fmt.Println("ℹ️  No open page tabs found in Chrome.")
			return
		}
		fmt.Printf("🌐 Open Chrome Tabs (%d active):\n", len(tabs))
		fmt.Println("-----------------------------------------------------------------")
		for _, t := range tabs {
			idPrefix := t.ID
			if len(idPrefix) > 8 {
				idPrefix = idPrefix[:8]
			}
			fmt.Printf("[%s] %-35s : %s\n", idPrefix, t.Title, t.URL)
		}
		fmt.Println("-----------------------------------------------------------------")

	case "open":
		if len(args) < 2 {
			fmt.Println("Usage: agyo browser open <url>")
			os.Exit(1)
		}
		targetURL := args[1]
		tab, err := profile.OpenTab(profile.DefaultDebugPort, targetURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening tab: %v\n", err)
			os.Exit(1)
		}
		idPrefix := tab.ID
		if len(idPrefix) > 8 {
			idPrefix = idPrefix[:8]
		}
		fmt.Printf("✅ Tab opened [%s]: %s\n", idPrefix, targetURL)

	case "close":
		if len(args) < 2 {
			fmt.Println("Usage: agyo browser close <tab-id>")
			os.Exit(1)
		}
		tabID := args[1]
		if err := profile.CloseTab(profile.DefaultDebugPort, tabID); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing tab: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Tab [%s] closed successfully.\n", tabID)

	case "eval":
		if len(args) < 2 {
			fmt.Println("Usage: agyo browser eval \"<javascript>\"")
			os.Exit(1)
		}
		expr := args[1]
		val, err := profile.Eval(profile.DefaultDebugPort, expr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error evaluating JS via CDP: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(val)

	case "shot", "screenshot":
		dest := "screenshot.png"
		if len(args) > 1 {
			dest = args[1]
		}
		if err := profile.Screenshot(profile.DefaultDebugPort, dest); err != nil {
			fmt.Fprintf(os.Stderr, "Error capturing screenshot: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("📸 Screenshot saved successfully to: %s\n", dest)

	default:
		fmt.Fprintf(os.Stderr, "Unknown browser subcommand: %s\n", sub)
		os.Exit(1)
	}
}


func runSync(info *platform.Info) {
	fmt.Println("🔄 Synchronizing rules and manifests into Google Antigravity...")
	res, err := installer.Sync(info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during synchronization: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Rules synchronized at: %s\n", res.RulesPath)
	if res.MCPPath != "" {
		fmt.Printf("✅ MCP manifest prepared at: %s\n", res.MCPPath)
	}
	if res.SkillsPath != "" {
		fmt.Printf("✅ Antigravity skill installed at: %s\n", res.SkillsPath)
	}
}

func runHook(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: agyo hook [install|uninstall] [dir]")
		os.Exit(1)
	}
	sub := args[0]
	targetDir := "."
	if len(args) > 1 {
		targetDir = args[1]
	}
	switch sub {
	case "install":
		path, err := hook.Install(targetDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error installing git hook: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Pre-commit session hook installed at: %s\n", path)
	case "uninstall":
		if err := hook.Uninstall(targetDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error uninstalling git hook: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Pre-commit session hook uninstalled successfully.")
	default:
		fmt.Fprintf(os.Stderr, "Unknown hook subcommand: %s\n", sub)
		os.Exit(1)
	}
}

