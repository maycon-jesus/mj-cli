package cube_cmd

import (
	"context"
	"strconv"

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
				"command.cube.scramble.short_description":           "Scramble command",
				"command.cube.scramble.long_description":            "Generates a random scramble for a cube.",
				"command.cube.scramble.args.move_count.description": "The number of moves in the scramble.",
			},
			"pt-BR": {
				"command.cube.scramble.short_description":           "Comando embaralhar",
				"command.cube.scramble.long_description":            "Gera um embaralhamento aleatório para um cubo.",
				"command.cube.scramble.args.move_count.description": "O número de movimentos no embaralhamento.",
			},
		},
		Args: []commands.Arg{
			{
				Name:           "move_count",
				DescriptionKey: "command.cube.scramble.args.move_count.description",
				Required:       true,
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			moveCount, err := strconv.Atoi(execData.Args[0])
			if err != nil {
				return err
			}
			scramble := cube.GenerateScramble333(moveCount)
			execData.Terminal.Println(scramble)
			return nil
		},
	}
}
