package cmd

import (
	"io"
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

type CommandOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func RunCommandWithOptions(cmdStr string, options CommandOptions) error {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = options.Stdin
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func GetCommandOutput(cmdStr string) (string, error) {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
