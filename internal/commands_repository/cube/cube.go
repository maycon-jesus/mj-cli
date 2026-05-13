package cube_cmd

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewCubeCommand() *commands.Command {
	return &commands.Command{
		Name:                "cube",
		ShortDescriptionKey: "command.cube.short_description",
		LongDescriptionKey:  "command.cube.long_description",
		RunHelpOnNoArgs:     true,
		Translations: intl.Translations{
			"en": {
				"command.cube.short_description": "Cube command",
				"command.cube.long_description":  "This command does something related to cubes.",
			},
			"pt-BR": {
				"command.cube.short_description": "Comando cubo",
				"command.cube.long_description":  "Este comando faz algo relacionado a cubos.",
			},
		},
		SubCommands: []*commands.Command{
			newCubeScrambleCommand(),
		},
	}
}
