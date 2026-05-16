package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasViewCommand() *commands.Command {
	return &commands.Command{
		Name:                "view",
		ShortDescriptionKey: "command.alias.view.short_description",
		Args: []commands.Arg{
			{Name: "name", DescriptionKey: "command.alias.view.arg.name", Required: true},
		},
		Translations: intl.Translations{
			"en": {
				"command.alias.view.short_description": "View a specific alias",
				"command.alias.view.arg.name":          "The name of the alias to view",
				"alias.view.not_found":                 "Alias '{{name}}' not found",
			},
			"pt-BR": {
				"command.alias.view.short_description": "Visualizar um alias específico",
				"command.alias.view.arg.name":          "O nome do alias a ser visualizado",
				"alias.view.not_found":                 "Alias '{{name}}' não encontrado",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := fmt.Sprintf("command.alias.aliases.%s", execData.Args[0])
			command := execData.ConfigGeneral.GetString(key)

			if command == "" {
				execData.Translator.Println("alias.view.not_found", map[string]string{"name": execData.Args[0]})
				return nil
			}

			fmt.Printf("%s = %s\n", execData.Args[0], command)
			return nil
		},
	}
}
