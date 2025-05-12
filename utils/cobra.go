package utils

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateFlagWithViperConfig(cmd *cobra.Command, name string, value string, usage string, configKey string) {
	configValue := viper.GetString(configKey)
	description := ""

	if configValue != "" {
		description = fmt.Sprintf("%s\nConfig Key: %s\nConfig Value: %s", usage, configKey, configValue)
	} else {
		description = fmt.Sprintf("%s\nConfig Key: %s", usage, configKey)
	}

	cmd.Flags().String(name, value, description)

	err := cmd.Flags().Lookup(name).Value.Set(configValue)
	cobra.CheckErr(err)
}

func CreateRequiredFlagWithViperConfig(cmd *cobra.Command, name string, value string, usage string, configKey string) {
	CreateFlagWithViperConfig(cmd, name, value, usage, configKey)
	configValue := viper.GetString(configKey)
	if configValue == "" {
		err := cmd.MarkFlagRequired(name)
		cobra.CheckErr(err)
	}
}
