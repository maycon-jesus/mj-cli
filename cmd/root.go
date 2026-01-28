package cmd

import (
	"os"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	commandsrepository "github.com/maycon-jesus/mj-cli/internal/commands_repository"
	alias_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/alias"
	config_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/config"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

// Execute executa o comando raiz
func Execute(configRegistry *config.ConfigRegistry, translator *intl.Translator) {
	root := commandsrepository.NewRootCommand()
	rootCmd := root.ToCobraCommand(configRegistry, translator)

	// Cria o registro de comandos
	registry := commands.NewRegistry()

	// Registra os comandos
	registry.RegisterMultiple(
		config_cmd.NewConfigCommand(),
		alias_cmd.NewAliasCommand(),
	)

	// Anexa todos os comandos ao rootCmd
	registry.AttachToRoot(rootCmd, configRegistry, translator)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
