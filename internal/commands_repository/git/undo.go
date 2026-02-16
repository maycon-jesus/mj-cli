package gitcmd

import (
	"context"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/ui"
)

func newGitUndoCommand() *commands.Command {
	return &commands.Command{
		Name:                "undo",
		ShortDescriptionKey: "command.git.undo.short_description",
		Translations: intl.Translations{
			"en": {
				"command.git.undo.short_description":   "Undo the last commit, keeping changes in the working directory",
				"command.git.undo.undoing_last_commit": "Undoing last commit...",
				"command.git.undo.failed_undo":         "Failed to undo last commit: {{error}}",
				"command.git.undo.success":             "Last commit undone successfully!",
			},
			"pt-BR": {
				"command.git.undo.short_description":   "Desfaz o último commit, mantendo as alterações no diretório de trabalho",
				"command.git.undo.undoing_last_commit": "Desfazendo último commit...",
				"command.git.undo.failed_undo":         "Falha ao desfazer último commit: {{error}}",
				"command.git.undo.success":             "Último commit desfeito com sucesso!",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			log := execData.Logger
			term := execData.Terminal
			log.Info("Starting git undo command")

			textUndoing := execData.Translator.T("command.git.undo.undoing_last_commit", map[string]string{})
			term.Println(ui.Title(textUndoing))
			term.StartSpinner("undo", textUndoing)
			gitService := services.NewGitService().WithLogger(log.WithGroup("gitservice"))

			err := gitService.UndoLastCommit()
			if err != nil {
				log.Error("Failed to undo last commit", "error", err)
				term.StopSpinner("undo", execData.Translator.T("command.git.undo.failed_undo", map[string]string{
					"error": err.Error(),
				}))
				return err
			}

			log.Info("Last commit undone successfully")
			term.StopSpinner("undo", execData.Translator.T("command.git.undo.success", map[string]string{}))
			return nil
		},
	}
}
