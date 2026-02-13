package cmd

import (
	"os"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	commandsrepository "github.com/maycon-jesus/mj-cli/internal/commands_repository"
	alias_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/alias"
	config_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/config"
	branch_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/git"
)

// Execute executa o comando raiz
func Execute(app *commands.App) {
	app.Logger.Log.Debug("Executando comando raiz")

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
