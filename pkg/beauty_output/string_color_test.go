package beautyoutput

import (
	"testing"
)

func TestNewStrBuilder(t *testing.T) {
	sb := NewStrBuilder()
	if sb == nil {
		t.Error("NewStrBuilder should not return nil")
	}
	if sb.content != "" {
		t.Errorf("NewStrBuilder content should be empty, got %q", sb.content)
	}
}

func TestTextStyles(t *testing.T) {
	tests := []struct {
		name     string
		apply    func(*StrBuilder) *StrBuilder
		expected string
	}{
		{"Bold", (*StrBuilder).Bold, "\033[1m"},
		{"Italic", (*StrBuilder).Italic, "\033[3m"},
		{"Underline", (*StrBuilder).Underline, "\033[4m"},
		{"Dim", (*StrBuilder).Dim, "\033[2m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			tt.apply(sb)
			if sb.content != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, sb.content, tt.expected)
			}
		})
	}
}

func TestStandardColors(t *testing.T) {
	tests := []struct {
		name     string
		apply    func(*StrBuilder) *StrBuilder
		expected string
	}{
		{"Black", (*StrBuilder).Black, "\033[30m"},
		{"Red", (*StrBuilder).Red, "\033[31m"},
		{"Green", (*StrBuilder).Green, "\033[32m"},
		{"Yellow", (*StrBuilder).Yellow, "\033[33m"},
		{"Blue", (*StrBuilder).Blue, "\033[34m"},
		{"Magenta", (*StrBuilder).Magenta, "\033[35m"},
		{"Cyan", (*StrBuilder).Cyan, "\033[36m"},
		{"White", (*StrBuilder).White, "\033[37m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			tt.apply(sb)
			if sb.content != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, sb.content, tt.expected)
			}
		})
	}
}

func TestBrightColors(t *testing.T) {
	tests := []struct {
		name     string
		apply    func(*StrBuilder) *StrBuilder
		expected string
	}{
		{"BrightBlack", (*StrBuilder).BrightBlack, "\033[90m"},
		{"BrightRed", (*StrBuilder).BrightRed, "\033[91m"},
		{"BrightGreen", (*StrBuilder).BrightGreen, "\033[92m"},
		{"BrightYellow", (*StrBuilder).BrightYellow, "\033[93m"},
		{"BrightBlue", (*StrBuilder).BrightBlue, "\033[94m"},
		{"BrightMagenta", (*StrBuilder).BrightMagenta, "\033[95m"},
		{"BrightCyan", (*StrBuilder).BrightCyan, "\033[96m"},
		{"BrightWhite", (*StrBuilder).BrightWhite, "\033[97m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			tt.apply(sb)
			if sb.content != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, sb.content, tt.expected)
			}
		})
	}
}

func TestRGB(t *testing.T) {
	tests := []struct {
		name     string
		r, g, b  int
		expected string
	}{
		{"Black RGB", 0, 0, 0, "\033[38;2;0;0;0m"},
		{"White RGB", 255, 255, 255, "\033[38;2;255;255;255m"},
		{"Red RGB", 255, 0, 0, "\033[38;2;255;0;0m"},
		{"Green RGB", 0, 255, 0, "\033[38;2;0;255;0m"},
		{"Blue RGB", 0, 0, 255, "\033[38;2;0;0;255m"},
		{"Custom RGB", 128, 64, 32, "\033[38;2;128;64;32m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			sb.RGB(tt.r, tt.g, tt.b)
			if sb.content != tt.expected {
				t.Errorf("RGB(%d, %d, %d) = %q, want %q", tt.r, tt.g, tt.b, sb.content, tt.expected)
			}
		})
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{"Simple text", "Hello", "Hello"},
		{"Empty text", "", ""},
		{"Text with spaces", "Hello World", "Hello World"},
		{"Text with special chars", "Hello\nWorld", "Hello\nWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			sb.Text(tt.text)
			if sb.content != tt.expected {
				t.Errorf("Text(%q) = %q, want %q", tt.text, sb.content, tt.expected)
			}
		})
	}
}

func TestReset(t *testing.T) {
	sb := NewStrBuilder()
	sb.Reset()
	expected := "\033[0m"
	if sb.content != expected {
		t.Errorf("Reset() = %q, want %q", sb.content, expected)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*StrBuilder)
		expected string
	}{
		{
			name:     "Empty builder",
			setup:    func(sb *StrBuilder) {},
			expected: "\033[0m",
		},
		{
			name: "With text only",
			setup: func(sb *StrBuilder) {
				sb.Text("Hello")
			},
			expected: "Hello\033[0m",
		},
		{
			name: "With color and text",
			setup: func(sb *StrBuilder) {
				sb.Red().Text("Error")
			},
			expected: "\033[31mError\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			tt.setup(sb)
			result := sb.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestMethodChaining(t *testing.T) {
	tests := []struct {
		name     string
		build    func() string
		expected string
	}{
		{
			name: "Bold red text",
			build: func() string {
				return NewStrBuilder().Bold().Red().Text("Error").String()
			},
			expected: "\033[1m\033[31mError\033[0m",
		},
		{
			name: "Underline blue text with reset",
			build: func() string {
				return NewStrBuilder().Underline().Blue().Text("Link").Reset().Text(" normal").String()
			},
			expected: "\033[4m\033[34mLink\033[0m normal\033[0m",
		},
		{
			name: "Multiple colors",
			build: func() string {
				return NewStrBuilder().Red().Text("R").Green().Text("G").Blue().Text("B").String()
			},
			expected: "\033[31mR\033[32mG\033[34mB\033[0m",
		},
		{
			name: "All styles combined",
			build: func() string {
				return NewStrBuilder().Bold().Italic().Underline().Dim().Text("styled").String()
			},
			expected: "\033[1m\033[3m\033[4m\033[2mstyled\033[0m",
		},
		{
			name: "RGB with style",
			build: func() string {
				return NewStrBuilder().Bold().RGB(255, 128, 0).Text("Orange").String()
			},
			expected: "\033[1m\033[38;2;255;128;0mOrange\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.build()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFluentInterface(t *testing.T) {
	sb := NewStrBuilder()

	// Test that each method returns the same pointer
	if sb.Bold() != sb {
		t.Error("Bold() should return the same StrBuilder instance")
	}
	if sb.Italic() != sb {
		t.Error("Italic() should return the same StrBuilder instance")
	}
	if sb.Red() != sb {
		t.Error("Red() should return the same StrBuilder instance")
	}
	if sb.Text("test") != sb {
		t.Error("Text() should return the same StrBuilder instance")
	}
	if sb.Reset() != sb {
		t.Error("Reset() should return the same StrBuilder instance")
	}
	if sb.RGB(0, 0, 0) != sb {
		t.Error("RGB() should return the same StrBuilder instance")
	}
}
