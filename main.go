package main

import (
	"github.com/maycon-jesus/mj-cli/cmd"
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
)

func main() {
	log, err := logger.NewWithTemporaryFile("mj-cli")
	if err != nil {
		panic(err)
	}
	defer log.Close()
	defer log.RecoverPanic()

	newViperAdapter := config.NewViperAdapter("mj-cli")
	newViperAdapter.SetEnvPrefix("MJ_CLI")

	configRegistry := config.NewConfigRegistry()
	configRegistry.RegisterModule("general", newViperAdapter)
	newViperAdapter.ReadInConfig()

	translator := intl.NewTranslator(newViperAdapter.GetString("lang"))

	app := &commands.App{
		Logger:     log,
		Config:     configRegistry,
		Translator: translator,
	}

	cmd.Execute(app)

	err = newViperAdapter.WriteConfig()
	if err != nil {
		panic(err)
	}
}
