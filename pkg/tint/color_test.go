package tint

import (
	"testing"
)

// withColorSupport executa uma função com um nível de suporte a cores específico
// e restaura o valor original após a execução.
func withColorSupport(t *testing.T, support ColorSupport, fn func()) {
	t.Helper()

	original := DetectedColorSupport
	DetectedColorSupport = support
	defer func() {
		DetectedColorSupport = original
	}()

	fn()
}

func TestTrueColor_GetForegroundCode(t *testing.T) {
	tests := []struct {
		name  string
		color TrueColor
		want  string
	}{
		{
			name:  "preto RGB(0,0,0)",
			color: TrueColor{R: 0, G: 0, B: 0},
			want:  "\033[38;2;0;0;0m",
		},
		{
			name:  "branco RGB(255,255,255)",
			color: TrueColor{R: 255, G: 255, B: 255},
			want:  "\033[38;2;255;255;255m",
		},
		{
			name:  "vermelho RGB(255,0,0)",
			color: TrueColor{R: 255, G: 0, B: 0},
			want:  "\033[38;2;255;0;0m",
		},
		{
			name:  "verde RGB(0,255,0)",
			color: TrueColor{R: 0, G: 255, B: 0},
			want:  "\033[38;2;0;255;0m",
		},
		{
			name:  "azul RGB(0,0,255)",
			color: TrueColor{R: 0, G: 0, B: 255},
			want:  "\033[38;2;0;0;255m",
		},
		{
			name:  "cor customizada RGB(128,64,32)",
			color: TrueColor{R: 128, G: 64, B: 32},
			want:  "\033[38;2;128;64;32m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.GetForegroundCode()
			if got != tt.want {
				t.Errorf("TrueColor.GetForegroundCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrueColor_GetBackgroundCode(t *testing.T) {
	tests := []struct {
		name  string
		color TrueColor
		want  string
	}{
		{
			name:  "preto RGB(0,0,0)",
			color: TrueColor{R: 0, G: 0, B: 0},
			want:  "\033[48;2;0;0;0m",
		},
		{
			name:  "branco RGB(255,255,255)",
			color: TrueColor{R: 255, G: 255, B: 255},
			want:  "\033[48;2;255;255;255m",
		},
		{
			name:  "vermelho RGB(255,0,0)",
			color: TrueColor{R: 255, G: 0, B: 0},
			want:  "\033[48;2;255;0;0m",
		},
		{
			name:  "cor customizada RGB(128,64,32)",
			color: TrueColor{R: 128, G: 64, B: 32},
			want:  "\033[48;2;128;64;32m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.GetBackgroundCode()
			if got != tt.want {
				t.Errorf("TrueColor.GetBackgroundCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColorAnsi256_GetForegroundCode(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{name: "código 0 (preto)", code: 0, want: "\033[38;5;0m"},
		{name: "código 15 (branco brilhante)", code: 15, want: "\033[38;5;15m"},
		{name: "código 196 (vermelho)", code: 196, want: "\033[38;5;196m"},
		{name: "código 208 (laranja)", code: 208, want: "\033[38;5;208m"},
		{name: "código 255", code: 255, want: "\033[38;5;255m"},
		{name: "código 128 (médio)", code: 128, want: "\033[38;5;128m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ColorAnsi256{Code: tt.code}
			got := c.GetForegroundCode()
			if got != tt.want {
				t.Errorf("ColorAnsi256{%d}.GetForegroundCode() = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestColorAnsi256_GetBackgroundCode(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{name: "código 0 (preto)", code: 0, want: "\033[48;5;0m"},
		{name: "código 15 (branco brilhante)", code: 15, want: "\033[48;5;15m"},
		{name: "código 196 (vermelho)", code: 196, want: "\033[48;5;196m"},
		{name: "código 255", code: 255, want: "\033[48;5;255m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ColorAnsi256{Code: tt.code}
			got := c.GetBackgroundCode()
			if got != tt.want {
				t.Errorf("ColorAnsi256{%d}.GetBackgroundCode() = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestColorAnsi_GetForegroundCode(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{name: "preto (30)", code: 30, want: "\033[30m"},
		{name: "vermelho (31)", code: 31, want: "\033[31m"},
		{name: "verde (32)", code: 32, want: "\033[32m"},
		{name: "amarelo (33)", code: 33, want: "\033[33m"},
		{name: "azul (34)", code: 34, want: "\033[34m"},
		{name: "magenta (35)", code: 35, want: "\033[35m"},
		{name: "ciano (36)", code: 36, want: "\033[36m"},
		{name: "branco (37)", code: 37, want: "\033[37m"},
		{name: "preto brilhante (90)", code: 90, want: "\033[90m"},
		{name: "vermelho brilhante (91)", code: 91, want: "\033[91m"},
		{name: "branco brilhante (97)", code: 97, want: "\033[97m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ColorAnsi{Code: tt.code}
			got := c.GetForegroundCode()
			if got != tt.want {
				t.Errorf("ColorAnsi{%d}.GetForegroundCode() = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestColorAnsi_GetBackgroundCode(t *testing.T) {
	tests := []struct {
		name string
		code int
		want string
	}{
		{name: "preto bg (30 -> 40)", code: 30, want: "\033[40m"},
		{name: "vermelho bg (31 -> 41)", code: 31, want: "\033[41m"},
		{name: "verde bg (32 -> 42)", code: 32, want: "\033[42m"},
		{name: "amarelo bg (33 -> 43)", code: 33, want: "\033[43m"},
		{name: "azul bg (34 -> 44)", code: 34, want: "\033[44m"},
		{name: "magenta bg (35 -> 45)", code: 35, want: "\033[45m"},
		{name: "ciano bg (36 -> 46)", code: 36, want: "\033[46m"},
		{name: "branco bg (37 -> 47)", code: 37, want: "\033[47m"},
		{name: "preto brilhante bg (90 -> 100)", code: 90, want: "\033[100m"},
		{name: "vermelho brilhante bg (91 -> 101)", code: 91, want: "\033[101m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ColorAnsi{Code: tt.code}
			got := c.GetBackgroundCode()
			if got != tt.want {
				t.Errorf("ColorAnsi{%d}.GetBackgroundCode() = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestCompleteColor_GetForegroundCode(t *testing.T) {
	testColor := CompleteColor{
		TrueColor: TrueColor{R: 255, G: 0, B: 0},
		Ansi256:   ColorAnsi256{Code: 196},
		Ansi:      ColorAnsi{Code: ANSI_COLOR_RED},
	}

	tests := []struct {
		name         string
		colorSupport ColorSupport
		want         string
	}{
		{
			name:         "ColorSupportNone retorna vazio",
			colorSupport: ColorSupportNone,
			want:         "",
		},
		{
			name:         "ColorSupport4bit usa Ansi",
			colorSupport: ColorSupport4bit,
			want:         "\033[31m",
		},
		{
			name:         "ColorSupport8bit usa Ansi256",
			colorSupport: ColorSupport8bit,
			want:         "\033[38;5;196m",
		},
		{
			name:         "ColorSupport24bit usa TrueColor",
			colorSupport: ColorSupport24bit,
			want:         "\033[38;2;255;0;0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, tt.colorSupport, func() {
				got := testColor.GetForegroundCode()
				if got != tt.want {
					t.Errorf("CompleteColor.GetForegroundCode() = %q, want %q", got, tt.want)
				}
			})
		})
	}
}

func TestCompleteColor_GetBackgroundCode(t *testing.T) {
	testColor := CompleteColor{
		TrueColor: TrueColor{R: 0, G: 0, B: 255},
		Ansi256:   ColorAnsi256{Code: 21},
		Ansi:      ColorAnsi{Code: ANSI_COLOR_BLUE},
	}

	tests := []struct {
		name         string
		colorSupport ColorSupport
		want         string
	}{
		{
			name:         "ColorSupportNone retorna vazio",
			colorSupport: ColorSupportNone,
			want:         "",
		},
		{
			name:         "ColorSupport4bit usa Ansi",
			colorSupport: ColorSupport4bit,
			want:         "\033[44m", // 34 + 10
		},
		{
			name:         "ColorSupport8bit usa Ansi256",
			colorSupport: ColorSupport8bit,
			want:         "\033[48;5;21m",
		},
		{
			name:         "ColorSupport24bit usa TrueColor",
			colorSupport: ColorSupport24bit,
			want:         "\033[48;2;0;0;255m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withColorSupport(t, tt.colorSupport, func() {
				got := testColor.GetBackgroundCode()
				if got != tt.want {
					t.Errorf("CompleteColor.GetBackgroundCode() = %q, want %q", got, tt.want)
				}
			})
		})
	}
}

func TestCompleteColor_DefaultCase(t *testing.T) {
	// Testa o case default do switch (usa TrueColor)
	testColor := CompleteColor{
		TrueColor: TrueColor{R: 128, G: 64, B: 32},
		Ansi256:   ColorAnsi256{Code: 136},
		Ansi:      ColorAnsi{Code: ANSI_COLOR_YELLOW},
	}

	// Usa um valor de ColorSupport inválido para testar o default
	withColorSupport(t, ColorSupport(999), func() {
		gotFg := testColor.GetForegroundCode()
		wantFg := "\033[38;2;128;64;32m"
		if gotFg != wantFg {
			t.Errorf("CompleteColor.GetForegroundCode() default = %q, want %q", gotFg, wantFg)
		}

		gotBg := testColor.GetBackgroundCode()
		wantBg := "\033[48;2;128;64;32m"
		if gotBg != wantBg {
			t.Errorf("CompleteColor.GetBackgroundCode() default = %q, want %q", gotBg, wantBg)
		}
	})
}

func TestPredefinedColors(t *testing.T) {
	tests := []struct {
		name          string
		color         ColorAnsi
		expectedCode  int
		expectedFgStr string
	}{
		{"Black", Black, ANSI_COLOR_BLACK, "\033[30m"},
		{"Red", Red, ANSI_COLOR_RED, "\033[31m"},
		{"Green", Green, ANSI_COLOR_GREEN, "\033[32m"},
		{"Yellow", Yellow, ANSI_COLOR_YELLOW, "\033[33m"},
		{"Blue", Blue, ANSI_COLOR_BLUE, "\033[34m"},
		{"Magenta", Magenta, ANSI_COLOR_MAGENTA, "\033[35m"},
		{"Cyan", Cyan, ANSI_COLOR_CYAN, "\033[36m"},
		{"White", White, ANSI_COLOR_WHITE, "\033[37m"},
		{"BrightBlack", BrightBlack, ANSI_COLOR_BLACK + ANSI_COLOR_BRIGHT_OFFSET, "\033[90m"},
		{"BrightRed", BrightRed, ANSI_COLOR_RED + ANSI_COLOR_BRIGHT_OFFSET, "\033[91m"},
		{"BrightGreen", BrightGreen, ANSI_COLOR_GREEN + ANSI_COLOR_BRIGHT_OFFSET, "\033[92m"},
		{"BrightYellow", BrightYellow, ANSI_COLOR_YELLOW + ANSI_COLOR_BRIGHT_OFFSET, "\033[93m"},
		{"BrightBlue", BrightBlue, ANSI_COLOR_BLUE + ANSI_COLOR_BRIGHT_OFFSET, "\033[94m"},
		{"BrightMagenta", BrightMagenta, ANSI_COLOR_MAGENTA + ANSI_COLOR_BRIGHT_OFFSET, "\033[95m"},
		{"BrightCyan", BrightCyan, ANSI_COLOR_CYAN + ANSI_COLOR_BRIGHT_OFFSET, "\033[96m"},
		{"BrightWhite", BrightWhite, ANSI_COLOR_WHITE + ANSI_COLOR_BRIGHT_OFFSET, "\033[97m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color.Code != tt.expectedCode {
				t.Errorf("%s.Code = %d, want %d", tt.name, tt.color.Code, tt.expectedCode)
			}

			gotFg := tt.color.GetForegroundCode()
			if gotFg != tt.expectedFgStr {
				t.Errorf("%s.GetForegroundCode() = %q, want %q", tt.name, gotFg, tt.expectedFgStr)
			}
		})
	}
}

func TestColorInterface(t *testing.T) {
	// Verifica que todas as implementações satisfazem a interface Color
	var _ Color = TrueColor{}
	var _ Color = ColorAnsi256{}
	var _ Color = ColorAnsi{}
	var _ Color = CompleteColor{}
}
