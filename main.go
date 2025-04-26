package main

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/commands"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/spf13/viper"
	"os"
)

func main() {
	utils.LoadViper()

	//vault := obsidian.NewVault("/mnt/c/Users/conta/OneDrive/Documentos/notes")
	//vault.LoadAllFiles()
	//
	//note := vault.GetNote("00 - MOCs/99 - Meta/02 - Tags/book.md")
	//note.AddProperty("teste-escrita", []string{"aaa"}, obsidian.FilePropertyMetadata{})
	//note.WriteFile()
	rootCmd := commands.GetCommandRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Try save on exists config file
	err := viper.WriteConfig()
	if err != nil {
		//If config file not exists try create new
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := viper.SafeWriteConfig()
			if err != nil {
				panic(err)

			}
		} else {
			panic(err)
		}
	}
}
