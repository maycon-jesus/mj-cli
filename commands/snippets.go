package commands

import (
	"github.com/maycon-jesus/mj-cli/commands/snippets"
	"github.com/spf13/cobra"
)

var SnippetsCommand = &cobra.Command{
	Use:     "snippets",
	Short:   "Snippets",
	Aliases: []string{"sn"},
}

func GetSnippetsCommand() *cobra.Command {
	SnippetsCommand.AddCommand(snippets.GetBulletJournalSnippetCommand())
	SnippetsCommand.AddCommand(snippets.GetPortForwardAutoblocsApiSnippetCommand())

	return SnippetsCommand
}
