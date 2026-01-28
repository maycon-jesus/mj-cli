package cmd

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:  "simple command",
			input: "ls -la",
			want:  []string{"ls", "-la"},
		},
		{
			name:  "command with multiple spaces",
			input: "echo   hello    world",
			want:  []string{"echo", "hello", "world"},
		},
		{
			name:  "command with double quotes",
			input: `echo "hello world"`,
			want:  []string{"echo", "hello world"},
		},
		{
			name:  "command with single quotes",
			input: `echo 'hello world'`,
			want:  []string{"echo", "hello world"},
		},
		{
			name:  "command with escaped double quote inside double quotes",
			input: `echo "hello \"world\""`,
			want:  []string{"echo", `hello "world"`},
		},
		{
			name:  "command with escaped backslash inside double quotes",
			input: `echo "hello\\world"`,
			want:  []string{"echo", `hello\world`},
		},
		{
			name:  "command with backslash in single quotes (literal)",
			input: `echo 'hello\world'`,
			want:  []string{"echo", `hello\world`},
		},
		{
			name:  "command with invalid escape in double quotes",
			input: `echo "hello\nworld"`,
			want:  []string{"echo", `hello\nworld`},
		},
		{
			name:  "command with escaped space outside quotes",
			input: `echo hello\ world`,
			want:  []string{"echo", "hello world"},
		},
		{
			name:  "command with mixed quotes",
			input: `echo "hello" 'world'`,
			want:  []string{"echo", "hello", "world"},
		},
		{
			name:  "command with single quote inside double quotes",
			input: `echo "hello 'world'"`,
			want:  []string{"echo", "hello 'world'"},
		},
		{
			name:  "command with double quote inside single quotes",
			input: `echo 'hello "world"'`,
			want:  []string{"echo", `hello "world"`},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "only spaces",
			input: "   ",
			want:  nil,
		},
		{
			name:  "command with escaped dollar sign",
			input: `echo "\$HOME"`,
			want:  []string{"echo", "$HOME"},
		},
		{
			name:  "command with escaped backtick",
			input: "echo \"\\`date\\`\"",
			want:  []string{"echo", "`date`"},
		},
		{
			name:  "complex git command",
			input: `git commit -m "feat: add new feature"`,
			want:  []string{"git", "commit", "-m", "feat: add new feature"},
		},
		{
			name:    "unclosed double quote",
			input:   `echo "hello world`,
			wantErr: true,
			errMsg:  "unclosed quote: \"",
		},
		{
			name:    "unclosed single quote",
			input:   `echo 'hello world`,
			wantErr: true,
			errMsg:  "unclosed quote: '",
		},
		{
			name:    "trailing backslash",
			input:   `echo hello\`,
			wantErr: true,
			errMsg:  "unexpected end of command: trailing backslash",
		},
		{
			name:  "multiple arguments with various quote types",
			input: `cmd arg1 "arg 2" 'arg 3' arg4`,
			want:  []string{"cmd", "arg1", "arg 2", "arg 3", "arg4"},
		},
		{
			name:  "leading and trailing spaces",
			input: "  echo hello  ",
			want:  []string{"echo", "hello"},
		},
		{
			name:  "command with tabs",
			input: "echo\thello\tworld",
			want:  []string{"echo", "hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCommand(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseCommand() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("parseCommand() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("parseCommand() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParseCommand(b *testing.B) {
	testCases := []string{
		"ls -la",
		`echo "hello world"`,
		`git commit -m "feat: add new feature"`,
		`cmd arg1 "arg 2" 'arg 3' arg4`,
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = parseCommand(tc)
			}
		})
	}
}
