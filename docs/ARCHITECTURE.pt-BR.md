# Arquitetura do Antigravity Operator (`agyo`)

## 1. Visão Geral do Sistema

O `antigravity-operator` opera como uma camada de runtime e governança entre o **Modelo de IA** (como Gemini ou Claude no Antigravity) e o **Sistema Operacional do Desenvolvedor** (macOS ou Linux).

```mermaid
graph TD
    User([Usuário / Sessão]) -->|Prompt| Agent[Antigravity Session Agent]
    
    subgraph Governance ["Outer Harness (Fowler)"]
        Agent -->|1. Consulta Guia| Rules[templates/rules/session-agent.md]
        Agent -->|2. Persiste Estado| Memory[".agents/session/{state, decisions, todo}.md"]
    end

    subgraph Runtime ["Agyo OS Engine (Go)"]
        Agent -->|3. Executa Ações| CLI[agyo CLI]
        CLI --> Platform[internal/platform]
        CLI --> Doctor[internal/doctor]
        CLI --> Profile[internal/profile]
        CLI --> Installer[internal/installer]
        CLI --> CDP[internal/cdp]
        CLI --> Hook[internal/hook]
        CLI --> Watcher[internal/watcher]
    end

    subgraph OS_Targets ["Alvos do Sistema Operacional"]
        Profile -->|CDP Port 9222 / PID| IsolatedChrome["Chrome Isolado (~/.gemini/antigravity-browser-profile)"]
        CDP -->|Pure Go RFC 6455 WS| IsolatedChrome
        Watcher -->|Tail JSONL & Alerta| Brain["Antigravity Brain (~/.gemini/antigravity/brain/*/transcript.jsonl)"]
        Hook -->|Pre-Commit Gate| Filesystem["Filesystem & Git Repo"]
        Installer -->|MCPs| DevToolsMCP["Chrome DevTools MCP & Playwright"]
    end
```

---

## 2. Ciclo de Vida da Sessão

1. **Inicialização (`agyo init`):**
   * Cria o diretório `.agents/session/`.
   * Cria os arquivos `state.md`, `decisions.md` e `todo.md` com templates canônicos caso não existam.
   * Cria `.agents/.gitignore` para impedir que credenciais e dados temporários de depuração sejam commitados acidentalmente.

2. **Diagnóstico da Máquina (`agyo doctor`):**
   * Avalia a saúde da máquina local: Git, Chrome, Node/NPX, exibição gráfica, processos do Antigravity, chaves BYOK e integração com o `harness-core`.
   * Classifica cada item como `OK`, `WARN`, `FAIL` ou `INFO`.

3. **Orquestração de Navegação e CDP Nativo (`agyo browser`):**
   * Lança o Google Chrome restrito ao diretório `~/.gemini/antigravity-browser-profile`.
   * Abre a porta de depuração remota `9222`.
   * Se executado em ambiente Linux sem servidor gráfico (`$DISPLAY`), automaticamente anexa as flags headless essenciais (`--headless=new`, `--disable-dev-shm-usage`, `--no-sandbox`).
   * Faz polling no endpoint `http://127.0.0.1:9222/json/version` até receber status 200 OK.
   * Fornece subcomandos nativos de inspeção CDP em Go puro (`tabs`, `open`, `close`, `eval`, `shot`) via WebSockets RFC 6455 sem dependências externas.

4. **Streaming de Raciocínio e Alertas Desktop (`agyo session watch`):**
   * Monitora em tempo real a transcrição da sessão ativa (`~/.gemini/antigravity/brain/*/transcript.jsonl`).
   * Formata pensamentos (`💭 Think`), chamadas de ferramentas (`🛠️ Tool`) e interação do usuário (`👤 User`).
   * Dispara alertas visuais e sonoros no sistema operacional (`osascript` no macOS, `notify-send` no Linux, bell terminal `\a`) quando o agente solicita resposta (`ask_question`).

5. **Hook de Continuidade Pré-Commit (`agyo hook`):**
   * Instala um sensor do outer harness em `.git/hooks/pre-commit`.
   * Bloqueia commits caso `.agents/session/state.md` e `todo.md` não tenham sido atualizados na sessão.

6. **Sincronização de Regras (`agyo sync`):**
   * Grava as regras canônicas do Session Agent em `~/.gemini/antigravity/rules/session-agent.md`.
   * Prepara os manifestos padrão de MCPs em `~/.gemini/antigravity/mcp/default-servers.json`.

---

## 3. Padrões de Projeto e Decisões de Engenharia

- **Single Responsibility Principle (SRP):** Cada pacote sob `internal/` possui um escopo estrito e não vaza detalhes de implementação para outros pacotes.
- **Embed Nativo (`//go:embed`):** Permite distribuição de binário único sem instaladores complexos ou necessidade de clonar o repositório em todas as máquinas.
- **Zero CGO (`CGO_ENABLED=0`):** Garante compatibilidade binária entre qualquer versão de kernel Linux e biblioteca C (glibc ou musl).
- **Zero Dependências Externas em Tempo de Execução:** Todo o código de rede, WebSockets (RFC 6455) e streaming de JSONL é implementado diretamente sobre a biblioteca padrão do Go.
