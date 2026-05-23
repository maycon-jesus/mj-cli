package config_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newGetCommand() *commands.Command {
	return &commands.Command{
		Name: "get",
		Args: []commands.Arg{
			{Name: "key", DescriptionKey: "command.config.get.arg.key", Required: true},
		},
		ShortDescriptionKey: "command.config.get.short_description",
		Translations: intl.Translations{
			"en": {
				"command.config.get.short_description": "Get a configuration value",
				"command.config.get.arg.key":           "The configuration key to get (dot-separated)",
				"config.get.invalid_args":              "expected 1 argument, got {{count}}",
				"config.get.not_found":                 "configuration key '{{key}}' is not registered",
			},
			"pt-BR": {
				"command.config.get.short_description": "Obter um valor de configuração",
				"command.config.get.arg.key":           "A chave de configuração a ser obtida (separada por pontos)",
				"config.get.invalid_args":              "esperado 1 argumento, recebidos {{count}}",
				"config.get.not_found":                 "chave de configuração '{{key}}' não está registrada",
			},
		},
		BeforeRun: func(ctx context.Context, execData *commands.ExecData) error {
			if len(execData.Args) != 1 {
				return execData.Translator.Errorf("config.get.invalid_args", map[string]string{
					"count": fmt.Sprintf("%d", len(execData.Args)),
				})
			}
			return nil
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := execData.Args[0]
			value, exists := execData.Config.Get(key)
			if !exists {
				return execData.Translator.Errorf("config.get.not_found", map[string]string{"key": key})
			}
			fmt.Printf("%v\n", value)
			return nil
		},
	}
}
