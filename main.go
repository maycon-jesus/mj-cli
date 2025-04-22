package main

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/commands"
	"os"
)

func main() {
	rootCmd := commands.GetCommandRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
