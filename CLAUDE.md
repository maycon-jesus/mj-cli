# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
# Build
make build                    # Build for current platform → build/mj-cli
make build-all               # Cross-compile for linux/windows (amd64/arm64)
VERSION=1.0.0 make build     # Build with custom version

# Test
make test                    # Run all tests
make test-coverage           # Generate coverage report (opens coverage.html)
go test ./pkg/cmd -v         # Test specific package
go test -run TestParseCommand ./pkg/cmd  # Run single test

# Lint & Format
make lint                    # Run go vet
go fmt ./...                 # Format code

# Dependencies
make deps                    # Download and tidy modules

# Other
make run                     # Run the application directly
make clean                   # Remove build artifacts and coverage files
```

## Architecture Overview

### Bootstrap Flow
`main.go` creates shared services (Terminal, Logger, ConfigRegistry, Translator, DatabaseService), builds the command tree via `commandsrepository.NewRootCommand()`, packs everything into an `App` struct (along with `Version` from the embedded `VERSION` file and `Name`), and passes it to `cmd.Execute()`. `cmd.Execute()` calls `app.RootCmd.ToCobraCommand(app)`, which recursively converts the custom `Command` tree into a Cobra command tree. Config is saved on exit.

The app data directory is `~/.mj-cli/` and is also where the SQLite database (`app.db`) lives.

### Command Framework
The CLI uses a custom command framework built on top of Cobra. The framework is defined in `internal/commands/` and commands are implemented in `internal/commands_repository/`.

**Adding a new command:**
1. Create a new package under `internal/commands_repository/` with a `New*Command()` constructor returning a `*Command`
2. Define the `Command` struct with handlers, flags, `Configs`, and i18n keys; nest children via the `SubCommands` field
3. Attach it to a parent's `SubCommands` slice — a group command in `commands_repository/`, or `NewRootCommand()` in `internal/commands_repository/root.go` for a top-level command

Commands use `CommandHandler` functions with signature:
```go
func(ctx context.Context, execData *ExecData) error
```

`ExecData` provides access to args, flags, config registry, translator, terminal, logger, and `RootCmd` (the full custom command tree). A single `Terminal` and `Logger` are created in `main.go` and shared across all command handlers via the `App` struct.

The `Command` struct also supports:
- `Configs []Config` — declarative list of config keys (with default values) the command reads. Use this instead of touching the registry ad-hoc; see `internal/commands_repository/cube/scramble.go` for an example. The `config generate` command walks the `ExecData.RootCmd` tree collecting these to emit a default config file.
- `BeforeRun` / `AfterRun` hooks — run before/after `Handler`. Returning a non-nil error aborts execution and is logged to the file handler. Hooks defined on a parent command are explicitly invoked by every subcommand's `PreRunE`/`PostRunE` (which walk `cmd.Parent()`), so a top-level hook fires for the whole subtree.
- `RunHelpOnNoArgs` — when true and no `Handler` is set, running the command prints its help (used by group commands like the root).

### Alias Argument Substitution
Aliases stored in config use `{{1}}`, `{{2}}`, etc. for positional arguments. When `alias run` is invoked, extra args after the alias name are substituted via `pkg/utils.ReplaceVariables()`.

### Terminal UI (pkg/mjterm)
Event-driven terminal manager for rich output:
- Thread-safe print/prompt handling via a single event channel processed by a goroutine
- Spinner support with automatic line clearing during input
- Implements `io.Writer` for use as stdout/stderr in command execution
- Use `mjterm.New()` to create, `Close()` when done

### Logging (pkg/logger)
JSON-structured logger built on `log/slog` with custom handlers:
- `FileHandler` — writes JSON lines to a temp file (`mj-cli-log-*` in os temp dir); groups are nested JSON objects; level is `Info` by default
- `MultiHandler` — fans out log records to multiple named handlers (currently only "file")
- `DEBUG` env var overrides the file handler level (e.g., `DEBUG=debug`, `DEBUG=trace`); if the value isn't a valid slog level, defaults to `Debug`
- Supports `WithAttrs()` and `WithGroup()` for structured context
- Thread-safe: cloned handlers share a `*sync.RWMutex` for file writes

### Configuration System
- Uses Viper adapter (`internal/config/`) with registry pattern
- Modules registered by name (e.g. `"general"`); accessed via `execData.Config.GetModule("general")`
- Config file is `config.yaml`, searched in the project dir (`.`) then `~/.mj-cli/` (user); format is YAML. `config generate` emits YAML via `yaml.Marshal`
- Environment override prefix: `MJ_CLI_` (dots replaced with underscores: `MJ_CLI_LANG`)
- Flags can be bound to config keys via `FlagConfigRegistry` on the `Flag` struct
- Config persists automatically on exit (WriteConfig in `main.go`)

### Internationalization (i18n)
- Thread-safe translator at `pkg/intl/intl.go`
- Supports `en` and `pt-BR`
- Variable substitution: `"Hello {{name}}"`
- Each command defines its translations in `Translations` field

### Services Layer (internal/services)
Shared business logic used by command handlers. Services wrap `pkg/cmd` calls and provide higher-level operations (e.g., `GitService` with `DetectMainBranch()`, `Checkout()`, `NewBranch()`, `Pull()`, `BranchExists()`). Services accept a `*slog.Logger` via `WithLogger()` for structured logging.

### Persistence (GORM + SQLite)
- `DatabaseService` (`internal/services/database.go`) opens a SQLite database via `gorm.io/gorm` + `github.com/glebarez/sqlite` (pure-Go driver, no CGO) at `<appDataDir>/<dbname>` and exposes `*gorm.DB` as `DB`.
- A single `*DatabaseService` is constructed in `main.go` and shared on the `App` struct. `Startup()` is lazy via `sync.Once` and must be called before `DB` is used.
- No feature service currently persists data. When one is added, it should take `*gorm.DB` in its constructor and call `db.AutoMigrate(...)` on its own models — there is no central migration step. Define models in the same file as the service that owns them.

### Command Execution (pkg/cmd)
Three modes for running shell commands:
- `RunCommand()` - standard streams
- `RunCommandWithOptions()` - custom I/O redirection (e.g., pipe to Terminal)
- `GetCommandOutput()` - capture output as string

Parser handles quotes, escapes, and bash-compatible parsing.

### Terminal Styling (pkg/tint) & UI Templates (pkg/ui)
`pkg/tint` — low-level ANSI color/styling library:
- Automatic terminal capability detection (4-bit, 8-bit, 24-bit colors)
- Style composition: `tint.NewStyle().Bold().Foreground(tint.Red)`
- NO_COLOR environment variable support

`pkg/ui` — high-level wrappers on top of tint for consistent UI:
- `Title()`, `Subtitle()`, `Success()`, `Error()`, `Warning()`, `Info()`, `Lambda()`

### Rubik's Cube Domain (pkg/cube)
Pure domain library for a 3x3 Rubik's cube, independent of CLI plumbing.
- `Cube333` represents state as 6 `*CubeFace` arrays; faces are addressed via `FaceU/D/F/B/L/R` constants.
- Moves are declared as `MovesCollection` (map of move name → `MoveSet`) where each `PieceMove` describes a sticker rotation between two faces. `Cube333MovesCollection` covers the standard `U/D/F/B/L/R` (and their `'`/`2` variants) used by the scramble generator and the `cube scramble` command.
- When extending moves or adding cube variants, keep state mutation inside this package and have command handlers only call high-level operations (`ApplyMoves`, scramble generators, etc.) — don't reach into faces from the command layer.

## Project Conventions

- **Go version**: 1.25.4
- **Commit messages**: Conventional commits in Portuguese-BR (pt-BR)
- **Language**: Default UI language is pt-BR, configurable via `lang` config key
- **Thread safety**: Translator, Terminal, Logger, and Style are all thread-safe (mutex-protected)
- **Test style**: Table-driven tests with `[]struct` pattern; benchmarks in `pkg/cmd`
- **Version**: Embedded from `VERSION` file via `//go:embed`

## Key Entry Points

| Purpose | Location |
|---------|----------|
| Main entry | `main.go` |
| CLI setup (`Execute`) | `cmd/root.go` |
| Command tree root | `internal/commands_repository/root.go` |
| Command framework | `internal/commands/` |
| Command implementations | `internal/commands_repository/{alias,config,cube,git,log,semver}/` |
| Services | `internal/services/` |
| Config management | `internal/config/` |
| Command execution | `pkg/cmd/` |
| i18n | `pkg/intl/` |
| Terminal UI | `pkg/mjterm/` |
| Logging | `pkg/logger/` |
| Terminal styling | `pkg/tint/` |
| UI templates | `pkg/ui/` |
| Rubik's cube domain | `pkg/cube/` |
