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
			{Name: "registry", Description: "The configuration registry to modify", Required: true},
			{Name: "key", Description: "The configuration key to set", Required: false},
		},
		ShortDescription: "Get a configuration value",
		Translations: intl.Translations{
			"en": {
				"config.get.invalid_args":  "expected at most 2 arguments, got {{count}}",
				"config.get.name_registry": "All settings in registry '{{registry}}':",
			},
			"pt-BR": {
				"config.get.invalid_args":  "esperado no máximo 2 argumentos, recebidos {{count}}",
				"config.get.name_registry": "Todas as configurações do registro '{{registry}}':",
			},
		},
		BeforeRun: func(ctx context.Context, execData *commands.ExecData) error {
			if len(execData.Args) > 2 {
				return fmt.Errorf("expected at most 2 arguments, got %d", len(execData.Args))
			}
			return nil
		},
		Handler: func(ctx context.Context, data *commands.ExecData) error {
			switch len(data.Args) {
			case 1:
				settings := data.Config.GetModule(data.Args[0]).AllSettings()
				fmt.Println(data.Translator.T("config.get.name_registry", map[string]string{"registry": data.Args[0]}))
				for k, v := range settings {
					fmt.Printf("%s: %v\n", k, v)
				}
			case 2:
				value := data.Config.GetModule(data.Args[0]).GetString(data.Args[1])
				fmt.Println(value)
			default:
				return fmt.Errorf("expected at most 2 arguments, got %d", len(data.Args))
			}

			return nil
		},
	}
}
