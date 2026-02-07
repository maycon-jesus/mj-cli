package commands

import (
	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
)

type App struct {
	Logger     *logger.Logger
	Config     *config.ConfigRegistry
	Translator *intl.Translator
}
