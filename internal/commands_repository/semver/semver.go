package semvercmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewSemverCommand() *commands.Command {
	return &commands.Command{
		Name:                "semver",
		ShortDescriptionKey: "command.semver.short_description",
		RunHelpOnNoArgs:     true,
		SubCommands: []*commands.Command{
			newSemverUpgradeCommand(),
		},
		Translations: intl.Translations{
			"en": {
				"command.semver.short_description": "Manage semantic versioning",
			},
			"pt-BR": {
				"command.semver.short_description": "Gerenciar versionamento semântico",
			},
		},
	}
}
