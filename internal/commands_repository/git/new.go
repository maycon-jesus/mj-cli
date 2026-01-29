package gitcmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	beautyoutput "github.com/maycon-jesus/mj-cli/pkg/beauty_output"
	"github.com/maycon-jesus/mj-cli/pkg/cmd"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newBranchNewCommand() *commands.Command {
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
				"command.git.new.failed_checkout":     "Failed to checkout {{branch}}",
				"command.git.new.failed_pull":         "Failed to update {{branch}}",
				"command.git.new.failed_create":       "Failed to create branch {{branch}}",
				"command.git.new.checkout_success":    "Checked out {{branch}} successfully!",
				"command.git.new.pull_success":        "Updated {{branch}} successfully!",
				"command.git.new.create_success":      "Branch {{branch}} created successfully!",
				"command.git.new.checking_out":        "Checking out {{branch}}...",
				"command.git.new.pulling":             "Updating {{branch}}...",
				"command.git.new.creating_branch":     "Creating branch {{branch}}...",
				"command.git.new.creating_new_branch": "Creating new branch: {{branch}}",
			},
			"pt-BR": {
				"command.git.new.short_description":   "Faz checkout em main/master, atualiza e cria uma nova branch",
				"command.git.new.arg.branch_name":     "O nome da nova branch a ser criada",
				"command.git.new.no_main_or_master":   "Nem a branch 'main' nem 'master' foram encontradas",
				"command.git.new.failed_detect_main":  "Falha ao detectar branch principal: {{error}}",
				"command.git.new.failed_checkout":     "Falha ao fazer checkout em {{branch}}",
				"command.git.new.failed_pull":         "Falha ao atualizar {{branch}}",
				"command.git.new.failed_create":       "Falha ao criar branch {{branch}}",
				"command.git.new.checkout_success":    "Checkout em {{branch}} realizado com sucesso!",
				"command.git.new.pull_success":        "Branch {{branch}} atualizada com sucesso!",
				"command.git.new.create_success":      "Branch {{branch}} criada com sucesso!",
				"command.git.new.checking_out":        "Fazendo checkout em {{branch}}...",
				"command.git.new.pulling":             "Atualizando {{branch}}...",
				"command.git.new.creating_branch":     "Criando branch {{branch}}...",
				"command.git.new.creating_new_branch": "Criando nova branch: {{branch}}",
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			branchName := execData.Args[0]

			// Detecta se a branch principal é main ou master
			mainBranch, err := detectMainBranch()
			if err != nil {
				if errors.Is(err, ErrNoMainBranch) {
					return execData.Translator.Errorf("command.git.new.no_main_or_master", nil)
				}
				return execData.Translator.Errorf("command.git.new.failed_detect_main", map[string]string{"error": err.Error()})
			}

			output := beautyoutput.NewStrBuilder()
			output.SetRealtimeOutput(true)
			output.TitleLine(execData.Translator.T("command.git.new.creating_new_branch", map[string]string{"branch": branchName}))

			// Faz checkout na branch principal
			output.Lambda(fmt.Sprintf("git checkout %s", mainBranch)).NewLine()
			spinner := beautyoutput.NewSpinner(execData.Translator.T("command.git.new.checking_out", map[string]string{"branch": mainBranch}))
			spinner.Start()
			err = cmd.RunCommandWithOptions(fmt.Sprintf("git checkout %s", mainBranch), cmd.CommandOptions{})
			if err != nil {
				spinner.StopWithError(execData.Translator.T("command.git.new.failed_checkout", map[string]string{"branch": mainBranch}))
				return err
			}
			spinner.Stop(execData.Translator.T("command.git.new.checkout_success", map[string]string{"branch": mainBranch}))

			// Atualiza a branch principal
			output.Lambda("git pull").NewLine()
			spinner = beautyoutput.NewSpinner(execData.Translator.T("command.git.new.pulling", map[string]string{"branch": mainBranch}))
			spinner.Start()
			err = cmd.RunCommandWithOptions("git pull", cmd.CommandOptions{})
			if err != nil {
				spinner.StopWithError(execData.Translator.T("command.git.new.failed_pull", map[string]string{"branch": mainBranch}))
				return err
			}
			spinner.Stop(execData.Translator.T("command.git.new.pull_success", map[string]string{"branch": mainBranch}))

			// Cria e faz checkout na nova branch
			output.Lambda(fmt.Sprintf("git checkout -b %s", branchName)).NewLine()
			spinner = beautyoutput.NewSpinner(execData.Translator.T("command.git.new.creating_branch", map[string]string{"branch": branchName}))
			spinner.Start()
			err = cmd.RunCommandWithOptions(fmt.Sprintf("git checkout -b %s", branchName), cmd.CommandOptions{})
			if err != nil {
				spinner.StopWithError(execData.Translator.T("command.git.new.failed_create", map[string]string{"branch": branchName}))
				return err
			}
			spinner.Stop(execData.Translator.T("command.git.new.create_success", map[string]string{"branch": branchName}))

			return nil
		},
	}
}

var ErrNoMainBranch = fmt.Errorf("no main or master branch found")

// detectMainBranch detecta se a branch principal é main ou master
func detectMainBranch() (string, error) {
	// Tenta verificar se existe a branch main
	_, err := cmd.GetCommandOutput("git rev-parse --verify main")
	if err == nil {
		return "main", nil
	}

	// Tenta verificar se existe a branch master
	_, err = cmd.GetCommandOutput("git rev-parse --verify master")
	if err == nil {
		return "master", nil
	}

	// Se nenhuma das duas existir, retorna erro
	return "", ErrNoMainBranch
}
