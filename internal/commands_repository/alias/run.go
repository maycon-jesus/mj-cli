package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/cmd"
	"github.com/maycon-jesus/mj-cli/pkg/utils"
)

func newAliasRunCommand() *commands.Command {
	return &commands.Command{
		Name:             "run",
		ShortDescription: "Run a command alias",
		Args: []commands.Arg{
			{Name: "name", Description: "The name of the alias to run", Required: true},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := fmt.Sprintf("aliases.%s", execData.Args[0])
			command := execData.Config.GetModule("general").GetString(key)

			variables := make(map[string]string)

			for i, value := range execData.Args[1:] {
				variables[fmt.Sprintf("%d", i+1)] = value
			}

			command = utils.ReplaceVariables(command, variables)

			cmd.RunCommand(command)
			return nil
		},
	}
}
