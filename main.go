package main

import (
	_ "embed"

	"github.com/maycon-jesus/mj-cli/cmd"
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
	"github.com/maycon-jesus/mj-cli/pkg/mjterm"
)

//go:embed VERSION
var appVersion string
var appName = "mj-cli"

func main() {
	term := mjterm.New()
	defer term.Close()

	logger, err := logger.NewLoggerComplete(appName, appVersion)
	if err != nil {
		panic(err)
	}
	defer logger.FileHandler.Close()
	log := logger.Log

	log.Info("Iniciando aplicação")

	newViperAdapter := config.NewViperAdapter("mj-cli")
	newViperAdapter.SetEnvPrefix("MJ_CLI")

	configRegistry := config.NewConfigRegistry()
	configRegistry.RegisterModule("general", newViperAdapter)

	if err := newViperAdapter.ReadInConfig(); err != nil {
		log.Warn("Falha ao carregar configuração", "error", err.Error())
	} else {
		log.Debug("Configuração carregada")
	}

	lang := newViperAdapter.GetString("lang")
	translator := intl.NewTranslator(lang)
	log.Debug("Tradutor inicializado", "lang", lang)

	app := &commands.App{
		Logger:     &logger,
		Config:     configRegistry,
		Translator: translator,
		Terminal:   term,
		Version:    appVersion,
		Name:       appName,
	}

	cmd.Execute(app)

	if err := newViperAdapter.WriteConfig(); err != nil {
		log.Error("Falha ao salvar configuração", "error", err.Error())
		panic(err)
	}

	log.Debug("Configuração salva")
	log.Info("Aplicação finalizada")
}
