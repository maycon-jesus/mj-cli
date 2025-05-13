package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// V4Command represents a Cobra command for generating a new UUID v4.
var V4Command = &cobra.Command{
	Use:   "v4",
	Short: "Generate a new UUID v4",
	Run:   RunV4Command,
}

// GetV4Command returns the cobra.Command instance for generating a UUID v4.
func GetV4Command() *cobra.Command {
	return V4Command
}

// RunV4Command executes the given Cobra command, generates a new UUID v4, and prints it to the standard output.
func RunV4Command(cmd *cobra.Command, args []string) {
	u := uuid.New()
	fmt.Println(u)
}
