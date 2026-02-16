package log_cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
)

func NewLogCommand() *commands.Command {
	return &commands.Command{
		Name:                "log",
		ShortDescriptionKey: "command.log.short_description",
		RunHelpOnNoArgs:     true,
		Translations: map[string]map[string]string{
			"en": {
				"command.log.short_description": "Manages the logs of this application in a structured and stylized way",
			},
			"pt-BR": {
				"command.log.short_description": "Gerencia os logs desta aplicação de forma estruturada e estilizada",
			},
		},
		SubCommands: []*commands.Command{
			NewLogViewCommand(),
		},
	}
}
