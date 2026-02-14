package commands

import (
	"github.com/spf13/cobra"
)

// Registry mantém todos os comandos registrados
type Registry struct {
	commands []*Command
}

// NewRegistry cria um novo registro de comandos
func NewRegistry() *Registry {
	return &Registry{
		commands: make([]*Command, 0),
	}
}

// Register adiciona um comando ao registro
func (r *Registry) Register(cmd *Command) {
	r.commands = append(r.commands, cmd)
}

// RegisterMultiple adiciona múltiplos comandos ao registro
func (r *Registry) RegisterMultiple(cmds ...*Command) {
	r.commands = append(r.commands, cmds...)
}

// GetCommands retorna todos os comandos registrados
func (r *Registry) GetCommands() []*Command {
	return r.commands
}

// AttachToRoot adiciona todos os comandos registrados a um comando raiz do Cobra
func (r *Registry) AttachToRoot(rootCmd *cobra.Command, app *App) {
	app.Logger.Log.Debug("Attaching commands to root command", "command_count", len(r.commands))
	for _, cmd := range r.commands {
		rootCmd.AddCommand(cmd.ToCobraCommand(app))
	}
}
