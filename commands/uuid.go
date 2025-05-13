package commands

import (
	"github.com/maycon-jesus/mj-cli/commands/uuid"
	"github.com/spf13/cobra"
)

var UuidCommand = &cobra.Command{
	Use:   "uuid",
	Short: "UUID tools",
}

func GetUuidCommand() *cobra.Command {
	UuidCommand.AddCommand(uuid.GetV4Command())
	UuidCommand.AddCommand(uuid.GetV7Command())

	return UuidCommand
}
