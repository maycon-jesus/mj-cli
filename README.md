# mj-cli

Uma ferramenta de linha de comando personalizavel para automatizar tarefas do dia a dia no terminal.

## Funcionalidades

- **Aliases** - Crie e gerencie atalhos para comandos frequentes, com suporte a argumentos posicionais (`{{1}}`, `{{2}}`, etc.)
- **Configuracao** - Sistema de configuracao baseado em TOML com suporte a variaveis de ambiente
- **Git** - Automacao de workflows Git (criacao de branches com checkout + pull automatico)
- **Internacionalizacao** - Interface em Portugues (pt-BR) e Ingles (en)
- **Terminal Rico** - Spinners, cores ANSI com deteccao automatica e formatacao semantica

## Instalacao

### Build a partir do codigo fonte

Requisitos: [Go](https://go.dev/) 1.25.4+

```bash
git clone https://github.com/maycon-jesus/mj-cli.git
cd mj-cli
make build
```

O binario sera gerado em `build/mj-cli`.

### Cross-compilation

```bash
make build-all  # Linux e Windows (amd64 + arm64)
```

## Uso

### Aliases

Crie atalhos para comandos que voce usa com frequencia:

```bash
# Definir um alias
mj-cli alias set gc 'git commit -m "{{1}}"'

# Executar o alias
mj-cli alias run gc "feat: minha alteracao"

# Listar aliases
mj-cli alias ls

# Ver um alias especifico
mj-cli alias view gc

# Remover alias
mj-cli alias rm gc
```

Os argumentos `{{1}}`, `{{2}}`, etc. sao substituidos pelos argumentos passados ao executar o alias.

### Configuracao

```bash
# Ver todas as configuracoes de um modulo
mj-cli config get general

# Ver uma configuracao especifica
mj-cli config get general lang

# Alterar uma configuracao
mj-cli config set general lang en
```

**Arquivos de configuracao:**
- Projeto: `./general.toml`
- Usuario: `~/.mj-cli/config.toml`

As configuracoes tambem podem ser definidas via variaveis de ambiente com o prefixo `MJ_CLI_` (ex: `MJ_CLI_LANG=en`).

### Git

```bash
# Criar uma nova branch (faz checkout para main, pull e cria a branch)
mj-cli git new feat/minha-feature
```

## Arquitetura

```
mj-cli/
├── cmd/                         # Setup do CLI (Cobra)
├── internal/
│   ├── commands/                # Framework de comandos
│   ├── commands_repository/     # Implementacao dos comandos
│   │   ├── alias/               # Comandos de alias
│   │   ├── config/              # Comandos de configuracao
│   │   └── git/                 # Comandos de git
│   └── config/                  # Sistema de configuracao (Viper)
├── pkg/
│   ├── cmd/                     # Execucao de comandos shell
│   ├── intl/                    # Internacionalizacao (i18n)
│   ├── mjterm/                  # Terminal UI (spinners, input)
│   ├── tint/                    # Estilizacao ANSI
│   ├── ui/                      # Templates de UI (Title, Success, Error, etc.)
│   └── utils/                   # Utilitarios
├── main.go                      # Ponto de entrada
└── Makefile                     # Build, test, lint
```

## Desenvolvimento

```bash
# Instalar dependencias
make deps

# Executar diretamente
make run

# Rodar testes
make test

# Cobertura de testes
make test-coverage

# Lint
make lint

# Formatar codigo
go fmt ./...

# Limpar artefatos
make clean
```

## Tecnologias

- [Go](https://go.dev/) 1.25.4
- [Cobra](https://github.com/spf13/cobra) - Framework CLI
- [Viper](https://github.com/spf13/viper) - Gerenciamento de configuracao
- [TOML](https://toml.io/) - Formato de configuracao

## Licenca

Este projeto esta licenciado sob a [MIT License](LICENSE).
