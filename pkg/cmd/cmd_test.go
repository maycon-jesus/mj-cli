package cmd

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestRunCommandWithOptions(t *testing.T) {
	tests := []struct {
		name        string
		cmdStr      string
		input       string
		wantOutput  string
		wantErr     bool
		skipOnOS    string
	}{
		{
			name:       "echo command",
			cmdStr:     "echo hello",
			wantOutput: "hello\n",
		},
		{
			name:       "echo with multiple arguments",
			cmdStr:     "echo hello world",
			wantOutput: "hello world\n",
		},
		{
			name:       "echo with quoted string",
			cmdStr:     `echo "hello world"`,
			wantOutput: "hello world\n",
		},
		{
			name:    "non-existent command",
			cmdStr:  "nonexistentcommand123",
			wantErr: true,
		},
		{
			name:     "printf command",
			cmdStr:   "printf 'test'",
			skipOnOS: "windows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnOS != "" && runtime.GOOS == tt.skipOnOS {
				t.Skipf("Skipping test on %s", runtime.GOOS)
			}

			var stdout, stderr bytes.Buffer
			var stdin strings.Reader
			if tt.input != "" {
				stdin = *strings.NewReader(tt.input)
			}

			options := CommandOptions{
				Stdin:  &stdin,
				Stdout: &stdout,
				Stderr: &stderr,
			}

			err := RunCommandWithOptions(tt.cmdStr, options)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RunCommandWithOptions() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("RunCommandWithOptions() unexpected error = %v, stderr = %s", err, stderr.String())
				return
			}

			if tt.wantOutput != "" {
				got := stdout.String()
				if got != tt.wantOutput {
					t.Errorf("RunCommandWithOptions() output = %q, want %q", got, tt.wantOutput)
				}
			}
		})
	}
}

func TestGetCommandOutput(t *testing.T) {
	tests := []struct {
		name       string
		cmdStr     string
		wantOutput string
		contains   string
		wantErr    bool
		skipOnOS   string
	}{
		{
			name:       "echo command",
			cmdStr:     "echo hello",
			wantOutput: "hello\n",
		},
		{
			name:       "echo with multiple arguments",
			cmdStr:     "echo hello world",
			wantOutput: "hello world\n",
		},
		{
			name:    "non-existent command",
			cmdStr:  "nonexistentcommand123",
			wantErr: true,
		},
		{
			name:     "printf command",
			cmdStr:   "printf 'test'",
			contains: "test",
			skipOnOS: "windows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnOS != "" && runtime.GOOS == tt.skipOnOS {
				t.Skipf("Skipping test on %s", runtime.GOOS)
			}

			got, err := GetCommandOutput(tt.cmdStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCommandOutput() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetCommandOutput() unexpected error = %v", err)
				return
			}

			if tt.wantOutput != "" && got != tt.wantOutput {
				t.Errorf("GetCommandOutput() = %q, want %q", got, tt.wantOutput)
			}

			if tt.contains != "" && !strings.Contains(got, tt.contains) {
				t.Errorf("GetCommandOutput() = %q, should contain %q", got, tt.contains)
			}
		})
	}
}

func TestGetCommandOutput_CapturesStderr(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows")
	}

	// Create a command that writes to stderr
	cmdStr := "sh -c 'echo stdout; echo stderr >&2'"
	output, err := GetCommandOutput(cmdStr)

	if err != nil {
		t.Errorf("GetCommandOutput() unexpected error = %v", err)
		return
	}

	if !strings.Contains(output, "stdout") {
		t.Errorf("GetCommandOutput() should contain stdout, got %q", output)
	}

	if !strings.Contains(output, "stderr") {
		t.Errorf("GetCommandOutput() should contain stderr, got %q", output)
	}
}

func TestRunCommandWithOptions_StdinRedirect(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows")
	}

	input := "test input\n"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	options := CommandOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stdout,
	}

	err := RunCommandWithOptions("cat", options)
	if err != nil {
		t.Errorf("RunCommandWithOptions() unexpected error = %v", err)
		return
	}

	if stdout.String() != input {
		t.Errorf("RunCommandWithOptions() output = %q, want %q", stdout.String(), input)
	}
}

func TestCommandOptions_NilStreams(t *testing.T) {
	var stdout bytes.Buffer

	options := CommandOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: nil,
	}

	err := RunCommandWithOptions("echo test", options)
	if err != nil {
		t.Errorf("RunCommandWithOptions() with nil streams unexpected error = %v", err)
	}
}

func TestParseCommandErrors(t *testing.T) {
	tests := []struct {
		name   string
		cmdStr string
	}{
		{
			name:   "unclosed quote in RunCommandWithOptions",
			cmdStr: `echo "unclosed`,
		},
		{
			name:   "trailing backslash in RunCommandWithOptions",
			cmdStr: `echo test\`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			options := CommandOptions{
				Stdout: &stdout,
				Stderr: &stdout,
			}

			err := RunCommandWithOptions(tt.cmdStr, options)
			if err == nil {
				t.Errorf("RunCommandWithOptions() expected error for %q but got none", tt.cmdStr)
			}
		})

		t.Run(tt.name+" GetCommandOutput", func(t *testing.T) {
			_, err := GetCommandOutput(tt.cmdStr)
			if err == nil {
				t.Errorf("GetCommandOutput() expected error for %q but got none", tt.cmdStr)
			}
		})
	}
}

func TestRunCommandWithOptions_CommandNotFound(t *testing.T) {
	var stdout bytes.Buffer
	options := CommandOptions{
		Stdout: &stdout,
		Stderr: &stdout,
	}

	err := RunCommandWithOptions("commandthatdoesnotexist123456", options)
	if err == nil {
		t.Error("RunCommandWithOptions() should return error for non-existent command")
		return
	}

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) && !errors.Is(err, exec.ErrNotFound) {
		t.Logf("Expected exec.ExitError or exec.ErrNotFound, got: %T", err)
	}
}

func TestCommandWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name       string
		cmdStr     string
		wantOutput string
		skipOnOS   string
	}{
		{
			name:       "command with dollar sign",
			cmdStr:     `echo '$HOME'`,
			wantOutput: "$HOME\n",
		},
		{
			name:       "command with asterisk in quotes",
			cmdStr:     `echo '*'`,
			wantOutput: "*\n",
		},
		{
			name:       "command with semicolon in quotes",
			cmdStr:     `echo 'test;test'`,
			wantOutput: "test;test\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnOS != "" && runtime.GOOS == tt.skipOnOS {
				t.Skipf("Skipping test on %s", runtime.GOOS)
			}

			output, err := GetCommandOutput(tt.cmdStr)
			if err != nil {
				t.Errorf("GetCommandOutput() unexpected error = %v", err)
				return
			}

			if output != tt.wantOutput {
				t.Errorf("GetCommandOutput() = %q, want %q", output, tt.wantOutput)
			}
		})
	}
}

func BenchmarkRunCommandWithOptions(b *testing.B) {
	var stdout bytes.Buffer
	options := CommandOptions{
		Stdout: &stdout,
		Stderr: &stdout,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdout.Reset()
		_ = RunCommandWithOptions("echo test", options)
	}
}

func BenchmarkGetCommandOutput(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetCommandOutput("echo test")
	}
}
