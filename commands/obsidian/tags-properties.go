package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"github.com/maycon-jesus/mj-cli/utils/obsidian/tagRuler"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
)

var TagsProperties = &cobra.Command{
	Use:     "format-tags",
	Short:   "Format all file bases on tags templates",
	Aliases: []string{"tags-properties"},
	Run:     run,
}

func GetCommandTagsProperties() *cobra.Command {
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

	for _, tag := range values {
		tagRule, ok := tagRuler.TagsRules[tag]
		if !ok {
			continue
		}
		messages = tagRule.ApplyRules(file)
		if len(messages) > 0 {
			file.WriteFile()
		}
	}

	if len(messages) > 0 {
		p, _ := filepath.Rel(vault.Path, file.Path)
		messages = append([]string{"======", "Nota: " + p}, messages...)
	}
}

func run(cmd *cobra.Command, args []string) {
	vaultDir, _ := cmd.Flags().GetString("vault-dir")

	wdDir, _ := os.Getwd()
	vaultDirAbs, _ := utils.NormalizePath(wdDir, vaultDir)

	vault := obsidian.NewVault(vaultDirAbs)
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
