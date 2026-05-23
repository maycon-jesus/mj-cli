package config_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/ui"
)

func newSetCommand() *commands.Command {
	return &commands.Command{
		Name: "set",
		Args: []commands.Arg{
			{Name: "key", DescriptionKey: "command.config.set.arg.key", Required: true},
			{Name: "value", DescriptionKey: "command.config.set.arg.value", Required: true},
		},
		ShortDescriptionKey: "command.config.set.short_description",
		Translations: intl.Translations{
			"en": {
				"command.config.set.short_description": "Set a configuration value",
				"command.config.set.arg.key":           "The configuration key to set (dot-separated)",
				"command.config.set.arg.value":         "The value to set for the configuration key",
				"config.set.invalid_args":              "expected 2 arguments, got {{count}}",
				"config.set.not_found":                 "configuration key '{{key}}' is not registered",
				"config.set.success":                   "Configuration '{{key}}' set to '{{value}}'",
			},
			"pt-BR": {
				"command.config.set.short_description": "Definir um valor de configuração",
				"command.config.set.arg.key":           "A chave de configuração a ser definida (separada por pontos)",
				"command.config.set.arg.value":         "O valor a ser definido para a chave de configuração",
				"config.set.invalid_args":              "esperado 2 argumentos, recebidos {{count}}",
				"config.set.not_found":                 "chave de configuração '{{key}}' não está registrada",
				"config.set.success":                   "Configuração '{{key}}' definida para '{{value}}'",
			},
		},
		BeforeRun: func(ctx context.Context, execData *commands.ExecData) error {
			if len(execData.Args) != 2 {
				return execData.Translator.Errorf("config.set.invalid_args", map[string]string{
					"count": fmt.Sprintf("%d", len(execData.Args)),
				})
			}
			return nil
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			key := execData.Args[0]
			value := execData.Args[1]

			if !execData.Config.Set(key, value) {
				return execData.Translator.Errorf("config.set.not_found", map[string]string{"key": key})
			}

			fmt.Println(ui.Success(execData.Translator.T("config.set.success", map[string]string{
				"key":   key,
				"value": value,
			})))
			return nil
		},
	}
}
