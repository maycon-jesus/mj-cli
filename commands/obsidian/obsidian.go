package obsidian

import (
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var Obsidian = &cobra.Command{
	Use:              "obsidian",
	Short:            "Utilitários para o obsidian",
	Aliases:          []string{"ob"},
	PersistentPreRun: CommandObsidianValidadePersistentFlags,
	TraverseChildren: true,
}

func CommandObsidianValidadePersistentFlags(cmd *cobra.Command, args []string) {
	vaultDir, err := cmd.Flags().GetString("vault-dir")
	cobra.CheckErr(err)

	wdDir, err := os.Getwd()
	cobra.CheckErr(err)

	dirNormalized, err := utils.NormalizePath(wdDir, vaultDir)
	cobra.CheckErr(err)

	err = cmd.Flags().Lookup("vault-dir").Value.Set(dirNormalized)
	cobra.CheckErr(err)
}

func GetCommandObsidian() *cobra.Command {
	Obsidian.AddCommand(GetCommandTagsProperties())
	Obsidian.AddCommand(GetCommandWeek())
	Obsidian.AddCommand(GetCommandMonth())
	Obsidian.AddCommand(GetCommandDaily())

	vaultDir := viper.GetString("obsidian-vault-dir")
	Obsidian.PersistentFlags().String("vault-dir", vaultDir, "Vault directory\nConfig key: obsidian-vault-dir\n")
	if vaultDir == "" {
		Obsidian.MarkPersistentFlagRequired("vault-dir")
	}

	return Obsidian
}
