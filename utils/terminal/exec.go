package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func clearStdin() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.Discard(reader.Buffered())
}

func ConvertStringToCmd(cmdStr string) ([]string, error) {
	var parts []string
	var current []rune
	insideQuotes := false
	escapeNext := false

	for _, c := range cmdStr {
		if escapeNext {
			fmt.Println(string(c))
			current = append(current, c)
			escapeNext = false
			continue
		}

		switch c {
		case '\\':
			escapeNext = true
		case '"':
			insideQuotes = !insideQuotes
			current = append(current, c)
		case ' ':
			if insideQuotes {
				current = append(current, c)
			} else if len(current) > 0 {
				parts = append(parts, string(current))
				current = []rune{}
			}
		default:
			current = append(current, c)
		}
	}

	if len(current) > 0 {
		parts = append(parts, string(current))
	}

	if insideQuotes {
		return nil, fmt.Errorf("string contém aspas não balanceadas")
	}

	return parts, nil
}

// RunCommandOptions defines configuration options for running a command, including input, output, and error stream settings.
// HideContent specifies whether to suppress content visibility or not when configuring output streams.
// Stdin represents the input stream source for the command execution.
// Stdout represents the output stream destination for standard output during execution.
// Stderr represents the output stream destination for standard error during execution.
type RunCommandOptions struct {
	Debug       bool
	HideContent bool
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
}

// RunCommandRealtime executes a command in real-time with specified arguments, using provided I/O configuration options.
// It supports redirecting input, output, and error streams based on the given RunCommandOptions settings.
// The function clears the standard input buffer before executing the command to ensure clean interactions.
// Returns an error if command execution fails.
func RunCommandRealtime(command string, opts RunCommandOptions) error {
	clearStdin()
	cmdArr, err := ConvertStringToCmd(command)
	if err != nil {
		return err
	}

	if opts.Debug {
		fmt.Println("[ARGS]")
		for i, v := range cmdArr {
			fmt.Printf("%d - %s\n", i, v)
		}
	}

	cmde := exec.Command(cmdArr[0], cmdArr[1:]...)
	cmde.Env = os.Environ()

	if !opts.HideContent {
		if opts.Debug {
			fmt.Println("[OUTPUT]")
		}
		if opts.Stdout != nil {
			cmde.Stdout = opts.Stdout
		} else {
			cmde.Stdout = os.Stdout
		}

		if opts.Stderr != nil {
			cmde.Stderr = opts.Stderr
		} else {
			cmde.Stderr = os.Stderr
		}
	}

	if opts.Stdin != nil {
		cmde.Stdin = opts.Stdin
	} else {
		cmde.Stdin = os.Stdin
	}

	err = cmde.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// RunCommand executes a shell command with the given arguments and returns its combined output and error, if any.
func RunCommand(command string, args []string) (string, error) {
	cmde := exec.Command(command, args...)
	content, err := cmde.CombinedOutput()
	return string(content), err
}
