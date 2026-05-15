package cube_cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/cube"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newCubeScrambleCommand() *commands.Command {
	return &commands.Command{
		Name:                "scramble",
		ShortDescriptionKey: "command.cube.scramble.short_description",
		LongDescriptionKey:  "command.cube.scramble.long_description",
		RunHelpOnNoArgs:     true,
		Translations: intl.Translations{
			"en": {
				"command.cube.scramble.short_description":      "Scramble command",
				"command.cube.scramble.long_description":       "Generates a random scramble for a cube.",
				"command.cube.scramble.flags.size.description": "The number of moves in the scramble.",
				"command.cube.scramble.scramble_output":        "Generated scramble: {{scramble}}",
				"command.cube.scramble.scranble_view_output":   "Scramble view:\n{{view}}",
			},
			"pt-BR": {
				"command.cube.scramble.short_description":      "Comando embaralhar",
				"command.cube.scramble.long_description":       "Gera um embaralhamento aleatório para um cubo.",
				"command.cube.scramble.flags.size.description": "O número de movimentos no embaralhamento.",
				"command.cube.scramble.scramble_output":        "Embaralhamento gerado: {{scramble}}",
				"command.cube.scramble.scranble_view_output":   "Visualização do embaralhamento:\n{{view}}",
			},
		},
		Configs: []commands.Config{
			{
				Key:            "command.cube.scramble.size",
				DescriptionKey: "command.cube.scramble.flags.size.description",
				DefaultValue:   20,
			},
		},
		Flags: []commands.Flag{
			{
				Name:           "size",
				Shorthand:      "c",
				DescriptionKey: "command.cube.scramble.flags.size.description",
				DefaultValue:   "20",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			moveCountStr, err := execData.Flags.GetString("size")
			if err != nil {
				return err
			}

			moveCount, err := strconv.Atoi(moveCountStr)
			if err != nil {
				return err
			}
			scramble := cube.GenerateScramble333(moveCount)
			virtualCube := cube.CreateCube333()
			virtualCube.ApplyMoves(strings.Split(scramble, " "))
			execData.Translator.Println("command.cube.scramble.scramble_output", map[string]string{"scramble": scramble})
			execData.Translator.Println("command.cube.scramble.scranble_view_output", map[string]string{"view": virtualCube.Render()})
			return nil
		},
	}
}
