package config_cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewConfigCommand() *commands.Command {
	setCommand := newSetCommand()
	getCommand := newGetCommand()
	generateCommand := newConfigGenerateCommand()

	return &commands.Command{
		Name:                "config",
		ShortDescriptionKey: "command.config.short_description",
		SubCommands:         []*commands.Command{setCommand, getCommand, generateCommand},
		RunHelpOnNoArgs:     true,
		Translations: intl.Translations{
			"en": {
				"command.config.short_description": "Manage settings",
			},
			"pt-BR": {
				"command.config.short_description": "Gerenciar configurações",
			},
		},
	}
}
