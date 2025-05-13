package snippets

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
)

// BulletJournalSnippetCommand is a Cobra command for managing bullet journal snippets with multiple subcommands execution.
var BulletJournalSnippetCommand = &cobra.Command{
	Use: "bullet-journal",
	Run: RunBulletJournalSnippetCommand,
}

// GetBulletJournalSnippetCommand returns the Cobra command for managing bullet journal snippets.
func GetBulletJournalSnippetCommand() *cobra.Command {
	return BulletJournalSnippetCommand
}

// MyPrinter is a struct that implements the io.Writer interface for handling content output with optional visibility toggle.
// It allows capturing written content in a buffer while optionally suppressing its display.
// The 'content' field stores all written data and 'hideContent' controls whether the data is printed to stdout.
type MyPrinter struct {
	content     string
	hideContent bool
}

// Write appends the provided byte slice to the content and conditionally prints it based on the hideContent flag.
// It returns the number of bytes written and any error encountered.
func (c MyPrinter) Write(p []byte) (int, error) {
	c.content += string(p)
	if !c.hideContent {
		fmt.Print(string(p))
	}
	return len(p), nil
}

// RunBulletJournalSnippetCommand executes a series of commands related to an Obsidian bullet journal workflow.
// It iterates over predefined command sets, runs them, and handles their output and errors.
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
