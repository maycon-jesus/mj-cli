package config_cmd

import "github.com/maycon-jesus/mj-cli/internal/commands"

func NewConfigCommand() *commands.Command {
	setCommand := newSetCommand()
	getCommand := newGetCommand()

	return &commands.Command{
		Name:             "config",
		ShortDescription: "Manage settings",
		SubCommands:      []*commands.Command{setCommand, getCommand},
		RunHelpOnNoArgs:  true,
	}
}
