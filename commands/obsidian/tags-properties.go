package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var TagsProperties = &cobra.Command{
	Use:     "tags-properties",
	Short:   "tags-properties short",
	Aliases: []string{"format-tags"},
	Run:     run,
}

func GetCommandTagsProperties() *cobra.Command {
	TagsProperties.Flags().StringP("templates-dir", "t", "", "tags templates directory")

	if val := viper.GetString("obsidianTagsDir"); val != "" {
		TagsProperties.Flags().Lookup("templates-dir").Value.Set(val)
	} else {
		TagsProperties.MarkFlagRequired("templates-dir")
	}

	return TagsProperties
}

func analyzeFile(ch chan<- []string, wg *sync.WaitGroup, vault *obsidian.Vault, file *obsidian.ObsidianFile) {
	messages := make([]string, 0)

	defer func() {
		ch <- messages
		wg.Done()
	}()

	values, ok := file.GetPropertyValues("tags")
	if !ok {
		return
	}
	fileUpdated := false

	for _, tag := range values {
		tagTemplateNote, ok := vault.GetTagTemplateNote(tag)
		if !ok {
			continue
		}

		for tagTemplateKey, tagTemplateValue := range tagTemplateNote.Frontmatter {
			if strings.HasPrefix(tagTemplateKey, "metadata.") {
				continue
			}

			_, ok = file.GetProperty(tagTemplateKey)
			if !ok {
				file.AddProperty(tagTemplateKey, tagTemplateValue.GetValues())
				fileUpdated = true
				message := fmt.Sprintf("A propriedade %s foi adicionada", tagTemplateKey)
				messages = append(messages, message)
			}
			propertyValue, _ := file.GetProperty(tagTemplateKey)

			//	required
			propertyMetaRequired := fmt.Sprintf("metadata.%s.required", tagTemplateKey)
			propertyRequired, ok := tagTemplateNote.GetProperty(propertyMetaRequired)
			if ok && propertyRequired.GetValues()[0] == "true" && len(propertyValue.GetValues()) == 0 {
				message := fmt.Sprintf("A propriedade %s é obrigatória", tagTemplateKey)
				messages = append(messages, message)
			}
		}
	}
	if fileUpdated {
		file.WriteFile()
	}

	if len(messages) > 0 {
		p, _ := filepath.Rel(vault.Path, file.Path)
		messages = append([]string{"======", "Nota: " + p}, messages...)
	}
}

func run(cmd *cobra.Command, args []string) {
	vaultDir, _ := cmd.Flags().GetString("vault-dir")
	templatesDir, _ := cmd.Flags().GetString("templates-dir")

	wdDir, _ := os.Getwd()
	vaultDirAbs, _ := utils.NormalizePath(wdDir, vaultDir)
	templatesDirAbs, _ := utils.NormalizePath(vaultDirAbs, templatesDir)

	vault := obsidian.NewVault(vaultDirAbs, templatesDirAbs)
	vault.LoadAllFiles()

	waitAnalyzeFiles := &sync.WaitGroup{}
	waitAnalyzeFiles.Add(len(vault.Notes))
	chMessages := make(chan []string)

	for _, file := range vault.Notes {
		go analyzeFile(chMessages, waitAnalyzeFiles, vault, file)
	}

	for range len(vault.Notes) {
		messages := <-chMessages
		for _, message := range messages {
			fmt.Println(message)
		}
	}
	close(chMessages)

	waitAnalyzeFiles.Wait()

}
