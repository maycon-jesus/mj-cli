package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasSetCommand() *commands.Command {
	return &commands.Command{
		Name:             "set",
		ShortDescription: "Set a command alias",
		Translations: intl.Translations{
			"en": {
				"alias.set.success": "Alias '{{name}}' set to '{{command}}'",
			},
			"pt-BR": {
				"alias.set.success": "Alias '{{name}}' definido para '{{command}}'",
			},
		},
		Args: []commands.Arg{
			{Name: "name", Description: "The name of the alias to set", Required: true},
			{Name: "command", Description: "The command that the alias will represent", Required: true},
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
