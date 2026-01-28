package commandsrepository

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
)

func NewRootCommand() *commands.Command {
	return &commands.Command{
		Name:             "mj-cli",
		ShortDescription: "MJ CLI - Uma ferramenta de linha de comando",
		LongDescription:  `MJ CLI é uma ferramenta de linha de comando construída com uma arquitetura modular e extensível.`,
		RunHelpOnNoArgs:  true,
	}
}
