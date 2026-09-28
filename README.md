# Antigravity Operator (`agyo`)

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-000000?style=for-the-badge&logo=apple&logoColor=white" alt="Platform" />
  <img src="https://img.shields.io/badge/Architecture-Single%20Binary%20(No%20CGO)-success?style=for-the-badge" alt="Binary" />
  <img src="https://img.shields.io/badge/Pattern-Fowler%20Outer%20Harness-blueviolet?style=for-the-badge" alt="Pattern" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License" />
</p>

> **The Autonomous Session Agent Engine & OS Runtime for Google Antigravity**  
> *Deterministic governance, filesystem operational memory, and Chrome DevTools isolation with seamless parity across macOS and Linux.*

<p align="center">
  <a href="README.pt-BR.md">🇧🇷 <b>Leia em Português</b></a> | <a href="#-getting-started--installation"><b>Getting Started</b></a> | <a href="#-student--google-ai-pro-edition"><b>Student Edition</b></a> | <a href="CONTRIBUTING.md"><b>Contributing</b></a>
</p>

> **Disclaimer:** *This is an open source community-driven companion project and is not an officially sponsored Google product. It is built to extend and empower the Google Antigravity & Google AI developer ecosystem.*

---

## 📑 Table of Contents

- [Overview](#-overview)
- [Tribute to the Community & Google AI Pro](#-tribute-to-the-community--google-ai-pro)
- [The Problem: Why Antigravity Needs an Operator](#-the-problem-why-antigravity-needs-an-operator)
- [The Solution: Core Capabilities](#-the-solution-what-antigravity-operator-solves)
- [Landscape & Benchmark](#-landscape-how-agyo-compares)
- [Student, Research & Google AI Pro Edition](#-student-research--google-ai-pro-edition)
- [System Architecture (SRP, KISS, YAGNI, DRY)](#-system-architecture-srp-kiss-yagni-dry)
- [Getting Started & Installation](#-getting-started--installation)
- [CLI Reference & Usage](#-cli-reference--usage)
- [Contributing & Community Standards](#-contributing--community-standards)
- [Security & License](#-security--license)

---

## 🔭 Overview

**Antigravity Operator** (`agyo`) turns the raw power of Google Antigravity into an autonomous, safe, and persistent **Operating System Operator** (Claude Computer Use / OS Agent style).

By implementing the canonical **Outer Harness (Martin Fowler)** model, `agyo` provides:
1. **Deterministic Session Memory:** State persists directly in `.agents/session/` on disk (`state.md`, `decisions.md`, `todo.md`), eliminating context amnesia.
2. **Zero-Pollution Chrome Isolation:** Automatically launches and supervises a dedicated Chrome instance on port `9222` (`~/.gemini/antigravity-browser-profile`), keeping your personal browsing safe and untouched.
3. **Headless & Server Linux Parity:** Automatically detects missing graphical environments (`$DISPLAY` / `$WAYLAND_DISPLAY`) and activates robust server flags (`--headless=new`, `--disable-dev-shm-usage`, `--no-sandbox`).
4. **Single-Binary Portability:** Written in pure Go with `CGO_ENABLED=0` and embedded templates (`//go:embed`), producing a self-contained ~6MB executable requiring zero dependencies.

---

## 🎁 Tribute to the Community & Google AI Pro

> *"This project is an open engineering contribution to the developer community, students, and researchers worldwide, and a special thank you to **Google** for the transformative student access program through **Google AI Pro**."*

Our mission is to democratize high-end agentic engineering: enabling every student and software engineer to leverage 100% of their **Gemini Pro and Antigravity** quotas with professional discipline, zero token waste, and seamless portability across any Linux or macOS machine.

---

## 🔍 The Problem: Why Antigravity Needs an Operator

Google Antigravity provides state-of-the-art atomic tooling: arbitrary bash execution, surgical file edits, subagents, and Model Context Protocol (MCP) integrations.

**However, out-of-the-box, it is a raw-power engine lacking an operational harness:**
* **Competitors bundle proprietary sandboxes:** Tools like Claude Code, Devin, or Cursor enforce pre-configured guardrails. Antigravity provides atomic tools (`run_command`, `write_to_file`), leaving session governance, persistence, and OS lifecycle to the developer.
* **Operational Amnesia:** Without deterministic filesystem state, agents lose context across compaction windows and repeated sessions.
* **Runtime Friction:** Developers must manually configure CDP ports, isolate browser profiles, and debug headless Linux edge cases.

### The 4 Critical Failure Modes in Local AI Agents:
1. **The "Drunken Agent" Syndrome:** Agents executing unchecked commands, hallucinating paths, assuming code works without testing, and entering infinite retry loops burning API quota.
2. **Operational Amnesia & Context Drift:** As token limits are reached, agents forget earlier architectural agreements and repeat solved mistakes.
3. **Personal Browser Hijacking:** Agents interacting with the web using the user's personal browser profile, risking banking cookies, private tabs, or crashes.
4. **The macOS vs Linux Chasm:** Scripts developed on macOS failing on Linux servers, VPSs, WSL2, or Docker due to missing graphical displays (`$DISPLAY`), `/dev/shm` memory constraints, or sandbox permission errors.

---

## 💡 The Solution: What `antigravity-operator` Solves

| Capability | Engineering Implementation |
|---|---|
| **Outer Harness (Fowler)** | **Guide × Sensor:** Deterministic directives guide the model; automated tests (`go test`, linters, runtime probes) validate every change before completion. |
| **Filesystem Memory** | **`.agents/session/`:** Real-time state (`state.md`), architecture log (`decisions.md`), and task tracker (`todo.md`) persist across chat resets. |
| **Browser Supervision** | **Chrome DevTools Protocol (CDP):** Dedicated profile on port `9222`, PID tracking, and graceful shutdown (`agyo browser stop`). |
| **Auto-Headless Mode** | **Dynamic Display Probe:** Injects `--headless=new`, `--disable-dev-shm-usage`, and `--no-sandbox` automatically in server environments. |
| **Zero Runtime Deps** | **Pure Go (`CGO_ENABLED=0`):** Single static ~6MB binary containing all embedded rules and templates. |

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

## 🎓 Student, Research & Google AI Pro Edition

For computer science students and researchers leveraging academic benefits such as **Google AI Pro**, `agyo` is the ultimate productivity multiplier:

1. **Token Quota Conservation:** Prevents infinite retry loops and verbose repetitive code outputs, ensuring your Gemini Pro quota lasts the entire semester.
2. **Zero-Root Portability in University Labs (Linux):** University labs often run locked-down Linux machines without `sudo` access to install Docker or system packages. The static `agyo-linux-amd64` binary runs directly from user space (`~/`).
3. **Academic Logbook & Portfolio:** The `.agents/session/` folder preserves architectural rationales and algorithm trade-offs, turning daily coding into documented learning logs.
4. **Safe Sandbox:** Isolated Chrome automation protects personal university credentials and institutional logins.

### 🎁 Bonus Student Skills Included (`skills/`):
This repository includes 3 canonical skills out-of-the-box:
* **`feynman-code-tutor`:** Senior tutor based on the Feynman Technique. Explains complex algorithms, data structures, and Big-O using real-world analogies and comprehension checkpoints.
* **`student-study-planner`:** Breaks down complex college syllabi, final projects, and technical interview prep into focused sprint cycles (20% theory, 80% deliberate coding).
* **`token-budget-guard`:** Surgical token optimizer ensuring context efficiency and zero repetitive code waste.

---

## 🏛️ System Architecture (SRP, KISS, YAGNI, DRY)

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
├── docs/                     # Architectural and cross-platform specifications
├── .github/workflows/        # Automated multi-OS CI (Ubuntu & macOS)
├── scripts/
│   ├── bootstrap.sh          # One-liner end-user setup
│   └── setup-dev.sh          # Developer setup with pre-commit hooks
├── AUTHORS                   # Project authors
├── CONTRIBUTORS              # Project contributors
├── CODE_OF_CONDUCT.md        # Google Open Source Community Guidelines
├── SECURITY.md               # Responsible vulnerability disclosure policy
├── CONTRIBUTING.md           # Contribution guide and testing protocol
└── Makefile                  # Native build and cross-compilation targets
```

---

## ⚡ Getting Started & Installation

### Prerequisites
- Go 1.22+ (to build from source) or download pre-built binary
- Google Chrome or Chromium installed on the host

### Build Locally
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

## 🚀 CLI Reference & Usage

### 1. Environment Diagnostics (`doctor`)
Inspects system readiness across OS, Git, Chrome, Node/NPX, and Harness connections:
```bash
agyo doctor
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

## 🤝 Contributing & Community Standards

We welcome contributions! Please review:
* [CONTRIBUTING.md](CONTRIBUTING.md) — Step-by-step contribution and testing workflow.
* [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — Adapted from Contributor Covenant and Google Open Source Guidelines.
* [AUTHORS](AUTHORS) & [CONTRIBUTORS](CONTRIBUTORS) — List of project maintainers and contributors.

---

## 🔒 Security & License

* **Security Policy:** Refer to [SECURITY.md](SECURITY.md) for vulnerability disclosure guidelines.
* **License:** Distributed under the [MIT](LICENSE) License.
