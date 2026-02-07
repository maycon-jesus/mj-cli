package main

import (
	"os"

	"github.com/maycon-jesus/mj-cli/cmd"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
)

func main() {
	logFile, _ := os.CreateTemp("", "mj-cli-log-*.json")
	log := logger.New(logFile)

	newViperAdapter := config.NewViperAdapter("mj-cli")
	newViperAdapter.SetEnvPrefix("MJ_CLI")

	configRegistry := config.NewConfigRegistry()
	configRegistry.RegisterModule("general", newViperAdapter)
	newViperAdapter.ReadInConfig()

	translator := intl.NewTranslator(newViperAdapter.GetString("lang"))

	app := &cmd.App{
		Logger:     log,
		Config:     configRegistry,
		Translator: translator,
	}

	cmd.Execute(app)

	err := newViperAdapter.WriteConfig()
	if err != nil {
		panic(err)
	}
}
