package cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	commandsrepository "github.com/maycon-jesus/mj-cli/internal/commands_repository"
	alias_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/alias"
	config_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/config"
	branch_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/git"
	log_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/log"
	semvercmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/semver"
)

// Execute executa o comando raiz
func Execute(app *commands.App) {
	app.Logger.Log.Debug("Executando comando raiz")

	root := commandsrepository.NewRootCommand()
	rootCmd := root.ToCobraCommand(app)
	rootCmd.SetOut(app.Terminal)
	rootCmd.SetErr(app.Terminal)

	// Cria o registro de comandos
	registry := commands.NewRegistry()

	// Registra os comandos
	registry.RegisterMultiple(
		config_cmd.NewConfigCommand(),
		alias_cmd.NewAliasCommand(),
		branch_cmd.NewGitCommand(),
		log_cmd.NewLogCommand(),
		semvercmd.NewSemverCommand(),
	)

	// Anexa todos os comandos ao rootCmd
	registry.AttachToRoot(rootCmd, app)

	rootCmd.Execute()
}
