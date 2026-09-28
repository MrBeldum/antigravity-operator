# Antigravity Operator (`agyo`)

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-000000?style=for-the-badge&logo=apple&logoColor=white" alt="Platform" />
  <img src="https://img.shields.io/badge/Architecture-Single%20Binary%20(No%20CGO)-success?style=for-the-badge" alt="Binary" />
  <img src="https://img.shields.io/badge/Pattern-Fowler%20Outer%20Harness-blueviolet?style=for-the-badge" alt="Pattern" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License" />
</p>

> **The Autonomous Session Agent Engine & OS Runtime**  
> *Transform Google Antigravity and AI coding agents into autonomous, safe, and persistent operating system operators (Claude Computer Use / OS Agent style) with deterministic governance, filesystem memory, and total parity across macOS and Linux.*

<p align="center">
  <a href="README.pt-BR.md">🇧🇷 <b>Leia em Português</b></a> | <a href="#-quick-start"><b>Quick Start</b></a> | <a href="#-student--google-ai-pro-edition"><b>Student Edition</b></a> | <a href="CONTRIBUTING.md"><b>Contributing</b></a>
</p>

---

### 🎁 An Engineering Tribute to the Community & Google

> *"This project is an open contribution to developers, students, and researchers worldwide, and a special tribute to **Google** for the transformative student access program through **Google AI Pro**."*  
>  
> Our goal is to democratize high-end agentic engineering: enabling every student and software engineer to leverage 100% of the power of **Google Antigravity and Gemini Pro** models with professional rigour, zero token waste, and seamless portability across any Linux or macOS environment.

---

## 🔍 The Problem: Why Google Antigravity Needs an Operator

Google Antigravity is one of the most powerful AI coding engines available today — featuring atomic tools for shell execution, surgical code edits, subagents, and Model Context Protocol (MCP) integrations.

**However, out of the box, Antigravity is a raw power engine without a built-in Session Harness:**
* **Competitors bundle proprietary sandboxes:** Tools like Claude Code, Devin, or Cursor enforce pre-configured guardrails. Antigravity provides atomic tools (`run_command`, `write_to_file`), but leaves session governance, persistence, and OS lifecycle entirely to the user.
* **Lack of Long-Term Memory:** Sessions lack deterministic filesystem persistence, causing operational amnesia across context window compactions.
* **Runtime Friction:** Developers have to manually configure CDP ports, browser profile isolation, and diagnose headless Linux edge cases.

### The 4 Critical Failure Modes in Local AI Agents:

1. **The "Drunken Agent" Syndrome:** Agents with bash access running wild commands, hallucinating paths, assuming code works without testing, and entering infinite retry loops burning tokens.
2. **Operational Amnesia & Context Drift:** As discussions grow, agents forget previous architectural agreements and repeat solved mistakes.
3. **Personal Browser Hijacking & Security Risks:** Agents interacting with web pages by hijacking the developer's personal Chrome profile, exposing sensitive cookies or crashing active tabs.
4. **The macOS vs Linux Chasm:** Automation scripts developed on macOS failing on Linux servers, VPSs, WSL2, or Docker due to missing graphical displays (`$DISPLAY`), `/dev/shm` memory constraints, or sandbox permission errors.

---

## 💡 The Solution: What `antigravity-operator` Solves

The **`antigravity-operator`** (`agyo`) wraps operating system infrastructure and governance around Antigravity:

* **Martin Fowler Outer Harness (Guide × Sensor):** The agent operates strictly with no unverified assumptions. Deterministic rules guide the model before generation; computational sensors (`go test`, linters, runtime probes) validate every change before completion.
* **Persistent Session Memory (`.agents/session/`):** State transitions are versioned directly in the project filesystem (`state.md`, `decisions.md`, `todo.md`). Memory survives context compactions and IDE restarts.
* **Total Chrome Isolation via DevTools MCP:** Automatically manages an isolated Google Chrome instance (`~/.gemini/antigravity-browser-profile`) on debug port `9222`, keeping your personal browsing completely untouched.
* **Dynamic Headless Mode (Linux & Servers):** Intelligently detects graphical displays (`$DISPLAY` / `$WAYLAND_DISPLAY`). If headless, it seamlessly injects `--headless=new`, `--disable-dev-shm-usage`, and `--no-sandbox`.
* **Zero Runtime Dependencies (Single Binary Go):** Pure Go (`CGO_ENABLED=0`), generating a single static ~6MB binary that runs instantly on any macOS (Apple Silicon / Intel) or Linux (x86_64 / ARM64) distribution.

---

## 🥊 Landscape: How `agyo` Compares

| Feature | Raw Bash Scripts | Claude Computer Use | Open-Interpreter | **Antigravity Operator (`agyo`)** |
|---|---|---|---|---|
| **Outer Harness Governance** | ❌ No | ❌ No | ❌ No | **✅ Native (Guide × Sensor)** |
| **Filesystem Session Memory** | ❌ No | ❌ No | ❌ No | **✅ Canonical `.agents/session/`** |
| **Isolated Browser Profile** | ❌ Uses personal | ⚠️ Heavy Docker | ❌ No | **✅ Dedicated Profile (`9222`)** |
| **macOS / Linux Parity** | ⚠️ Fragile | ⚠️ Docker-only | ⚠️ Dep conflicts | **✅ Native & Auto-Headless** |
| **Runtime Footprint** | Multi-tooling | Docker / APIs | Python / venv | **✅ Single Static Binary (~6MB)** |
| **Integrated Diagnostics (`doctor`)** | ❌ No | ❌ No | ❌ No | **✅ Built into CLI** |

---

## 🎓 Spotlight: Students, Researchers & Google AI Pro

For computer science students and researchers benefiting from academic programs like **Google AI Pro**, `agyo` is the ultimate productivity multiplier:

1. **Token Quota Conservation:** Prevents infinite loops and verbose repetitive code outputs, ensuring your Gemini Pro quota lasts the entire semester.
2. **Zero-Root Portability in University Labs (Linux):** University labs often run locked-down Linux machines where students cannot install Docker or global packages. The static `agyo-linux-amd64` binary runs directly from user space (`~/`).
3. **Academic Logbook & Portfolio:** The `.agents/session/` folder preserves architectural rationales and algorithm trade-offs, turning daily coding into documented learning logs.
4. **Safe Sandbox:** Isolated Chrome automation protects personal university credentials and institutional logins.

### 🎁 Bonus Student Skills Included (`skills/`):
This repository includes 3 canonical skills out-of-the-box:
* **`feynman-code-tutor`:** Senior tutor based on the Feynman Technique. Explains complex algorithms, data structures, and Big-O using real-world analogies and comprehension checkpoints.
* **`student-study-planner`:** Breaks down complex college syllabi, final projects, and technical interview prep into focused sprint cycles (20% theory, 80% deliberate coding).
* **`token-budget-guard`:** Surgical token optimizer ensuring context efficiency and zero repetitive code waste.

---

## 🏛️ Architecture (SRP, KISS, YAGNI, DRY)

```text
antigravity-operator/
├── cmd/agyo/                 # CLI entrypoint (main.go)
├── internal/
│   ├── platform/             # SRP: OS detection, X11/Wayland check, Chrome binary resolution
│   ├── session/              # SRP: .agents/session/ scaffold & gitignore protection
│   ├── profile/              # SRP: Chrome lifecycle management, PID tracking & CDP port
│   ├── installer/            # SRP: Idempotent rule and MCP manifesto synchronization
│   └── doctor/               # SRP: Machine diagnostic computational sensors
├── templates/                # Embedded static assets via //go:embed (zero external deps)
│   ├── rules/                # Canonical Session Agent rules
│   ├── session/              # Templates for state.md, decisions.md, and todo.md
│   └── mcps/                 # Default MCP servers manifest (DevTools, Playwright)
├── agents/                   # Specialized AI personas (operator-architect, cdp-engineer, qa-sentinel)
├── skills/                   # Bonus student skills (feynman tutor, study planner, token guard)
├── docs/                     # In-depth architectural & cross-platform specs
├── .github/workflows/        # Automated multi-OS CI (Ubuntu & macOS)
├── scripts/
│   ├── bootstrap.sh          # One-liner end-user setup
│   └── setup-dev.sh          # Developer setup with pre-commit hooks
└── Makefile                  # Native build and cross-compilation targets
```

---

## ⚡ Quick Start

### Build Locally (Go 1.22+)
```bash
git clone https://github.com/tiagovilasboas/antigravity-operator.git
cd antigravity-operator
make build
```

The compiled binary will be placed at `bin/agyo` (with an `antigravity-operator` symlink). To install system-wide:
```bash
make install
```

### Cross-Compile for Linux from macOS
```bash
make build-linux
# Static binaries generated at bin/agyo-linux-amd64 and bin/agyo-linux-arm64
```

---

## 🚀 CLI Usage

### 1. Environment Diagnostics (`doctor`)
Inspects system readiness across OS, Git, Chrome, Node/NPX, and Harness connections:
```bash
agyo doctor
```

Example output:
```text
🔍 Antigravity Operator Doctor [SO: darwin | Arch: arm64]
🖥️  Display Environment: Detected (Desktop GUI)
-----------------------------------------------------------------
✅ Git                          : git version 2.39.5 (Apple Git-154) (Tiago Vilas Boas <tcarvalhovb@gmail.com>)
✅ Google Chrome                : Found at: /Applications/Google Chrome.app/Contents/MacOS/Google Chrome
⚠️  Chrome DevTools (Port 9222)  : Inactive (run 'agyo browser start' to launch)
✅ NPX (MCP Runtime)            : Version 10.8.2 available
✅ Harness Core                 : Connected at /Users/tiago.boas/Github/harness-core
-----------------------------------------------------------------
```

### 2. Scaffold Operational Memory (`init`)
Creates the `.agents/session/` memory structure in your current project:
```bash
cd my-project
agyo init
```

Files created:
- `.agents/session/state.md` (Real-time objective and status)
- `.agents/session/decisions.md` (Architecture log and trade-offs)
- `.agents/session/todo.md` (Task tracker)
- `.agents/.gitignore` (Protects runtime logs and sensitive credentials)

### 3. Manage Isolated Chrome Lifecycle (`browser`)
Full process supervision with PID tracking and graceful shutdown:
```bash
# Launch isolated Chrome on port 9222 (Desktop GUI):
agyo browser start

# Launch on a custom port:
agyo browser start --port 9223

# Force headless mode (automatic on headless Linux servers, VPS, or WSL2):
agyo browser start --headless

# Check status and PID:
agyo browser status

# Gracefully terminate isolated Chrome instance (SIGTERM):
agyo browser stop
```

### 4. Sync Rules and MCP Manifestos (`sync`)
Provisions canonical rules and automation manifests into Google Antigravity:
```bash
agyo sync
```

### 5. Project Manifesto (`about`)
```bash
agyo about
```

---

## 🛡️ Canonical Operational Rules

When operating under `agyo`, agents adhere to 5 core rules:
1. **Autonomous Investigation:** Seek ground truth via shell, browser, and logs before asking trivia.
2. **Multi-Tool Orchestration:** Identify -> Investigate -> Implement -> Test -> Validate in Browser.
3. **Rigorous Validation:** Tasks are complete only after verified end-to-end evidence.
4. **Profile Isolation:** Never compromise or access the user's personal browser profile.
5. **Concise Communication:** Direct, technical, and grounded in empirical facts.

---

## 📄 License
Distributed under the [MIT](LICENSE) License.
