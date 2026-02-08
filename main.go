package main

import (
	_ "embed"
	"os"

	"github.com/maycon-jesus/mj-cli/cmd"
	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
)

//go:embed VERSION
var appVersion string
var appName = "mj-cli"

func main() {
	log, err := logger.NewWithTemporaryFile("mj-cli")
	log = log.WithAttrs(logger.Metadata{
		"app":     appName,
		"version": appVersion,
	})
	if err != nil {
		panic(err)
	}
	defer log.Close()
	defer log.RecoverPanic()

	deugEnv := os.Getenv("DEBUG")
	if deugEnv == "true" || deugEnv == "1" {
		log.TogglePrintStdout(true)
	}

	log.Info("Iniciando aplicação")

	newViperAdapter := config.NewViperAdapter("mj-cli")
	newViperAdapter.SetEnvPrefix("MJ_CLI")

	configRegistry := config.NewConfigRegistry()
	configRegistry.RegisterModule("general", newViperAdapter)

	if err := newViperAdapter.ReadInConfig(); err != nil {
		log.Warn("Falha ao carregar configuração", logger.Metadata{"error": err.Error()})
	} else {
		log.Debug("Configuração carregada")
	}

	lang := newViperAdapter.GetString("lang")
	translator := intl.NewTranslator(lang)
	log.Debug("Tradutor inicializado", logger.Metadata{"lang": lang})

	app := &commands.App{
		Logger:     log,
		Config:     configRegistry,
		Translator: translator,
	}

	cmd.Execute(app)

	if err := newViperAdapter.WriteConfig(); err != nil {
		log.Error("Falha ao salvar configuração", logger.Metadata{"error": err.Error()})
		panic(err)
	}

	log.Debug("Configuração salva")
	log.Info("Aplicação finalizada")
}
