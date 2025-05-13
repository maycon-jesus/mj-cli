package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var V4Command = &cobra.Command{
	Use: "v4",
	Run: RunV4Command,
}

func GetV4Command() *cobra.Command {
	return V4Command
}

func RunV4Command(cmd *cobra.Command, args []string) {
	u := uuid.New()
	fmt.Println(u)
}
