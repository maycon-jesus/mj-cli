package utils

import (
	"github.com/spf13/viper"
	"os"
	"regexp"
)

type ConfigValueValidator = func(value string) (bool, string)

type ConfigValue struct {
	Description  string
	DefaultValue string
	Validators   []ConfigValueValidator
}

var AllConfigs = map[string]ConfigValue{
	"obsidianVaultDir": {
		Description: "Directory default of obsidian vault",
	},
	"obsidianBulletJournalDir": {
		Description:  "Directory default of bullet journal on obsidian vault",
		DefaultValue: "03 - Journal",
	},
	"obsidianBulletJournalMonthDir": {
		Description:  "",
		DefaultValue: "01 - Meses",
	},
	"obsidianBulletJournalWeeklyDir": {
		Description:  "",
		DefaultValue: "02 - Semanas",
	},
	"obsidianBulletJournalDailyDir": {
		Description:  "",
		DefaultValue: "03 - Diários",
	},
	"obsidianTemplatesDir": {
		Description:  "",
		DefaultValue: "99 - Meta/00 - Templates",
	},
	"obsidian-vault-dir": {
		Description: "",
	},
	"obsidian-daily-template-path": {
		Description: "",
	},
	"obsidian-daily-note-dir": {
		Description: "",
	},
	"obsidian-weekly-template-path": {
		Description: "",
	},
	"obsidian-weekly-note-dir": {
		Description: "",
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
		if val := viper.Get(k); val == "" || val == nil {
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
