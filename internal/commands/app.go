package commands

import (
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
	"github.com/maycon-jesus/mj-cli/pkg/mjterm"
)

type App struct {
	Logger     *logger.Logger
	Config     *config.ConfigManager
	Translator *intl.Translator
	Terminal   *mjterm.Terminal
	Database   *services.DatabaseService
	RootCmd    *Command
	Version    string
	Name       string
}
