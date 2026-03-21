package semvercmd

import (
	"context"
	"os"

	"github.com/maycon-jesus/mj-cli/internal/commands"
	"github.com/maycon-jesus/mj-cli/internal/services"
	"github.com/maycon-jesus/mj-cli/pkg/filesystem"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
)

func newSemverUpgradeCommand() *commands.Command {
	return &commands.Command{
		Name:                "upgrade",
		ShortDescriptionKey: "command.semver.upgrade.short_description",
		Translations: intl.Translations{
			"en": {
				"command.semver.upgrade.short_description":     "Upgrade the project's semantic version",
				"command.semver.upgrade.args.path.description": "Path to the project",
				"command.semver.upgrade.args.type.description": "Upgrade type (major, minor, patch)",
			},
			"pt-BR": {
				"command.semver.upgrade.short_description":     "Atualizar a versão semântica do projeto",
				"command.semver.upgrade.args.path.description": "Caminho para o projeto",
				"command.semver.upgrade.args.type.description": "Tipo de atualização (major, minor, patch)",
			},
		},
		Args: []commands.Arg{
			{
				Name:           "path",
				DescriptionKey: "command.semver.upgrade.args.path.description",
				Required:       true,
			},
			{
				Name:           "type",
				DescriptionKey: "command.semver.upgrade.args.type.description",
				Required:       true,
			},
		},
		Handler: func(ctx context.Context, execData *commands.ExecData) error {
			execData.Logger.Info("Executing semver upgrade command")

			codeService := services.NewCodeService().WithLogger(execData.Logger.WithGroup("code-service"))
			semverService := services.NewSemverService()

			path, err := filesystem.GetAbsolutePath(execData.Args[0])
			if err != nil {
				execData.Logger.Error(err.Error())
				return err
			}

			currentVersion, err := codeService.GetVersion(path)
			if err != nil {
				execData.Logger.Error(err.Error())
				return err
			}

			upgradeType := execData.Args[1]

			version, err := semverService.ParseVersion(currentVersion)

			switch upgradeType {
			case "major":
				version.IncrementMajor()
			case "minor":
				version.IncrementMinor()
			case "patch":
				version.IncrementPatch()
			default:
				return os.ErrInvalid
			}

			codeService.SetVersion(path, version.String())

			return nil
		},
	}
}
