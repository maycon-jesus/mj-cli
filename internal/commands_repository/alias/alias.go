package alias_cmd

import "github.com/maycon-jesus/mj-cli/internal/commands"

func NewAliasCommand() *commands.Command {
	return &commands.Command{
		Name:             "alias",
		ShortDescription: "Manage command aliases",
		RunHelpOnNoArgs:  true,
		SubCommands:      []*commands.Command{newAliasSetCommand(), newAliasRunCommand(), newAliasLsCommand(), newAliasViewCommand(), newAliasRmCommand()},
	}
}
