package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var V7Command = &cobra.Command{
	Use: "v7",
	Run: RunV7Command,
}

func GetV7Command() *cobra.Command {
	return V7Command
}

func RunV7Command(cmd *cobra.Command, args []string) {
	u, e := uuid.NewV7()
	if e != nil {
		panic(e)
	}
	fmt.Println(u)
}
