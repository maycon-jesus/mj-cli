package alias_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newAliasLsCommand() *commands.Command {
	return &commands.Command{
		Name:                "ls",
		ShortDescriptionKey: "command.alias.ls.short_description",
		Translations: intl.Translations{
			"en": {
				"command.alias.ls.short_description": "List all aliases",
				"alias.ls.not_found":                 "No aliases found",
				"alias.ls.header":                    "Aliases:",
			},
			"pt-BR": {
				"command.alias.ls.short_description": "Listar todos os aliases",
				"alias.ls.not_found":                 "Nenhum alias encontrado",
				"alias.ls.header":                    "Aliases:",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			settings, _ := execData.Config.Get("command.alias.aliases")
			if settings == nil {
				execData.Translator.Println("alias.ls.not_found", map[string]string{})
				return nil
			}

			aliases, ok := settings.(map[string]interface{})
			if !ok || len(aliases) == 0 {
				execData.Translator.Println("alias.ls.not_found", map[string]string{})
				return nil
			}

			execData.Translator.Println("alias.ls.header", map[string]string{})
			for name, command := range aliases {
				fmt.Printf("  %s = %s\n", name, command)
			}
			return nil
		},
	}
}
