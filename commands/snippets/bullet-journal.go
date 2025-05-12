package snippets

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
)

var BulletJournalSnippetCommand = &cobra.Command{
	Use: "bullet-journal",
	Run: RunBulletJournalSnippetCommand,
}

func GetBulletJournalSnippetCommand() *cobra.Command {
	return BulletJournalSnippetCommand
}

type MyPrinter struct {
	content     string
	hideContent bool
}

func (c MyPrinter) Write(p []byte) (int, error) {
	c.content += string(p)
	if !c.hideContent {
		fmt.Print(string(p))
	}
	return len(p), nil
}

func RunBulletJournalSnippetCommand(cmd *cobra.Command, args []string) {

	programPath := os.Args[0]
	commands := [][]string{
		[]string{"obsidian", "month", "--soft"},
		[]string{"obsidian", "week", "--soft"},
		[]string{"obsidian", "daily", "--soft"},
	}

	for i, command := range commands {
		fmt.Println("====[" + strconv.Itoa(i) + "]====")
		cmde := exec.Command(programPath, command...)
		cmde.Stdout = MyPrinter{}
		cmde.Stderr = MyPrinter{}
		err := cmde.Run()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
	}
}
