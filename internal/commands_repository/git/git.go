package gitcmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewGitCommand() *commands.Command {
	return &commands.Command{
		Name:                "git",
		ShortDescriptionKey: "command.git.short_description",
		RunHelpOnNoArgs:     true,
		SubCommands:         []*commands.Command{newGitNewCommand()},
		Translations: intl.Translations{
			"en": {
				"command.git.short_description": "Manage git repository",
			},
			"pt-BR": {
				"command.git.short_description": "Gerenciar repositório git",
			},
		},
	}
}
