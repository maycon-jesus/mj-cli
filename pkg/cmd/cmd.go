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

func RunCommandTest(cmdStr string, writer io.Writer) error {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = writer
	cmd.Stderr = writer

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
