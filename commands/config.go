package commands

import (
	"errors"
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Obsidian = &cobra.Command{
	Use:   "config <config-name> [config-value]",
	Short: "Obter ou definir configuração",
	Args:  cobra.RangeArgs(0, 2),
	RunE:  runE,
}

func GetCommandConfig() *cobra.Command {
	return Obsidian
}

func runE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("Possible configs:")
		for k, v := range utils.AllConfigs {
			fmt.Println(fmt.Sprintf("- %s: %s", k, v.Description))
		}
		return nil
	}

	configName := args[0]
	config, ok := utils.AllConfigs[configName]
	if !ok {
		return errors.New("Config not found")
	}

	if len(args) == 1 {
		fmt.Println(viper.Get(configName))
		return nil
	}

	configSetValue := args[1]

	for _, validator := range config.Validators {
		ok, msg := validator(configSetValue)

		if !ok {
			return errors.New(msg)

		}
	}

	viper.Set(configName, configSetValue)

	return nil
}
