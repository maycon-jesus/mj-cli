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
	"time"
)

var TagsProperties = &cobra.Command{
	Use:     "format-tags",
	Short:   "Format all file bases on tags templates",
	Aliases: []string{"tags-properties"},
	Run:     run,
}

func GetCommandTagsProperties() *cobra.Command {
	TagsProperties.Flags().BoolP("watch", "w", false, "Watch for changes")
	return TagsProperties
}

func run(cmd *cobra.Command, args []string) {
	watch, _ := cmd.Flags().GetBool("watch")
	vaultDir, _ := cmd.Flags().GetString("vault-dir")
	wdDir, _ := os.Getwd()
	vaultDirAbs, _ := utils.NormalizePath(wdDir, vaultDir)

	vault := obsidian.NewVault(vaultDirAbs)
	for run := true; run != false; {

		filesModified := vault.LoadAllFiles()

		waitAnalyzeFiles := &sync.WaitGroup{}
		waitAnalyzeFiles.Add(len(filesModified))
		chMessages := make(chan []string)

		for _, file := range filesModified {
			go analyzeFile(chMessages, waitAnalyzeFiles, vault, file)
		}

		for range len(filesModified) {
			messages := <-chMessages
			for _, message := range messages {
				fmt.Println(message)
			}
		}
		close(chMessages)
		waitAnalyzeFiles.Wait()

		if watch {
			fmt.Println("watch")
			time.Sleep(time.Second * 5)
		} else {
			run = false
		}
	}
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
