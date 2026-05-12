package commandsrepository

import (
	"github.com/maycon-jesus/mj-cli/internal/commands"
	alias_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/alias"
	config_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/config"
	cube_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/cube"
	gitcmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/git"
	log_cmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/log"
	semvercmd "github.com/maycon-jesus/mj-cli/internal/commands_repository/semver"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func NewRootCommand() *commands.Command {
	return &commands.Command{
		Name:                "mj-cli",
		ShortDescriptionKey: "command.root.short_description",
		LongDescriptionKey:  "command.root.long_description",
		RunHelpOnNoArgs:     true,
		SubCommands: []*commands.Command{
			config_cmd.NewConfigCommand(),
			alias_cmd.NewAliasCommand(),
			gitcmd.NewGitCommand(),
			log_cmd.NewLogCommand(),
			semvercmd.NewSemverCommand(),
			cube_cmd.NewCubeCommand(),
		},
		Translations: intl.Translations{
			"en": {
				"command.root.short_description": "MJ CLI - A command line tool",
				"command.root.long_description":  "MJ CLI is a powerful command line interface tool to interact with MJ services.",
			},
			"pt-BR": {
				"command.root.short_description": "MJ CLI - Uma ferramenta de linha de comando",
				"command.root.long_description":  "MJ CLI é uma poderosa ferramenta de linha de comando para interagir com os serviços MJ.",
			},
		},
	}
}
