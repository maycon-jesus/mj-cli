package tint

import (
	"os"
	"testing"
)

func TestCreateAnsiCodeSequence(t *testing.T) {
	tests := []struct {
		name  string
		codes []int
		want  string
	}{
		{
			name:  "sem códigos retorna reset",
			codes: []int{},
			want:  "\033[m",
		},
		{
			name:  "código único",
			codes: []int{1},
			want:  "\033[1m",
		},
		{
			name:  "dois códigos",
			codes: []int{1, 4},
			want:  "\033[1;4m",
		},
		{
			name:  "múltiplos códigos",
			codes: []int{1, 4, 31},
			want:  "\033[1;4;31m",
		},
		{
			name:  "código zero (reset explícito)",
			codes: []int{0},
			want:  "\033[0m",
		},
		{
			name:  "código de cor ANSI vermelho",
			codes: []int{31},
			want:  "\033[31m",
		},
		{
			name:  "negrito com cor",
			codes: []int{1, 31},
			want:  "\033[1;31m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createAnsiCodeSequence(tt.codes...)
			if got != tt.want {
				t.Errorf("createAnsiCodeSequence(%v) = %q, want %q", tt.codes, got, tt.want)
			}
		})
	}
}

// withEnvVars executa uma função com variáveis de ambiente temporárias
// e restaura o estado original após a execução.
func withEnvVars(t *testing.T, vars map[string]string, fn func()) {
	t.Helper()

	// Lista de variáveis relevantes para manipular
	relevantVars := []string{"NO_COLOR", "COLORTERM", "TERM"}

	// Salva valores originais
	original := make(map[string]*string)
	for _, key := range relevantVars {
		if val, exists := os.LookupEnv(key); exists {
			v := val
			original[key] = &v
		} else {
			original[key] = nil
		}
	}

	// Limpa todas as variáveis relevantes primeiro
	for _, key := range relevantVars {
		os.Unsetenv(key)
	}

	// Define as variáveis do teste
	for key, val := range vars {
		os.Setenv(key, val)
	}

	// Executa o teste
	defer func() {
		// Restaura variáveis originais
		for _, key := range relevantVars {
			os.Unsetenv(key)
		}
		for key, val := range original {
			if val != nil {
				os.Setenv(key, *val)
			}
		}
	}()

	fn()
}

func TestDetectColorSupport(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    ColorSupport
	}{
		// NO_COLOR tem prioridade máxima
		{
			name:    "NO_COLOR definido retorna ColorSupportNone",
			envVars: map[string]string{"NO_COLOR": "1"},
			want:    ColorSupportNone,
		},
		{
			name:    "NO_COLOR com valor qualquer retorna ColorSupportNone",
			envVars: map[string]string{"NO_COLOR": "true"},
			want:    ColorSupportNone,
		},

		// COLORTERM truecolor
		{
			name:    "COLORTERM=truecolor retorna ColorSupport24bit",
			envVars: map[string]string{"COLORTERM": "truecolor"},
			want:    ColorSupport24bit,
		},
		{
			name:    "COLORTERM=24bit retorna ColorSupport24bit",
			envVars: map[string]string{"COLORTERM": "24bit"},
			want:    ColorSupport24bit,
		},

		// TERM patterns - truecolor
		{
			name:    "TERM contendo truecolor retorna ColorSupport24bit",
			envVars: map[string]string{"TERM": "xterm-truecolor"},
			want:    ColorSupport24bit,
		},
		{
			name:    "TERM contendo 24bit retorna ColorSupport24bit",
			envVars: map[string]string{"TERM": "xterm-24bit"},
			want:    ColorSupport24bit,
		},

		// TERM patterns - 256 cores
		{
			name:    "TERM contendo 256color retorna ColorSupport8bit",
			envVars: map[string]string{"TERM": "xterm-256color"},
			want:    ColorSupport8bit,
		},

		// TERM patterns - terminais modernos (assume 8-bit)
		{
			name:    "TERM=xterm retorna ColorSupport8bit",
			envVars: map[string]string{"TERM": "xterm"},
			want:    ColorSupport8bit,
		},
		{
			name:    "TERM=screen retorna ColorSupport8bit",
			envVars: map[string]string{"TERM": "screen"},
			want:    ColorSupport8bit,
		},
		{
			name:    "TERM=tmux retorna ColorSupport8bit",
			envVars: map[string]string{"TERM": "tmux"},
			want:    ColorSupport8bit,
		},
		{
			name:    "TERM=ansi retorna ColorSupport8bit",
			envVars: map[string]string{"TERM": "ansi"},
			want:    ColorSupport8bit,
		},
		{
			name:    "TERM contendo color retorna ColorSupport8bit",
			envVars: map[string]string{"TERM": "linux-color"},
			want:    ColorSupport8bit,
		},

		// Fallback
		{
			name:    "sem variáveis relevantes retorna ColorSupport4bit",
			envVars: map[string]string{},
			want:    ColorSupport4bit,
		},
		{
			name:    "TERM desconhecido retorna ColorSupport4bit",
			envVars: map[string]string{"TERM": "dumb"},
			want:    ColorSupport4bit,
		},
		{
			name:    "TERM vazio retorna ColorSupport4bit",
			envVars: map[string]string{"TERM": ""},
			want:    ColorSupport4bit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withEnvVars(t, tt.envVars, func() {
				got := detectColorSupport()
				if got != tt.want {
					t.Errorf("detectColorSupport() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

func TestDetectColorSupport_NO_COLOR_Priority(t *testing.T) {
	// NO_COLOR deve ter prioridade sobre COLORTERM e TERM
	withEnvVars(t, map[string]string{
		"NO_COLOR":  "1",
		"COLORTERM": "truecolor",
		"TERM":      "xterm-256color",
	}, func() {
		got := detectColorSupport()
		if got != ColorSupportNone {
			t.Errorf("NO_COLOR deve ter prioridade, got %v, want ColorSupportNone", got)
		}
	})
}

func TestDetectColorSupport_COLORTERM_Priority(t *testing.T) {
	// COLORTERM deve ter prioridade sobre TERM
	withEnvVars(t, map[string]string{
		"COLORTERM": "truecolor",
		"TERM":      "xterm",
	}, func() {
		got := detectColorSupport()
		if got != ColorSupport24bit {
			t.Errorf("COLORTERM deve ter prioridade sobre TERM, got %v, want ColorSupport24bit", got)
		}
	})
}
