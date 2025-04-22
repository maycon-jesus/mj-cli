package obsidian

import (
	"github.com/spf13/cobra"
)

var Obsidian = &cobra.Command{
	Use:              "obsidian",
	Short:            "Obsidian",
	TraverseChildren: true,
}

func GetCommandObsidian() *cobra.Command {
	var vaultDir string
	Obsidian.AddCommand(GetCommandTagsProperties())
	Obsidian.PersistentFlags().StringVarP(&vaultDir, "vault-dir", "v", "./", "vault directory")
	return Obsidian
}
