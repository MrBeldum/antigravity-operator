# Decisões de Arquitetura e Sessão

## Premissas Acordadas
- **Stack:** Go 1.22+ com compilação 100% estática (`CGO_ENABLED=0`), zero dependências de runtime.
- **Princípios Canônicos:** Outer Harness (Guia × Sensor de Martin Fowler), SRP (pacotes cirúrgicos: `platform`, `session`, `profile`, `installer`, `doctor`), KISS, YAGNI, DRY, Sem Afirmação Sem Fonte.
- **Portabilidade:** Suporte nativo e transparente a macOS (Darwin arm64/amd64) e Linux (amd64/arm64) com adaptação automática a ambientes sem display (`--headless=new`, `--disable-dev-shm-usage`, `--no-sandbox`).
- **Templates Embutidos:** Uso de `//go:embed` para que o binário funcione offline e isolado (Modo Standalone) ou conectado ao `harness-core`.

## Log de Decisões
### [2026-09-28] MVP e Arquitetura Inicial
- **Contexto:** Necessidade de transformar as convenções de Session Agent do Antigravity em um produto de engenharia reproduzível e versionado.
- **Decisão:** Criado o repositório `antigravity-operator` com a CLI `agyo` contendo os subcomandos `init`, `doctor`, `browser` e `sync`.
- **Trade-offs:** Escolha de Go em vez de Python/Bash para garantir portabilidade instantânea e zero dependência de interpretadores nas máquinas de destino.
