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

const Version = "0.2.0"

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
	case "doctor":
		runDoctor(info)
	case "browser":
		runBrowser(info, os.Args[2:])
	case "sync":
		runSync(info)
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
  doctor              Audit host readiness (OS, Chrome, DevTools 9222, Git, Node/NPX)
  browser start       Launch isolated Chrome instance with remote debugging flags
  browser status      Inspect DevTools port (9222) readiness and Chrome process PID
  browser stop        Gracefully terminate isolated Chrome process (SIGTERM)
  sync                Synchronize canonical rules and MCP manifests to Google Antigravity
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
}
