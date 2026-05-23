package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/maycon-jesus/mj-cli/cmd"
	"github.com/maycon-jesus/mj-cli/internal/commands"
	commandsrepository "github.com/maycon-jesus/mj-cli/internal/commands_repository"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/config"
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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	appDataDir := filepath.Join(homeDir, "."+appName)
	if err := os.MkdirAll(appDataDir, os.ModePerm); err != nil {
		panic(err)
	}
	log.Debug("appDataDir setado", "appDataDir", appDataDir)

	configPath := filepath.Join(appDataDir, "config.yaml")
	yamlPersister := config.NewYamlModule(configPath)
	if err := yamlPersister.Load(); err != nil {
		log.Warn("Falha ao carregar configuração", "error", err.Error())
	} else {
		log.Debug("Configuração carregada", "path", configPath)
	}

	configManager := config.NewConfigManager().WithPersister(yamlPersister)

	if err := configManager.AddEntry("lang", "Idioma da UI", "pt-BR"); err != nil {
		log.Warn("Falha ao registrar entrada de configuração", "key", "lang", "error", err.Error())
	}

	langValue, _ := configManager.Get("lang")
	lang, _ := langValue.(string)
	translator := intl.NewTranslator(lang)
	log.Debug("Tradutor inicializado", "lang", lang)

	rootCmd := commandsrepository.NewRootCommand()

	app := &commands.App{
		Logger:     &logger,
		Config:     configManager,
		Translator: translator,
		Terminal:   term,
		Version:    appVersion,
		Name:       appName,
		Database:   services.NewDatabaseService(appDataDir, "app.db").WithLogger(logger.Log.WithGroup("database")),
		RootCmd:    rootCmd,
	}

	cmd.Execute(app)

	if err := configManager.Save(); err != nil {
		log.Warn("Falha ao salvar configuração", "error", err.Error())
	} else {
		log.Debug("Configuração salva")
	}
	log.Info("Aplicação finalizada")
}
