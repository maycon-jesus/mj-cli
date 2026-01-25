package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasViewCommand() *commands.Command {
	return &commands.Command{
		Name:             "view",
		ShortDescription: "View a specific alias",
		Translations: intl.Translations{
			"en": {
				"alias.view.not_found": "Alias '{{name}}' not found",
			},
			"pt-BR": {
				"alias.view.not_found": "Alias '{{name}}' não encontrado",
			},
		},
		Args: []commands.Arg{
			{Name: "name", Description: "The name of the alias to view", Required: true},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := fmt.Sprintf("aliases.%s", execData.Args[0])
			command := execData.Config.GetModule("general").GetString(key)

			if command == "" {
				execData.Translator.Println("alias.view.not_found", map[string]string{"name": execData.Args[0]})
				return nil
			}

			fmt.Printf("%s = %s\n", execData.Args[0], command)
			return nil
		},
	}
}
