package obsidian

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Obsidian = &cobra.Command{
	Use:              "obsidian",
	Short:            "Utilitários para o obsidian",
	Aliases:          []string{"ob"},
	TraverseChildren: true,
}

func GetCommandObsidian() *cobra.Command {
	var vaultDir string
	Obsidian.AddCommand(GetCommandTagsProperties())
	Obsidian.PersistentFlags().StringVarP(&vaultDir, "vault-dir", "v", "", "vault directory")

	if val := viper.GetString("obsidianVaultDir"); val != "" {
		Obsidian.PersistentFlags().Lookup("vault-dir").Value.Set(val)
	} else {
		Obsidian.MarkPersistentFlagRequired("vault-dir")
	}
	return Obsidian
}
