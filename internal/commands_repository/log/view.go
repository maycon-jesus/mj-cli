package log_cmd

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/tint"
)

var StyleTime = tint.NewStyle().Foreground(tint.BrightBlack)
var StyleLevelDebug = tint.NewStyle().Foreground(tint.BrightBlue)
var StyleLevelInfo = tint.NewStyle().Foreground(tint.BrightGreen)
var StyleLevelWarn = tint.NewStyle().Foreground(tint.BrightYellow)
var StyleLevelError = tint.NewStyle().Foreground(tint.BrightRed)
var StyleLevel = map[string]*tint.Style{
	slog.LevelDebug.String(): StyleLevelDebug,
	slog.LevelInfo.String():  StyleLevelInfo,
	slog.LevelWarn.String():  StyleLevelWarn,
	slog.LevelError.String(): StyleLevelError,
}

func NewLogViewCommand() *commands.Command {
	return &commands.Command{
		Name:                "view",
		ShortDescriptionKey: "command.log.view.short_description",
		RunHelpOnNoArgs:     true,
		Translations: map[string]map[string]string{
			"en": {
				"command.log.view.short_description": "View the content of a log file",
				"command.log.arg.file_path":          "The path of the log file to be displayed",
				"command.log.arg.file_path.required": "The file_path argument is required",
			},
			"pt-BR": {
				"command.log.view.short_description": "Visualiza o conteúdo de um arquivo de log",
				"command.log.arg.file_path":          "O caminho do arquivo de log a ser exibido",
				"command.log.arg.file_path.required": "O argumento file_path é obrigatório",
			},
		},
		Args: []commands.Arg{
			{
				Name:           "file_path",
				DescriptionKey: "command.log.arg.file_path",
				Required:       true,
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			log := execData.Logger
			log.Info("Executing log view command", "args", execData.Args)
			if len(execData.Args) == 0 {
				return errors.New(execData.Translator.T("command.log.arg.file_path.required", map[string]string{}))
			}

			filePath := execData.Args[0]
			if filePath == "" {
				return errors.New(execData.Translator.T("command.log.arg.file_path.required", map[string]string{}))
			}

			logService := services.NewLogService().WithLogger(execData.Logger.WithGroup("logservice"))
			ch, err := logService.GetLogs(filePath)
			if err != nil {
				return err
			}
			for logLine := range ch {
				if logLine.Error != nil {
					execData.Terminal.Printf("Error reading log line: %v\n", logLine.Error)
					continue
				}
				line := logLine.Log

				timeIso := line["time"].(string)
				level := line["level"].(string)
				msg := line["message"].(string)

				timeParsed, err := time.Parse(time.RFC3339, timeIso)
				timeParsed = timeParsed.Local()
				if err != nil {
					execData.Terminal.Printf("Error parsing time: %v\n", err)
					continue
				}

				timeRendered := StyleTime.Render(timeParsed.Format("2006-01-02 15:04:05"))
				levelStyle, ok := StyleLevel[level]
				if !ok {
					levelStyle = StyleLevelInfo
				}
				levelRendered := levelStyle.Render(level)

				execData.Terminal.Printf("[%s] [%s] %s\n", timeRendered, levelRendered, msg)

				delete(line, "time")
				delete(line, "level")
				delete(line, "message")

				if len(line) == 0 {
					continue
				}
				jsonLine, err := json.MarshalIndent(line, "", "  ")
				if err != nil {
					execData.Terminal.Printf("Error marshaling log line to JSON: %v\n", err)
					continue
				}

				execData.Terminal.Println(string(jsonLine))
			}

			return nil
		},
	}
}
