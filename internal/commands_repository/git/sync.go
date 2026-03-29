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
				"git.sync.short_description":                        "Sync the current branch with the main branch",
				"command.git.sync.syncing":                          "Syncing branch {{currentBranch}} with {{mainBranch}}",
				"command.git.sync.success":                          "Branch {{currentBranch}} successfully synced with {{mainBranch}} branch",
				"command.git.sync.already_up_to_date":               "Current branch is already up to date with main branch",
				"command.git.sync.failed_get_current_branch":        "Failed to get current branch",
				"command.git.sync.failed_detect_main_branch":        "Failed to detect main branch",
				"command.git.sync.failed_checkout_main_branch":      "Failed to checkout main branch",
				"command.git.sync.failed_checkout_current_branch":   "Failed to checkout current branch",
				"command.git.sync.failed_rebase_current_branch":     "Failed to rebase current branch",
				"command.git.sync.failed_check_new_commits":         "Failed to check if current branch has new commits compared to main branch",
				"command.git.sync.failed_pull_main_branch":          "Failed to pull latest changes on main branch",
				"command.git.sync.stash_message":                    "mj-cli: git sync stash",
				"command.git.sync.failed_stash_changes":             "Failed to stash changes before syncing branches",
				"command.git.sync.failed_pop_stash_changes":         "Failed to unstash changes after syncing branches",
				"command.git.sync.failed_check_uncommitted_changes": "Failed to check for uncommitted changes",
				"command.git.sync.resolve_conflicts_hint":           "Failed rebase/merge detected. Please resolve the conflicts manually, commit the changes and then run 'git merge/rebase --continue' to continue the rebase/merge process.",
				"command.git.sync.stash_pop_hint":                   "There are stashed changes that need to be reapplied. Please run 'git stash pop' to reapply the stashed changes.",
			},
			"pt-BR": {
				"git.sync.short_description":                        "Sincroniza a branch atual com a branch principal",
				"command.git.sync.syncing":                          "Sincronizando branch {{currentBranch}} com {{mainBranch}}",
				"command.git.sync.success":                          "Branch {{currentBranch}} sincronizada com a branch {{mainBranch}} com sucesso",
				"command.git.sync.already_up_to_date":               "A branch atual já está atualizada com a branch principal",
				"command.git.sync.failed_get_current_branch":        "Falha ao obter a branch atual",
				"command.git.sync.failed_detect_main_branch":        "Falha ao detectar a branch principal",
				"command.git.sync.failed_checkout_main_branch":      "Falha ao fazer checkout da branch principal",
				"command.git.sync.failed_checkout_current_branch":   "Falha ao fazer checkout da branch atual",
				"command.git.sync.failed_rebase_current_branch":     "Falha ao rebasear a branch atual",
				"command.git.sync.failed_check_new_commits":         "Falha ao verificar se a branch atual tem novos commits",
				"command.git.sync.failed_pull_main_branch":          "Falha ao puxar as últimas mudanças da branch principal",
				"command.git.sync.stash_message":                    "mj-cli: stash do git sync",
				"command.git.sync.failed_stash_changes":             "Falha ao stashear as mudanças antes de sincronizar as branches",
				"command.git.sync.failed_pop_stash_changes":         "Falha ao unstashear as mudanças depois de sincronizar as branches",
				"command.git.sync.failed_check_uncommitted_changes": "Falha ao verificar por mudanças não comitadas",
				"command.git.sync.resolve_conflicts_hint":           "Rebase/merge com falha detectado. Por favor, resolva os conflitos manualmente, faça o commit das mudanças e então rode 'git merge/rebase --continue' para continuar o processo de rebase/merge.",
				"command.git.sync.stash_pop_hint":                   "Há mudanças stashed que precisam ser reaplicadas. Por favor, execute 'git stash pop' para reaplicar as mudanças stashed.",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			log := execData.Logger
			term := execData.Terminal
			t := execData.Translator.T
			tErr := execData.Translator.Errorf

			log.Info("Starting git sync process")
			gitService := services.NewGitService().WithLogger(execData.Logger.WithGroup("gitservice"))

			onMain, stashed := false, false

			term.StartSpinner("sync", t("command.git.sync.syncing", map[string]string{
				"currentBranch": "...",
				"mainBranch":    "...",
			}))

			currentBranch, err := gitService.GetCurrentBranch()
			if err != nil {
				execData.Logger.Error("Failed to get current branch", "error", err)
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_get_current_branch", map[string]string{}))
				return tErr("command.git.sync.failed_get_current_branch", map[string]string{})
			}

			mainBranch, err := gitService.DetectMainBranch()
			if err != nil {
				execData.Logger.Error("Failed to detect main branch", "error", err)
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_detect_main_branch", map[string]string{}))
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

			// check if there are uncommitted changes
			hasUncommittedChanges, err := gitService.HasTrackedUncommitedChanges()
			if err != nil {
				execData.Logger.Error("Failed to check for uncommitted changes", "error", err)
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_check_uncommitted_changes", map[string]string{}))
				return tErr("command.git.sync.failed_check_uncommitted_changes", map[string]string{})
			}

			// rollback
			defer func() {
				if onMain {
					err = gitService.Checkout(currentBranch)
					if err != nil {
						execData.Logger.Error("Failed to checkout back to current branch during rollback", "branch", currentBranch, "error", err.Error())
						term.Println(ui.Info(t("command.git.sync.failed_checkout_current_branch", map[string]string{})))
						if stashed {
							term.Println(ui.Info(t("command.git.sync.stash_pop_hint", map[string]string{})))
						}
						return
					}
				}

				if stashed {
					err = unstashChanges(ctx, execData, gitService)
					if err != nil {
						term.Println(ui.Info(t("command.git.sync.stash_pop_hint", map[string]string{})))
						return
					}
				}
			}()

			// stash push
			if hasUncommittedChanges {
				term.Println(ui.Lambdaf("git stash push -m \"%s\"", t("command.git.sync.stash_message", map[string]string{})))
				err = gitService.StashPush(t("command.git.sync.stash_message", map[string]string{}))
				if err != nil {
					execData.Logger.Error("Failed to stash changes", "error", err)
					term.StopSpinnerWithError("sync", t("command.git.sync.failed_stash_changes", map[string]string{}))
					return tErr("command.git.sync.failed_stash_changes", map[string]string{})
				}
				stashed = true
				execData.Logger.Info("Uncommitted changes stashed successfully")
			}

			// checkout on main
			term.Println(ui.Lambdaf("git checkout %s", mainBranch))
			err = gitService.Checkout(mainBranch)
			if err != nil {
				execData.Logger.Error("Failed to checkout main branch", "branch", mainBranch, "error", err.Error())
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_checkout_main_branch", map[string]string{}))
				return tErr("command.git.sync.failed_checkout_main_branch", map[string]string{})
			}
			onMain = true

			// pull on main
			term.Println(ui.Lambdaf("git pull"))
			err = gitService.Pull()
			if err != nil {
				execData.Logger.Error("Failed to pull latest changes on main branch", "branch", mainBranch, "error", err.Error())
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_pull_main_branch", map[string]string{}))
				return tErr("command.git.sync.failed_pull_main_branch", map[string]string{})
			}

			// checkout back to current branch
			term.Println(ui.Lambdaf("git checkout %s", currentBranch))
			err = gitService.Checkout(currentBranch)
			if err != nil {
				execData.Logger.Error("Failed to checkout back to current branch", "branch", currentBranch, "error", err.Error())
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_checkout_current_branch", map[string]string{}))
				return tErr("command.git.sync.failed_checkout_current_branch", map[string]string{})
			}
			onMain = false

			// check if current branch has new commits compared to main branch
			hasNewCommits, err := gitService.BranchHasNewCommitsFor(mainBranch)
			if err != nil {
				execData.Logger.Error("Failed to check if current branch has new commits compared to main branch", "branch", currentBranch, "mainBranch", mainBranch, "error", err.Error())
				term.StopSpinnerWithError("sync", t("command.git.sync.failed_check_new_commits", map[string]string{}))
				return tErr("command.git.sync.failed_check_new_commits", map[string]string{})
			}

			if !hasNewCommits {
				execData.Logger.Info("Current branch is already up to date with main branch", "branch", currentBranch, "mainBranch", mainBranch)
				term.Println(ui.Info(t("command.git.sync.already_up_to_date", map[string]string{
					"currentBranch": currentBranch,
					"mainBranch":    mainBranch,
				})))

				// unstash changes if there are stashed changes
				if stashed {
					err = unstashChanges(ctx, execData, gitService)
					if err != nil {
						log.Error("Failed to unstash changes after syncing with main branch", "error", err.Error())
					} else {
						log.Info("Stashed changes unstashed successfully")
						stashed = false
					}
				}

				term.StopSpinner("sync", t("command.git.sync.success", map[string]string{
					"currentBranch": currentBranch,
					"mainBranch":    mainBranch,
				}))
				return nil
			}

			// rebase current branch with main branch
			term.Println(ui.Lambdaf("git rebase %s", mainBranch))
			err = gitService.Rebase(mainBranch)
			if err != nil {
				execData.Logger.Error("Failed to rebase current branch with main branch", "branch", currentBranch, "mainBranch", mainBranch, "error", err.Error())

				// abort rebase if it failed
				abortErr := gitService.RebaseAbort()
				if abortErr != nil {
					execData.Logger.Error("Failed to abort rebase after failed rebase attempt", "branch", currentBranch, "mainBranch", mainBranch, "error", abortErr.Error())
					term.Println(ui.Info(t("command.git.sync.resolve_conflicts_hint", map[string]string{})))
				}

				term.StopSpinnerWithError("sync", t("command.git.sync.failed_rebase_current_branch", map[string]string{}))
				return tErr("command.git.sync.failed_rebase_current_branch", map[string]string{})
			}

			// unstash changes if there are stashed changes
			if stashed {
				err = unstashChanges(ctx, execData, gitService)
				if err != nil {
					log.Error("Failed to unstash changes after rebasing", "error", err.Error())
				} else {
					execData.Logger.Info("Stashed changes unstashed successfully")
					stashed = false
				}
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

func unstashChanges(ctx context.Context, execData *commands.ExecData, gitService *services.GitService) error {
	term := execData.Terminal
	tErr := execData.Translator.Errorf

	term.Println(ui.Lambdaf("git stash pop"))
	err := gitService.StashPop()
	if err != nil {
		execData.Logger.Error("Failed to pop stashed changes", "error", err)
		return tErr("command.git.sync.failed_pop_stash_changes", map[string]string{})
	}
	return nil
}
