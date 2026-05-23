package alias_cmd

import (
	"context"

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
			name := execData.Args[0]

			settings, _ := execData.Config.Get("command.alias.aliases")
			aliases, ok := settings.(map[string]interface{})
			if !ok || aliases == nil {
				execData.Translator.Println("alias.rm.not_found", map[string]string{"name": name})
				return nil
			}

			if _, exists := aliases[name]; !exists {
				execData.Translator.Println("alias.rm.not_found", map[string]string{"name": name})
				return nil
			}

			delete(aliases, name)
			execData.Config.Set("command.alias.aliases", aliases)
			execData.Translator.Println("alias.rm.removed", map[string]string{"name": name})
			return nil
		},
	}
}
