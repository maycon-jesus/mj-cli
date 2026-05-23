package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/cmd"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/ui"
	"github.com/maycon-jesus/mj-cli/pkg/utils"
)

func newAliasRunCommand() *commands.Command {
	return &commands.Command{
		Name:                "run",
		ShortDescriptionKey: "command.alias.run.short_description",
		Args: []commands.Arg{
			{Name: "name", DescriptionKey: "command.alias.run.arg.name", Required: true},
		},
		Translations: intl.Translations{
			"en": {
				"command.alias.run.short_description": "Run a command alias",
				"command.alias.run.arg.name":          "The name of the alias to run",
			},
			"pt-BR": {
				"command.alias.run.short_description": "Executar um alias de comando",
				"command.alias.run.arg.name":          "O nome do alias a ser executado",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			name := execData.Args[0]

			settings, _ := execData.Config.Get("command.alias.aliases")
			aliases, _ := settings.(map[string]interface{})
			command, _ := aliases[name].(string)

			variables := make(map[string]string)
			for i, value := range execData.Args[1:] {
				variables[fmt.Sprintf("%d", i+1)] = value
			}

			command = utils.ReplaceVariables(command, variables)

			output := ui.Title("Executando alias " + name)
			execData.Terminal.Println(output)
			execData.Terminal.StartSpinner("cmd", "Executando alias...")
			output = ui.Lambda(command)
			execData.Terminal.Println(output)

			err := cmd.RunCommandWithOptions(command, cmd.CommandOptions{
				Stdout: execData.Terminal,
				Stderr: execData.Terminal,
			})
			if err != nil {
				execData.Terminal.StopSpinnerWithError("cmd", err.Error())
				return err
			} else {
				execData.Terminal.StopSpinner("cmd", "Alias executado com sucesso!")
			}

			return nil
		},
	}
}
