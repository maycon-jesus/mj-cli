package cmd

import (
	"os"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	config_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/config"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mj-cli",
	Short: "MJ CLI - Uma ferramenta de linha de comando",
	Long:  `MJ CLI é uma ferramenta de linha de comando construída com uma arquitetura modular e extensível.`,
}

// Execute executa o comando raiz
func Execute(configRegistry *config.ConfigRegistry, translator *intl.Translator) {
	// Cria o registro de comandos
	registry := commands.NewRegistry()

	// Registra os comandos
	registry.RegisterMultiple(
		config_cmd.NewConfigCommand(),
	)

	// Anexa todos os comandos ao rootCmd
	registry.AttachToRoot(rootCmd, configRegistry, translator)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
