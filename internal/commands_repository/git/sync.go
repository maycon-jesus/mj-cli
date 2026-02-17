package gitcmd

import (
	"context"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/ui"
)

func newGitSyncCommand() *commands.Command {
	return &commands.Command{
		Name:                "sync",
		ShortDescriptionKey: "git.sync.short_description",
		Translations: intl.Translations{
			"en": {
				"git.sync.short_description":                      "Sync the current branch with the main branch",
				"command.git.sync.syncing":                        "Syncing branch {{currentBranch}} with {{mainBranch}}",
				"command.git.sync.success":                        "Branch {{currentBranch}} successfully synced with {{mainBranch}} branch",
				"command.git.sync.already_up_to_date":             "Current branch is already up to date with main branch",
				"command.git.sync.failed_get_current_branch":      "Failed to get current branch",
				"command.git.sync.failed_detect_main_branch":      "Failed to detect main branch",
				"command.git.sync.failed_checkout_main_branch":    "Failed to checkout main branch",
				"command.git.sync.failed_checkout_current_branch": "Failed to checkout current branch",
				"command.git.sync.failed_rebase_current_branch":   "Failed to rebase current branch",
				"command.git.sync.failed_check_new_commits":       "Failed to check if current branch has new commits compared to main branch",
				"command.git.sync.failed_pull_main_branch":        "Failed to pull latest changes on main branch",
			},
			"pt-BR": {
				"git.sync.short_description":                      "Sincroniza a branch atual com a branch principal",
				"command.git.sync.syncing":                        "Sincronizando branch {{currentBranch}} com {{mainBranch}}",
				"command.git.sync.success":                        "Branch {{currentBranch}} sincronizada com a branch {{mainBranch}} com sucesso",
				"command.git.sync.already_up_to_date":             "A branch atual já está atualizada com a branch principal",
				"command.git.sync.failed_get_current_branch":      "Falha ao obter a branch atual",
				"command.git.sync.failed_detect_main_branch":      "Falha ao detectar a branch principal",
				"command.git.sync.failed_checkout_main_branch":    "Falha ao fazer checkout da branch principal",
				"command.git.sync.failed_checkout_current_branch": "Falha ao fazer checkout da branch atual",
				"command.git.sync.failed_rebase_current_branch":   "Falha ao rebasear a branch atual",
				"command.git.sync.failed_check_new_commits":       "Falha ao verificar se a branch atual tem novos commits",
				"command.git.sync.failed_pull_main_branch":        "Falha ao puxar as últimas mudanças da branch principal",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			log := execData.Logger
			term := execData.Terminal
			t := execData.Translator.T
			tErr := execData.Translator.Errorf

			log.Info("Starting git sync process")
			gitService := services.NewGitService().WithLogger(execData.Logger.WithGroup("gitservice"))

			term.StartSpinner("sync", t("command.git.sync.syncing", map[string]string{
				"currentBranch": "...",
				"mainBranch":    "...",
			}))

			currentBranch, err := gitService.GetCurrentBranch()
			if err != nil {
				execData.Logger.Error("Failed to get current branch", "error", err)
				return tErr("command.git.sync.failed_get_current_branch", map[string]string{})
			}

			mainBranch, err := gitService.DetectMainBranch()
			if err != nil {
				execData.Logger.Error("Failed to detect main branch", "error", err)
				return tErr("command.git.sync.failed_detect_main_branch", map[string]string{})
			}

			term.UpdateSpinnerMessage("sync", t("command.git.sync.syncing", map[string]string{
				"currentBranch": currentBranch,
				"mainBranch":    mainBranch,
			}))

			term.Println(ui.Title(t("command.git.sync.syncing", map[string]string{
				"currentBranch": currentBranch,
				"mainBranch":    mainBranch,
			})))
			execData.Logger.Info("Syncing current branch with main branch", "mainBranch", mainBranch, "currentBranch", currentBranch)

			term.Println(ui.Lambdaf("git checkout %s", mainBranch))
			err = gitService.Checkout(mainBranch)
			if err != nil {
				execData.Logger.Error("Failed to checkout main branch", "branch", mainBranch, "error", err.Error())
				return tErr("command.git.sync.failed_checkout_main_branch", map[string]string{})
			}

			term.Println(ui.Lambdaf("git pull"))
			err = gitService.Pull()
			if err != nil {
				execData.Logger.Error("Failed to pull latest changes on main branch", "branch", mainBranch, "error", err.Error())
				return tErr("command.git.sync.failed_pull_main_branch", map[string]string{})
			}

			term.Println(ui.Lambdaf("git checkout %s", currentBranch))
			err = gitService.Checkout(currentBranch)
			if err != nil {
				execData.Logger.Error("Failed to checkout back to current branch", "branch", currentBranch, "error", err.Error())
				return tErr("command.git.sync.failed_checkout_current_branch", map[string]string{})
			}

			hasNewCommits, err := gitService.BranchHasNewCommitsFor(mainBranch)
			if err != nil {
				execData.Logger.Error("Failed to check if current branch has new commits compared to main branch", "branch", currentBranch, "mainBranch", mainBranch, "error", err.Error())
				return tErr("command.git.sync.failed_check_new_commits", map[string]string{})
			}

			if !hasNewCommits {
				execData.Logger.Info("Current branch is already up to date with main branch", "branch", currentBranch, "mainBranch", mainBranch)
				term.Println(ui.Info(t("command.git.sync.already_up_to_date", map[string]string{
					"currentBranch": currentBranch,
					"mainBranch":    mainBranch,
				})))
				term.StopSpinner("sync", t("command.git.sync.success", map[string]string{
					"currentBranch": currentBranch,
					"mainBranch":    mainBranch,
				}))
				return nil
			}

			term.Println(ui.Lambdaf("git rebase %s", mainBranch))
			err = gitService.Rebase(mainBranch)
			if err != nil {
				execData.Logger.Error("Failed to rebase current branch with main branch", "branch", currentBranch, "mainBranch", mainBranch, "error", err.Error())
				return tErr("command.git.sync.failed_rebase_current_branch", map[string]string{})
			}

			term.StopSpinner("sync", t("command.git.sync.success", map[string]string{
				"currentBranch": currentBranch,
				"mainBranch":    mainBranch,
			}))
			execData.Logger.Info("Current branch synced with main branch successfully")

			return nil
		},
	}
}
