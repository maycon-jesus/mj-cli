// Package logger provides JSON-structured logging built on [log/slog].
//
// It ships two custom [slog.Handler] implementations — [FileHandler] and
// [MultiHandler] — that can be composed to fan out log records to multiple
// destinations.
//
// The typical setup is done through [NewLoggerComplete], which creates a
// temporary log file, wires up a [MultiHandler], and respects the DEBUG
// environment variable for level overrides.
//
//	lgr, err := logger.NewLoggerComplete("mj-cli", "1.0.0")
//	if err != nil { ... }
//	defer lgr.FileHandler.Close()
//	lgr.Log.Info("application started")
package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger wraps a [slog.Logger] together with the underlying [FileHandler],
// so callers can access the log file path or close the file when done.
type Logger struct {
	Log         *slog.Logger
	FileHandler *FileHandler
}

// NewLoggerComplete creates a fully-configured [Logger].
//
// It creates a temporary log file named after appName (via [NewTemporaryFileHandler]),
// defaults to [slog.LevelDebug], and attaches "app" and "version"
// attributes to every record.
//
// The DEBUG environment variable can override the file handler level.
// If the value is a valid [slog.Level] string (e.g. "debug", "warn"), that level
// is used; otherwise the level defaults to [slog.LevelDebug].
func NewLoggerComplete(appName string, appVersion string) (Logger, error) {
	logMultiHandler := NewMultiHandler()

	fileHandler, err := NewTemporaryFileHandler(appName)
	if err != nil {
		return Logger{}, fmt.Errorf("failed to create file handler: %w", err)
	}

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
