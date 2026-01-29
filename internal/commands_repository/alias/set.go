package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasSetCommand() *commands.Command {
	return &commands.Command{
		Name:                "set",
		ShortDescriptionKey: "command.alias.set.short_description",
		Args: []commands.Arg{
			{Name: "name", DescriptionKey: "command.alias.set.arg.name", Required: true},
			{Name: "command", DescriptionKey: "command.alias.set.arg.command", Required: true},
		},
		Translations: intl.Translations{
			"en": {
				"command.alias.set.short_description": "Set a command alias",
				"command.alias.set.arg.name":          "The name of the alias to set",
				"command.alias.set.arg.command":       "The command that the alias will represent",
				"alias.set.success":                   "Alias '{{name}}' set to '{{command}}'",
			},
			"pt-BR": {
				"command.alias.set.short_description": "Definir um alias de comando",
				"command.alias.set.arg.name":          "O nome do alias a ser definido",
				"command.alias.set.arg.command":       "O comando que o alias representará",
				"alias.set.success":                   "Alias '{{name}}' definido para '{{command}}'",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := fmt.Sprintf("aliases.%s", execData.Args[0])
			execData.Config.GetModule("general").Set(key, execData.Args[1])
			execData.Config.GetModule("general").WriteConfig()
			execData.Translator.Println("alias.set.success", map[string]string{"name": execData.Args[0], "command": execData.Args[1]})
			return nil
		},
	}
}
