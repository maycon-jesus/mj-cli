package config_cmd

import (
	"context"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/tint"
)

func newSetCommand() *commands.Command {
	return &commands.Command{
		Name: "set",
		Args: []commands.Arg{
			{Name: "registry", DescriptionKey: "command.config.set.arg.registry", Required: true},
			{Name: "key", DescriptionKey: "command.config.set.arg.key", Required: true},
			{Name: "value", DescriptionKey: "command.config.set.arg.value", Required: true},
		},
		ShortDescriptionKey: "command.config.set.short_description",
		Translations: intl.Translations{
			"en": {
				"command.config.set.short_description": "Set a configuration value",
				"command.config.set.arg.registry":      "The configuration registry to modify",
				"command.config.set.arg.key":           "The configuration key to set",
				"command.config.set.arg.value":         "The value to set for the configuration key",
				"config.set.expected_args":             "expected 3 arguments, got {{count}}",
				"config.set.failed_write":              "failed to write configuration: {{error}}",
				"config.set.success":                   "Configuration '{{key}}' set to '{{value}}' in registry '{{registry}}'",
			},
			"pt-BR": {
				"command.config.set.short_description": "Definir um valor de configuração",
				"command.config.set.arg.registry":      "O registro de configuração a ser modificado",
				"command.config.set.arg.key":           "A chave de configuração a ser definida",
				"command.config.set.arg.value":         "O valor a ser definido para a chave de configuração",
				"config.set.expected_args":             "esperado 3 argumentos, recebido {{count}}",
				"config.set.failed_write":              "falha ao gravar a configuração: {{error}}",
				"config.set.success":                   "Configuração '{{key}}' definida para '{{value}}' no registro '{{registry}}'",
			},
		},
		BeforeRun: func(ctx context.Context, execData *commands.ExecData) error {
			if len(execData.Args) != 3 {
				return execData.Translator.Errorf("config.set.expected_args", map[string]string{"count": fmt.Sprintf("%d", len(execData.Args))})
			}
			return nil
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			registry := execData.Args[0]
			key := execData.Args[1]
			value := execData.Args[2]

			configModule := execData.Config.GetModule(registry)
			configModule.Set(key, value)
			err := configModule.WriteConfig()
			if err != nil {
				return execData.Translator.Errorf("config.set.failed_write", map[string]string{"error": err.Error()})
			}
			str := execData.Translator.T("config.set.success", map[string]string{
				"key":      key,
				"value":    value,
				"registry": registry,
			})
			fmt.Println(tint.PresetSuccess(str))

			return nil
		},
	}
}
