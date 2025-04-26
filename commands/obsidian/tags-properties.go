package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var TagsProperties = &cobra.Command{
	Use:   "tags-properties",
	Short: "tags-properties short",
	Run:   run,
}

func GetCommandTagsProperties() *cobra.Command {
	TagsProperties.Flags().StringP("templates-dir", "t", "99 - Meta/02 - Tags", "tags templates directory")

	if val := viper.GetString("obsidianTagsDir"); val != "" {
		TagsProperties.Flags().Lookup("templates-dir").Value.Set(val)
	} else {
		TagsProperties.MarkFlagRequired("templates-dir")
	}

	return TagsProperties
}

func run(cmd *cobra.Command, args []string) {
	vaultDir, _ := cmd.Flags().GetString("vault-dir")
	templatesDir, _ := cmd.Flags().GetString("templates-dir")
	var vaultDirAbs string
	var templatesDirAbs string

	if !filepath.IsAbs(vaultDir) {
		wdDir, _ := os.Getwd()
		vaultDirAbs = filepath.Join(wdDir, vaultDir)
	} else {
		vaultDirAbs = vaultDir
	}

	if !filepath.IsAbs(templatesDir) {
		templatesDirAbs = filepath.Join(vaultDirAbs, templatesDir)
	} else {
		templatesDirAbs = templatesDir
	}

	vault := obsidian.NewVault(vaultDirAbs, templatesDirAbs)
	vault.LoadAllFiles()

	for _, file := range vault.Notes {
		values, ok := file.GetPropertyValues("tags")
		if !ok {
			continue
		}
		fileUpdated := false

		for _, tag := range values {
			tagTemplateNote, ok := vault.GetTagTemplateNote(tag)
			if !ok {
				continue
			}

			for tagTemplateKey, tagTemplateValue := range tagTemplateNote.Frontmatter {
				_, ok = file.GetProperty(tagTemplateKey)
				if !ok {
					file.AddProperty(tagTemplateKey, tagTemplateValue.GetValues())
					fileUpdated = true
				}
			}
		}
		if fileUpdated {
			relPath, _ := filepath.Rel(vault.Path, file.Path)
			fmt.Println("Arquivo atualizado:", relPath)
			file.WriteFile()
		}
	}

}
