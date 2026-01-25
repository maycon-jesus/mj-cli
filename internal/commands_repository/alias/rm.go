package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasRmCommand() *commands.Command {
	return &commands.Command{
		Name:             "rm",
		ShortDescription: "Remove an alias",
		Translations: intl.Translations{
			"en": {
				"alias.rm.not_found": "Alias '{{name}}' not found",
				"alias.rm.removed":   "Alias '{{name}}' removed successfully",
			},
			"pt-BR": {
				"alias.rm.not_found": "Alias '{{name}}' não encontrado",
				"alias.rm.removed":   "Alias '{{name}}' removido com sucesso",
			},
		},
		Args: []commands.Arg{
			{Name: "name", Description: "The name of the alias to remove", Required: true},
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
