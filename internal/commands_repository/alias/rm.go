package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasRmCommand() *commands.Command {
	return &commands.Command{
		Name:                "rm",
		ShortDescriptionKey: "command.alias.rm.short_description",
		Translations: intl.Translations{
			"en": {
				"command.alias.rm.short_description": "Remove an alias",
				"command.alias.rm.arg.name":          "The name of the alias to remove",
				"alias.rm.not_found":                 "Alias '{{name}}' not found",
				"alias.rm.removed":                   "Alias '{{name}}' removed successfully",
			},
			"pt-BR": {
				"command.alias.rm.short_description": "Remover um alias",
				"command.alias.rm.arg.name":          "O nome do alias a ser removido",
				"alias.rm.not_found":                 "Alias '{{name}}' não encontrado",
				"alias.rm.removed":                   "Alias '{{name}}' removido com sucesso",
			},
		},
		Args: []commands.Arg{
			{Name: "name", DescriptionKey: "command.alias.rm.arg.name", Required: true},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := fmt.Sprintf("aliases.%s", execData.Args[0])
			command := execData.Config.GetModule("general").GetString(key)

			if command == "" {
				execData.Translator.Println("alias.rm.not_found", map[string]string{"name": execData.Args[0]})
				return nil
			}

			execData.Config.GetModule("general").Set(key, nil)
			execData.Config.GetModule("general").WriteConfig()
			execData.Translator.Println("alias.rm.removed", map[string]string{"name": execData.Args[0]})
			return nil
		},
	}
}
