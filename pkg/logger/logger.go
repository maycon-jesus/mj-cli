package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	Log         *slog.Logger
	FileHandler *FileHandler
}

func NewLoggerComplete(appName string, appVersion string) (Logger, error) {
	logMultiHandler := NewMultiHandler()

	fileHandler, err := NewTemporaryFileHandler(appName)
	if err != nil {
		return Logger{}, fmt.Errorf("failed to create file handler: %w", err)
	}
	fileHandler.SetLevel(slog.LevelInfo)

	debugEnv := os.Getenv("DEBUG")
	if debugEnv != "" {
		var level slog.Level
		if err := level.UnmarshalText([]byte(debugEnv)); err != nil {
			level = slog.LevelDebug
		}
		fileHandler.SetLevel(level)
	}

	logMultiHandler.AddHandler("file", fileHandler)

	logMultiHandler = logMultiHandler.WithAttrs([]slog.Attr{
		slog.String("app", appName),
		slog.String("version", appVersion),
	}).(*MultiHandler)

	log := slog.New(logMultiHandler)
	return Logger{
		Log:         log,
		FileHandler: fileHandler,
	}, nil
}
