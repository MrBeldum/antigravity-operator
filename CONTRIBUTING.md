# Contributing to Antigravity Operator (`agyo`)

Thank you for your interest in contributing to **Antigravity Operator**! This project is an open-source initiative designed to bring a robust, production-grade **Outer Harness** and session persistence engine to Google Antigravity and the broader Google AI developer ecosystem.

---

## 🧭 Code of Conduct & Core Philosophy

We follow the **Martin Fowler Outer Harness** model:
1. **Guide × Sensor:** Provide deterministic instructions (Guides) and automated tests (Sensors) before shipping.
2. **KISS & YAGNI:** Keep it simple, do not over-engineer. Avoid premature abstractions and unnecessary dependencies.
3. **Zero CGO:** Everything must compile statically (`CGO_ENABLED=0`) across Darwin (macOS) and Linux.
4. **No Assumptions:** Write code backed by real OS inspections and verified unit tests.

---

## 🛠️ Development Setup

### Prerequisites
- Go 1.22+
- Git
- Google Chrome or Chromium (optional for runtime browser testing)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/tiagovilasboas/antigravity-operator.git
cd antigravity-operator

# Run the automated development setup (configures pre-commit hooks and builds the binary)
./scripts/setup-dev.sh
```

---

## 🧪 Testing Guidelines

Before opening a pull request, ensure all sensors pass:

```bash
# Static analysis
go vet ./...

# Unit tests with race detector
go test -v -race ./...

# Cross-compilation verification for Linux
make build-linux

# Run doctor smoke test
./bin/agyo doctor
```

---

## 📝 Commit Conventions

We strictly follow **Conventional Commits** in English:
- `feat(...)`: A new feature
- `fix(...)`: A bug fix
- `docs(...)`: Documentation changes
- `test(...)`: Adding or refactoring tests
- `refactor(...)`: Code change that neither fixes a bug nor adds a feature

Example:
```bash
git commit -m "feat(browser): support dynamic CDP port detection"
```

---

## 🚀 Submitting a Pull Request

1. Fork the repository.
2. Create your feature branch (`git checkout -b feat/my-new-feature`).
3. Commit your changes adhering to the commit guidelines.
4. Push to the branch (`git push origin feat/my-new-feature`).
5. Open a Pull Request with a clear summary of changes and test evidence.

Thank you for helping empower students, researchers, and developers worldwide! 🎓
