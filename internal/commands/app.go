package commands

import (
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
	"github.com/maycon-jesus/mj-cli/pkg/mjterm"
)

type App struct {
	Logger     *logger.Logger
	Config     *config.ConfigRegistry
	Translator *intl.Translator
	Terminal   *mjterm.Terminal
	Database   *services.DatabaseService
	Version    string
	Name       string
}
