package alias_cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewAliasCommand() *commands.Command {
	return &commands.Command{
		Name:                "alias",
		ShortDescriptionKey: "command.alias.short_description",
		RunHelpOnNoArgs:     true,
		SubCommands:         []*commands.Command{newAliasSetCommand(), newAliasRunCommand(), newAliasLsCommand(), newAliasViewCommand(), newAliasRmCommand()},
		Translations: intl.Translations{
			"en": {
				"command.alias.short_description": "Manage command aliases",
			},
			"pt-BR": {
				"command.alias.short_description": "Gerenciar aliases de comandos",
			},
		},
		Configs: map[string]commands.Config{
			"command.alias.aliases": {},
		},
	}
}
