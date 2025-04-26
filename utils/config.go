package utils

import (
	"github.com/spf13/viper"
	"os"
	"regexp"
)

type ConfigValueValidator = func(value string) (bool, string)

type ConfigValue struct {
	Name         string
	Description  string
	DefaultValue string
	Validators   []ConfigValueValidator
}

var AllConfigs = map[string]ConfigValue{
	"obsidianVaultDir": {
		Description: "Directory default of obsidian vault",
	},
	"obsidianTagsDir": {
		Description: "Directory default of tags template on obsidian",
	},
}

func LoadViper() {
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(home)
	viper.SetConfigName(".mj-cli")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		//panic(err)
	}

	for k, v := range AllConfigs {
		if viper.Get(k) == "" {
			viper.Set(k, v.DefaultValue)
		}
	}

}

func validatorRegex(regex string) ConfigValueValidator {
	return func(value string) (bool, string) {
		match, _ := regexp.Match(regex, []byte(value))
		if match {
			return match, ""
		}
		return match, "Valor não bate com o Regexp"
	}
}
