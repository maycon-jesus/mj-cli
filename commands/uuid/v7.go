package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// V7Command represents a Cobra command to generate and display a UUID v7.
var V7Command = &cobra.Command{
	Use:   "v7",
	Short: "Generate a new UUID v7",
	Run:   RunV7Command,
}

// GetV7Command returns the cobra.Command instance for generating a UUID v7.
func GetV7Command() *cobra.Command {
	return V7Command
}

// RunV7Command executes the given Cobra command, generates a UUID v7, and prints it to the standard output.
func RunV7Command(cmd *cobra.Command, args []string) {
	u, e := uuid.NewV7()
	if e != nil {
		panic(e)
	}
	fmt.Println(u)
}
