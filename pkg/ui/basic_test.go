package ui

import (
	"strings"
	"testing"

	"github.com/maycon-jesus/mj-cli/pkg/tint"
)

func withColorSupport(t *testing.T, support tint.ColorSupport, fn func()) {
	t.Helper()

	original := tint.DetectedColorSupport
	tint.DetectedColorSupport = support
	defer func() {
		tint.DetectedColorSupport = original
	}()

	fn()
}

func TestPresets_SemCores(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		msg  string
		want string
	}{
		{"Title", Title, "título", "=== título ==="},
		{"Subtitle", Subtitle, "subtítulo", "--- subtítulo ---"},
		{"Success", Success, "sucesso", "✔ sucesso"},
		{"Error", Error, "erro", "✖ erro"},
		{"Warning", Warning, "aviso", "⚠ aviso"},
		{"Info", Info, "info", "ℹ info"},
		{"Lambda", Lambda, "lambda", "λ lambda"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, tint.ColorSupportNone, func() {
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
		text  string
		color tint.ColorAnsi
	}{
		{"Title", Title, "título", "=== título ===", tint.Cyan},
		{"Subtitle", Subtitle, "subtítulo", "--- subtítulo ---", tint.White},
		{"Success", Success, "sucesso", "✔ sucesso", tint.Green},
		{"Error", Error, "erro", "✖ erro", tint.Red},
		{"Warning", Warning, "aviso", "⚠ aviso", tint.Yellow},
		{"Info", Info, "info", "ℹ info", tint.Blue},
		{"Lambda", Lambda, "lambda", "λ lambda", tint.Magenta},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, tint.ColorSupport4bit, func() {
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
		{"Title", Title, "\033[0m\033[1m\033[36m=== x ===\033[0m"},
		{"Subtitle", Subtitle, "\033[0m\033[1m\033[37m--- x ---\033[0m"},
		{"Success", Success, "\033[0m\033[1m\033[32m✔ x\033[0m"},
		{"Error", Error, "\033[0m\033[1m\033[31m✖ x\033[0m"},
		{"Warning", Warning, "\033[0m\033[1m\033[33m⚠ x\033[0m"},
		{"Info", Info, "\033[0m\033[1m\033[34mℹ x\033[0m"},
		{"Lambda", Lambda, "\033[0m\033[1m\033[35mλ x\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, tint.ColorSupport4bit, func() {
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
		{"Title", Title, "===  ==="},
		{"Subtitle", Subtitle, "---  ---"},
		{"Success", Success, "✔ "},
		{"Error", Error, "✖ "},
		{"Warning", Warning, "⚠ "},
		{"Info", Info, "ℹ "},
		{"Lambda", Lambda, "λ "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, tint.ColorSupportNone, func() {
				got := tt.fn("")
				if got != tt.want {
					t.Errorf("%s(\"\") = %q, want %q", tt.name, got, tt.want)
				}
			})
		})
	}
}

func TestPresets_CaracteresEspeciais(t *testing.T) {
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
		{"Title", Title},
		{"Subtitle", Subtitle},
		{"Success", Success},
		{"Error", Error},
		{"Warning", Warning},
		{"Info", Info},
		{"Lambda", Lambda},
	}

	for _, msg := range messages {
		for _, preset := range presets {
			t.Run(preset.name+"/"+msg.name, func(t *testing.T) {
				withColorSupport(t, tint.ColorSupportNone, func() {
					got := preset.fn(msg.text)
					if !strings.Contains(got, msg.text) {
						t.Errorf("%s(%q) não preservou a mensagem, got %q", preset.name, msg.text, got)
					}
				})
			})
		}
	}
}
