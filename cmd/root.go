package cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	commandsrepository "github.com/maycon-jesus/mj-cli/internal/commands_repository"
)

// Execute executa o comando raiz
func Execute(app *commands.App) {
	app.Logger.Log.Debug("Executando comando raiz")

	root := commandsrepository.NewRootCommand()
	rootCmd := root.ToCobraCommand(app)
	rootCmd.SetOut(app.Terminal)
	rootCmd.SetErr(app.Terminal)

	rootCmd.Execute()
}
