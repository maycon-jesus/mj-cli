package cmd

import (
	"os"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	commandsrepository "github.com/maycon-jesus/mj-cli/internal/commands_repository"
	alias_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/alias"
	config_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/config"
	branch_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/git"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
)

// Execute executa o comando raiz
func Execute(app *commands.App) {
	app.Logger.Debug("Executando comando raiz", logger.Metadata{"event": "execute_root_command"})

	root := commandsrepository.NewRootCommand()
	rootCmd := root.ToCobraCommand(app)

	// Cria o registro de comandos
	registry := commands.NewRegistry()

	// Registra os comandos
	registry.RegisterMultiple(
		config_cmd.NewConfigCommand(),
		alias_cmd.NewAliasCommand(),
		branch_cmd.NewGitCommand(),
	)

	// Anexa todos os comandos ao rootCmd
	registry.AttachToRoot(rootCmd, app)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
