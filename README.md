# Antigravity Operator (`agyo`)

> **Autonomous Session Agent Engine & OS Runtime**  
> Empacota, padroniza e provisiona o Antigravity como um **Session Agent autônomo** (estilo Claude Computer Use / OS Agent) com governança rígida, suporte a múltiplos sistemas (macOS e Linux) e zero dependências externas.

---

## 🎯 Por que o Antigravity Operator?

A maioria das soluções de agentes locais comete um de dois erros:
1. **Agente solto sem governança:** Executa bash sem limites, alucina dados, assume estados sem checar e entra em loops infinitos.
2. **Ambiente não reprodutível:** Funciona na máquina de um desenvolvedor, mas quebra no Linux, VPS, WSL ou sem display gráfico.

O **`antigravity-operator`** implementa o modelo canônico do **Harness (Martin Fowler)**:
- **Outer Harness (Guia × Sensor):** Regras claras antes da execução e sensores computacionais de validação contínua.
- **Sem Afirmação Sem Fonte:** Nenhuma suposição sem teste, log ou evidência verificada.
- **Memória Operacional de Sessão (`.agents/session/`):** O estado da tarefa persiste no disco (`state.md`, `decisions.md`, `todo.md`) em vez de depender da janela volátil de tokens.
- **Browser & DevTools com Perfil Isolado:** Chrome DevTools MCP opera em porta de debug (`9222`) com perfil dedicado (`~/.gemini/antigravity-browser-profile`), sem misturar com sua navegação pessoal e adaptando-se a ambientes com ou sem interface gráfica (Linux Headless automático).
- **Single-Binary Portability:** Compilado em Go puro (`CGO_ENABLED=0`), gerando um único executável estático que roda idêntico no macOS (Apple Silicon/Intel) e Linux (x86_64/ARM64).

---

## 🏗️ Arquitetura Canônica (SRP, KISS, YAGNI, DRY)

```text
antigravity-operator/
├── cmd/agyo/                 # Entrypoint da CLI
├── internal/
│   ├── platform/             # SRP: Detecção de SO, display X11/Wayland e caminhos do Chrome
│   ├── session/              # SRP: Scaffold da memória operacional (.agents/session/)
│   ├── profile/              # SRP: Gerenciamento do Chrome DevTools e remote debugging
│   ├── installer/            # SRP: Sincronização idempotente de regras e manifestos MCP
│   └── doctor/               # SRP: Sensor computacional de diagnóstico da máquina
├── templates/                # Embutido no binário estático via //go:embed
│   ├── rules/                # Regras canônicas de Session Agent
│   ├── session/              # Templates de state.md, decisions.md e todo.md
│   └── mcps/                 # Manifesto de servidores MCP (DevTools, Playwright)
├── agents/                   # Personas especializadas para engenharia assistida por IA
├── docs/                     # Documentação de arquitetura e matriz cross-platform
├── .github/workflows/        # CI automatizado de paridade macOS & Linux
├── scripts/
│   ├── bootstrap.sh          # Setup para usuários finais
│   └── setup-dev.sh          # Setup de desenvolvimento e pre-commit hooks
└── Makefile                  # Build nativo e cross-compilação para Linux
```

---

## 🤖 Assisted-IA & Ecossistema de Agentes

O projeto foi desenhado sob o modelo **Agent-as-Code** e contém governança nativa para agentes:
- **`AGENTS.md`:** Contrato de conduta, diretrizes de código Go, checklist de sensores e convenções de commit para qualquer IA (Antigravity, Cursor, Claude, Copilot).
- **Roster de Especialistas (`agents/`):**
  - **`operator-architect`:** Guardião do sistema operacional, paridade macOS/Linux e princípios KISS/YAGNI.
  - **`cdp-engineer`:** Especialista no Chrome DevTools Protocol, flags de browser e sockets de depuração.
  - **`qa-sentinel`:** Responsável pelos testes automatizados e sensores de regressão.

---

## ⚡ Instalação Rápida

### Compilar localmente (Go 1.22+)
```bash
git clone https://github.com/tiagoboas/antigravity-operator.git
cd antigravity-operator
make build
```

O binário estará disponível em `bin/agyo`. Para instalar no seu `PATH`:
```bash
make install
```

### Cross-compilar para Linux
```bash
make build-linux
# Binários estáticos gerados em bin/agyo-linux-amd64 e bin/agyo-linux-arm64
```

---

## 🚀 Como Usar

### 1. Diagnosticar o ambiente da máquina (`doctor`)
Audita se o sistema operacional, Git, Chrome, Node/NPX e conexões de harness estão prontos:
```bash
agyo doctor
```

Exemplo de saída:
```text
🔍 Antigravity Operator Doctor [SO: darwin | Arch: arm64]
🖥️  Ambiente Gráfico: Detectado (Desktop GUI)
-----------------------------------------------------------------
✅ Git                          : git version 2.39.5 (Tiago Vilas Boas <tcarvalhovb@gmail.com>)
✅ Google Chrome                : Localizado em: /Applications/Google Chrome.app/Contents/MacOS/Google Chrome
⚠️  Chrome DevTools (Port 9222)  : Inativo (execute 'agyo browser start' para iniciar)
✅ NPX (MCP Runtime)            : Versão 10.8.2 disponível
✅ Harness Core                 : Conectado em /Users/tiago.boas/Github/harness-core
-----------------------------------------------------------------
```

### 2. Inicializar a memória operacional no projeto atual (`init`)
Cria a pasta `.agents/session/` com rastreamento de estado, decisões e tarefas:
```bash
cd meu-projeto
agyo init
```

Estrutura criada:
- `.agents/session/state.md`
- `.agents/session/decisions.md`
- `.agents/session/todo.md`
- `.agents/.gitignore` (protegendo logs e credenciais contra commits acidentais)

### 3. Iniciar o Chrome isolado para o DevTools MCP (`browser start`)
Sobe uma instância exclusiva do Chrome para o agente inspecionar DOM, rede e console:
```bash
# Iniciar normalmente (abre janela no Mac/Linux desktop):
agyo browser start

# Forçar modo headless (útil para servidores, CI ou WSL2 sem tela):
agyo browser start --headless

# Verificar se a porta de debug está ativa:
agyo browser status
```

### 4. Sincronizar regras e MCPs no Antigravity (`sync`)
Garante que as regras de governança e servidores de automação estejam instalados:
```bash
agyo sync
```

---

## 🛡️ Princípios Operacionais do Agente
Quando o Antigravity opera sob o `agyo`, ele segue 5 mandamentos:
1. **Autonomia de Investigação:** Busca fatos no terminal, browser e logs antes de fazer perguntas triviais.
2. **Orquestração Multiferramenta:** Identifica -> Investiga -> Implementa -> Testa -> Valida no Browser.
3. **Validação Rigorosa:** A tarefa só termina quando o resultado foi validado de ponta a ponta com evidências.
4. **Perfil Isolado:** Zero interferência ou exposição no Chrome pessoal do usuário.
5. **Comunicação Concisa:** Direta ao ponto, técnica e fundamentada em dados.
