package commandsrepository

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewRootCommand() *commands.Command {
	return &commands.Command{
		Name:                "mj-cli",
		ShortDescriptionKey: "command.root.short_description",
		LongDescriptionKey:  "command.root.long_description",
		RunHelpOnNoArgs:     true,
		Translations: intl.Translations{
			"en": {
				"command.root.short_description": "MJ CLI - A command line tool",
			},
			"pt-BR": {
				"command.root.short_description": "MJ CLI - Uma ferramenta de linha de comando",
			},
		},
	}
}
