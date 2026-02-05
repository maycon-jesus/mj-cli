package tint

import (
	"strings"
	"testing"
)

func TestPresets_SemCores(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		msg  string
		want string
	}{
		{"PresetTitle", PresetTitle, "título", "=== título ==="},
		{"PresetSubtitle", PresetSubtitle, "subtítulo", "--- subtítulo ---"},
		{"PresetSuccess", PresetSuccess, "sucesso", "✔ sucesso"},
		{"PresetError", PresetError, "erro", "✖ erro"},
		{"PresetWarning", PresetWarning, "aviso", "⚠ aviso"},
		{"PresetInfo", PresetInfo, "info", "ℹ info"},
		{"PresetLambda", PresetLambda, "lambda", "λ lambda"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, ColorSupportNone, func() {
				got := tt.fn(tt.msg)
				if got != tt.want {
					t.Errorf("%s(%q) = %q, want %q", tt.name, tt.msg, got, tt.want)
				}
			})
		})
	}
}

func TestPresets_ComCores(t *testing.T) {
	tests := []struct {
		name  string
		fn    func(string) string
		msg   string
		text  string    // texto formatado esperado dentro da saída
		color ColorAnsi // cor de foreground esperada
	}{
		{"PresetTitle", PresetTitle, "título", "=== título ===", Cyan},
		{"PresetSubtitle", PresetSubtitle, "subtítulo", "--- subtítulo ---", White},
		{"PresetSuccess", PresetSuccess, "sucesso", "✔ sucesso", Green},
		{"PresetError", PresetError, "erro", "✖ erro", Red},
		{"PresetWarning", PresetWarning, "aviso", "⚠ aviso", Yellow},
		{"PresetInfo", PresetInfo, "info", "ℹ info", Blue},
		{"PresetLambda", PresetLambda, "lambda", "λ lambda", Magenta},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, ColorSupport4bit, func() {
				got := tt.fn(tt.msg)

				if !strings.Contains(got, tt.text) {
					t.Errorf("%s(%q) deveria conter %q, got %q", tt.name, tt.msg, tt.text, got)
				}
				if !strings.Contains(got, "\033[1m") {
					t.Errorf("%s(%q) deveria conter código bold, got %q", tt.name, tt.msg, got)
				}
				if !strings.Contains(got, tt.color.GetForegroundCode()) {
					t.Errorf("%s(%q) deveria conter código de cor %q, got %q", tt.name, tt.msg, tt.color.GetForegroundCode(), got)
				}
				if !strings.HasPrefix(got, "\033[0m") {
					t.Errorf("%s(%q) deveria começar com reset, got %q", tt.name, tt.msg, got)
				}
				if !strings.HasSuffix(got, "\033[0m") {
					t.Errorf("%s(%q) deveria terminar com reset, got %q", tt.name, tt.msg, got)
				}
			})
		})
	}
}

func TestPresets_SaidaExata(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{"PresetTitle", PresetTitle, "\033[0m\033[1m\033[36m=== x ===\033[0m"},
		{"PresetSubtitle", PresetSubtitle, "\033[0m\033[1m\033[37m--- x ---\033[0m"},
		{"PresetSuccess", PresetSuccess, "\033[0m\033[1m\033[32m✔ x\033[0m"},
		{"PresetError", PresetError, "\033[0m\033[1m\033[31m✖ x\033[0m"},
		{"PresetWarning", PresetWarning, "\033[0m\033[1m\033[33m⚠ x\033[0m"},
		{"PresetInfo", PresetInfo, "\033[0m\033[1m\033[34mℹ x\033[0m"},
		{"PresetLambda", PresetLambda, "\033[0m\033[1m\033[35mλ x\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, ColorSupport4bit, func() {
				got := tt.fn("x")
				if got != tt.want {
					t.Errorf("%s(\"x\") = %q, want %q", tt.name, got, tt.want)
				}
			})
		})
	}
}

func TestPresets_MensagemVazia(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{"PresetTitle", PresetTitle, "===  ==="},
		{"PresetSubtitle", PresetSubtitle, "---  ---"},
		{"PresetSuccess", PresetSuccess, "✔ "},
		{"PresetError", PresetError, "✖ "},
		{"PresetWarning", PresetWarning, "⚠ "},
		{"PresetInfo", PresetInfo, "ℹ "},
		{"PresetLambda", PresetLambda, "λ "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, ColorSupportNone, func() {
				got := tt.fn("")
				if got != tt.want {
					t.Errorf("%s(\"\") = %q, want %q", tt.name, got, tt.want)
				}
			})
		})
	}
}

func TestPresets_CaractereEspeciais(t *testing.T) {
	messages := []struct {
		name string
		text string
	}{
		{"espaços", "texto com espaços"},
		{"newlines", "texto\ncom\nnewlines"},
		{"tabs", "texto\tcom\ttabs"},
		{"unicode", "日本語テスト"},
		{"emojis", "🎉🚀✨"},
	}

	presets := []struct {
		name string
		fn   func(string) string
	}{
		{"PresetTitle", PresetTitle},
		{"PresetSubtitle", PresetSubtitle},
		{"PresetSuccess", PresetSuccess},
		{"PresetError", PresetError},
		{"PresetWarning", PresetWarning},
		{"PresetInfo", PresetInfo},
		{"PresetLambda", PresetLambda},
	}

	for _, msg := range messages {
		for _, preset := range presets {
			t.Run(preset.name+"/"+msg.name, func(t *testing.T) {
				withColorSupport(t, ColorSupportNone, func() {
					got := preset.fn(msg.text)
					if !strings.Contains(got, msg.text) {
						t.Errorf("%s(%q) não preservou a mensagem, got %q", preset.name, msg.text, got)
					}
				})
			})
		}
	}
}
