package config_cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/utils"
	"go.yaml.in/yaml/v3"
)

func newConfigGenerateCommand() *commands.Command {
	return &commands.Command{
		Name:                "generate",
		ShortDescriptionKey: "command.config.generate.short_description",
		LongDescriptionKey:  "command.config.generate.long_description",
		Translations: intl.Translations{
			"en": {
				"command.config.generate.short_description": "Generate a configuration file",
				"command.config.generate.long_description":  "Generate a configuration file with default values and print it to stdout",
			},
			"pt-BR": {
				"command.config.generate.short_description": "Gerar um arquivo de configuração",
				"command.config.generate.long_description":  "Gera um arquivo de configuração com os valores padrão e imprime no stdout",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			out, err := yaml.Marshal(makeObject(execData.RootCmd))
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
}

func makeObject(cmd *commands.Command) map[string]interface{} {
	obj := make(map[string]interface{})

	if cmd.Configs != nil {
		for key, config := range cmd.Configs {
			configMap := utils.CreateDeepMap(strings.Split(key, "."), config.DefaultValue)
			utils.MergeDeepMaps(obj, configMap)
		}
	}

	if cmd.SubCommands != nil {
		for _, subCmd := range cmd.SubCommands {
			subObj := makeObject(subCmd)
			utils.MergeDeepMaps(obj, subObj)
		}
	}

	return obj
}
