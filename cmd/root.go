package cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
)

// Execute executa o comando raiz
func Execute(app *commands.App) {
	app.Logger.Log.Debug("Executando comando raiz")

	rootCmd := app.RootCmd.ToCobraCommand(app)
	rootCmd.SetOut(app.Terminal)
	rootCmd.SetErr(app.Terminal)

	rootCmd.Execute()
}
