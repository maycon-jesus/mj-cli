package gitcmd

import (
	"context"
	"errors"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/ui"
)

func newGitNewCommand() *commands.Command {
	return &commands.Command{
		Name:                "new",
		ShortDescriptionKey: "command.git.new.short_description",
		Args: []commands.Arg{
			{Name: "branch-name", DescriptionKey: "command.git.new.arg.branch_name", Required: true},
		},
		Translations: intl.Translations{
			"en": {
				"command.git.new.short_description":   "Checkout to main/master, update it and create a new branch",
				"command.git.new.arg.branch_name":     "The name of the new branch to create",
				"command.git.new.no_main_or_master":   "Neither 'main' nor 'master' branch found",
				"command.git.new.failed_detect_main":  "Failed to detect main branch: {{error}}",
				"command.git.new.create_error":        "Failed to create branch {{branch}}",
				"command.git.new.create_success":      "Branch {{branch}} created successfully!",
				"command.git.new.creating_branch":     "Creating branch {{branch}}...",
				"command.git.new.creating_new_branch": "Creating new branch: {{branch}}",
				"command.git.new.failed_checkout":     "Failed to checkout branch {{branch}}",
				"command.git.new.failed_pull":         "Failed to pull latest changes on branch {{branch}}",
			},
			"pt-BR": {
				"command.git.new.short_description":   "Faz checkout em main/master, atualiza e cria uma nova branch",
				"command.git.new.arg.branch_name":     "O nome da nova branch a ser criada",
				"command.git.new.no_main_or_master":   "Nem a branch 'main' nem 'master' foram encontradas",
				"command.git.new.failed_detect_main":  "Falha ao detectar branch principal: {{error}}",
				"command.git.new.create_error":        "Falha ao criar branch {{branch}}",
				"command.git.new.create_success":      "Branch {{branch}} criada com sucesso!",
				"command.git.new.creating_branch":     "Criando branch {{branch}}...",
				"command.git.new.creating_new_branch": "Criando nova branch: {{branch}}",
				"command.git.new.failed_checkout":     "Falha ao fazer checkout da branch {{branch}}",
				"command.git.new.failed_pull":         "Falha ao puxar últimas mudanças da branch {{branch}}",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			log := execData.Logger
			term := execData.Terminal
			translator := execData.Translator

			branchName := execData.Args[0]
			log.Info("Starting git new command", "branch", branchName)

			gitService := services.NewGitService().WithLogger(log.WithGroup("gitservice"))

			// Detecta se a branch principal é main ou master
			mainBranch, err := gitService.DetectMainBranch()
			if err != nil {
				if errors.Is(err, services.ErrNoMainBranch) {
					log.Error("No main or master branch found")
					return translator.Errorf("command.git.new.no_main_or_master", nil)
				}
				log.Error("Failed to detect main branch", "error", err)
				return translator.Errorf("command.git.new.failed_detect_main", map[string]string{"error": err.Error()})
			}
			log.Debug("Main branch detected", "main_branch", mainBranch)

			// Checkout main branch
			term.StartSpinner("checkout", translator.T("command.git.new.creating_branch", map[string]string{"branch": mainBranch}))
			err = gitService.Checkout(mainBranch)
			if err != nil {
				log.Error("Failed to checkout main branch", "main_branch", mainBranch, "error", err)
				term.StopSpinnerWithError("checkout", translator.T("command.git.new.failed_checkout", map[string]string{"branch": mainBranch}))
				return err
			}

			// Puxa as últimas mudanças da branch principal
			err = gitService.Pull()
			if err != nil {
				log.Error("Failed to pull latest changes on main branch", "main_branch", mainBranch, "error", err)
				term.StopSpinnerWithError("checkout", translator.T("command.git.new.failed_pull", map[string]string{"branch": mainBranch}))
				return err
			}

			// Cria a nova branch a partir da principal
			term.Println(ui.Title(translator.T("command.git.new.creating_branch", map[string]string{"branch": branchName})))
			err = gitService.NewBranch(branchName)
			if err != nil {
				log.Error("Failed to create branch", "branch", branchName, "error", err)
				term.StopSpinnerWithError("checkout", translator.T("command.git.new.create_error", map[string]string{"branch": branchName}))
				return err
			}
			log.Info("Branch created successfully", "branch", branchName, "from", mainBranch)
			term.StopSpinner("checkout", translator.T("command.git.new.create_success", map[string]string{"branch": branchName}))
			return nil
		},
	}
}
