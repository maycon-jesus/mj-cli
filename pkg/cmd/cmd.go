package cmd

import (
	"os"
	"os/exec"
)

func RunCommand(cmdStr string) error {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
