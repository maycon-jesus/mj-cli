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
				"command.root.long_description":  "MJ CLI is a powerful command line interface tool to interact with MJ services.",
			},
			"pt-BR": {
				"command.root.short_description": "MJ CLI - Uma ferramenta de linha de comando",
				"command.root.long_description":  "MJ CLI é uma poderosa ferramenta de linha de comando para interagir com os serviços MJ.",
			},
		},
	}
}
