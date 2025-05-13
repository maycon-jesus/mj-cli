package commands

import (
	"github.com/maycon-jesus/mj-cli/commands/obsidian"
	"github.com/spf13/cobra"
)

var CmdRoot = &cobra.Command{
	Use:               "mj-cli",
	Short:             "Conjunto de ferramentas CLI do Maycon",
	Long:              `mj-cli é um utilitário CLI desenvolvido para ajudar no dia a dia do Maycon.`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func GetCommandRoot() *cobra.Command {
	CmdRoot.AddCommand(obsidian.GetCommandObsidian())
	CmdRoot.AddCommand(GetSnippetsCommand())
	CmdRoot.AddCommand(GetCommandConfig())
	CmdRoot.AddCommand(GetUuidCommand())
	return CmdRoot
}
