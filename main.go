package main

import (
	"github.com/maycon-jesus/mj-cli/cmd"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func main() {
	newViperAdapter := config.NewViperAdapter("mj-cli")
	newViperAdapter.SetEnvPrefix("MJ_CLI")

	configRegistry := config.NewConfigRegistry()
	configRegistry.RegisterModule("general", newViperAdapter)
	newViperAdapter.ReadInConfig()

	translator := intl.NewTranslator(newViperAdapter.GetString("lang"))
	cmd.Execute(configRegistry, translator)

	err := newViperAdapter.WriteConfig()
	if err != nil {
		panic(err)
	}
}
