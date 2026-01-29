package beautyoutput

import (
	"testing"
)

func TestNewStrBuilder(t *testing.T) {
	sb := NewStrBuilder()
	if sb == nil {
		t.Error("NewStrBuilder should not return nil")
	}
	if sb.content.String() != "" {
		t.Errorf("NewStrBuilder content should be empty, got %q", sb.content.String())
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
			if sb.content.String() != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, sb.content.String(), tt.expected)
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
			if sb.content.String() != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, sb.content.String(), tt.expected)
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
			if sb.content.String() != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, sb.content.String(), tt.expected)
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
			if sb.content.String() != tt.expected {
				t.Errorf("RGB(%d, %d, %d) = %q, want %q", tt.r, tt.g, tt.b, sb.content.String(), tt.expected)
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
			if sb.content.String() != tt.expected {
				t.Errorf("Text(%q) = %q, want %q", tt.text, sb.content.String(), tt.expected)
			}
		})
	}
}

func TestReset(t *testing.T) {
	sb := NewStrBuilder()
	sb.Reset()
	expected := "\033[0m"
	if sb.content.String() != expected {
		t.Errorf("Reset() = %q, want %q", sb.content.String(), expected)
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

func TestTextf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{"Simple format", "Hello %s", []any{"World"}, "Hello World"},
		{"Multiple args", "%s is %d years old", []any{"John", 30}, "John is 30 years old"},
		{"No args", "Plain text", []any{}, "Plain text"},
		{"Float format", "Value: %.2f", []any{3.14159}, "Value: 3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStrBuilder()
			sb.Textf(tt.format, tt.args...)
			if sb.content.String() != tt.expected {
				t.Errorf("Textf(%q, %v) = %q, want %q", tt.format, tt.args, sb.content.String(), tt.expected)
			}
		})
	}
}

func TestNewLine(t *testing.T) {
	sb := NewStrBuilder()
	sb.Text("Line1").NewLine().Text("Line2")
	expected := "Line1\nLine2"
	if sb.content.String() != expected {
		t.Errorf("NewLine() = %q, want %q", sb.content.String(), expected)
	}
}

func TestNewLineFluentInterface(t *testing.T) {
	sb := NewStrBuilder()
	if sb.NewLine() != sb {
		t.Error("NewLine() should return the same StrBuilder instance")
	}
	if sb.Textf("test %s", "arg") != sb {
		t.Error("Textf() should return the same StrBuilder instance")
	}
}

func TestTemplateSuccess(t *testing.T) {
	sb := NewStrBuilder()
	sb.Success("Done")
	result := sb.String()
	if result != "\033[32m\033[1m✔ Done\033[0m\033[0m" {
		t.Errorf("Success() = %q, want %q", result, "\033[32m\033[1m✔ Done\033[0m\033[0m")
	}
}

func TestTemplateFailure(t *testing.T) {
	sb := NewStrBuilder()
	sb.Failure("Error occurred")
	result := sb.String()
	if result != "\033[31m\033[1m✖ Error occurred\033[0m\033[0m" {
		t.Errorf("Failure() = %q, want %q", result, "\033[31m\033[1m✖ Error occurred\033[0m\033[0m")
	}
}

func TestTemplateWarning(t *testing.T) {
	sb := NewStrBuilder()
	sb.Warning("Be careful")
	result := sb.String()
	if result != "\033[33m\033[1m⚠ Be careful\033[0m\033[0m" {
		t.Errorf("Warning() = %q, want %q", result, "\033[33m\033[1m⚠ Be careful\033[0m\033[0m")
	}
}

func TestTemplateInfo(t *testing.T) {
	sb := NewStrBuilder()
	sb.Info("Information")
	result := sb.String()
	if result != "\033[34m\033[1mℹ Information\033[0m\033[0m" {
		t.Errorf("Info() = %q, want %q", result, "\033[34m\033[1mℹ Information\033[0m\033[0m")
	}
}

func TestTemplatePending(t *testing.T) {
	sb := NewStrBuilder()
	sb.Pending("Loading")
	result := sb.String()
	if result != "\033[36m\033[1m… Loading\033[0m\033[0m" {
		t.Errorf("Pending() = %q, want %q", result, "\033[36m\033[1m… Loading\033[0m\033[0m")
	}
}

func TestTemplateLambda(t *testing.T) {
	sb := NewStrBuilder()
	sb.Lambda("Function")
	result := sb.String()
	if result != "\033[35m\033[1mλ Function\033[0m\033[0m" {
		t.Errorf("Lambda() = %q, want %q", result, "\033[35m\033[1mλ Function\033[0m\033[0m")
	}
}

func TestTemplateLambdaf(t *testing.T) {
	sb := NewStrBuilder()
	sb.Lambdaf("Value: %d", 42)
	result := sb.String()
	if result != "\033[35m\033[1mλ Value: 42\033[0m\033[0m" {
		t.Errorf("Lambdaf() = %q, want %q", result, "\033[35m\033[1mλ Value: 42\033[0m\033[0m")
	}
}

func TestTemplateTitleLine(t *testing.T) {
	sb := NewStrBuilder()
	sb.TitleLine("My Title")
	result := sb.String()
	expected := "\033[4m\033[1mMy Title\n\033[0m\033[0m"
	if result != expected {
		t.Errorf("TitleLine() = %q, want %q", result, expected)
	}
}

func TestTemplateTitleLinef(t *testing.T) {
	sb := NewStrBuilder()
	sb.TitleLinef("Title %s", "Test")
	result := sb.String()
	expected := "\033[4m\033[1mTitle Test\n\033[0m\033[0m"
	if result != expected {
		t.Errorf("TitleLinef() = %q, want %q", result, expected)
	}
}

func TestTemplateSectionLine(t *testing.T) {
	sb := NewStrBuilder()
	sb.SectionLine("Section")
	result := sb.String()
	expected := "\033[1mSection\033[0m\n\033[0m\033[0m"
	if result != expected {
		t.Errorf("SectionLine() = %q, want %q", result, expected)
	}
}

func TestTemplateList(t *testing.T) {
	sb := NewStrBuilder()
	sb.List([]string{"Item 1", "Item 2", "Item 3"})
	result := sb.String()
	expected := "• Item 1\n• Item 2\n• Item 3\n\033[0m"
	if result != expected {
		t.Errorf("List() = %q, want %q", result, expected)
	}
}

func TestTemplateListEmpty(t *testing.T) {
	sb := NewStrBuilder()
	sb.List([]string{})
	result := sb.String()
	expected := "\033[0m"
	if result != expected {
		t.Errorf("List([]) = %q, want %q", result, expected)
	}
}

func TestTemplateNumberedList(t *testing.T) {
	sb := NewStrBuilder()
	sb.NumberedList([]string{"First", "Second", "Third"})
	result := sb.String()
	expected := "1. First\n2. Second\n3. Third\n\033[0m"
	if result != expected {
		t.Errorf("NumberedList() = %q, want %q", result, expected)
	}
}

func TestTemplateNumberedListEmpty(t *testing.T) {
	sb := NewStrBuilder()
	sb.NumberedList([]string{})
	result := sb.String()
	expected := "\033[0m"
	if result != expected {
		t.Errorf("NumberedList([]) = %q, want %q", result, expected)
	}
}

func TestTemplateFluentInterface(t *testing.T) {
	sb := NewStrBuilder()

	if sb.Success("msg") != sb {
		t.Error("Success() should return the same StrBuilder instance")
	}

	sb2 := NewStrBuilder()
	if sb2.Failure("msg") != sb2 {
		t.Error("Failure() should return the same StrBuilder instance")
	}

	sb3 := NewStrBuilder()
	if sb3.Warning("msg") != sb3 {
		t.Error("Warning() should return the same StrBuilder instance")
	}

	sb4 := NewStrBuilder()
	if sb4.Info("msg") != sb4 {
		t.Error("Info() should return the same StrBuilder instance")
	}

	sb5 := NewStrBuilder()
	if sb5.Pending("msg") != sb5 {
		t.Error("Pending() should return the same StrBuilder instance")
	}

	sb6 := NewStrBuilder()
	if sb6.Lambda("msg") != sb6 {
		t.Error("Lambda() should return the same StrBuilder instance")
	}

	sb7 := NewStrBuilder()
	if sb7.Lambdaf("msg %s", "arg") != sb7 {
		t.Error("Lambdaf() should return the same StrBuilder instance")
	}

	sb8 := NewStrBuilder()
	if sb8.TitleLine("msg") != sb8 {
		t.Error("TitleLine() should return the same StrBuilder instance")
	}

	sb9 := NewStrBuilder()
	if sb9.TitleLinef("msg %s", "arg") != sb9 {
		t.Error("TitleLinef() should return the same StrBuilder instance")
	}

	sb10 := NewStrBuilder()
	if sb10.SectionLine("msg") != sb10 {
		t.Error("SectionLine() should return the same StrBuilder instance")
	}

	sb11 := NewStrBuilder()
	if sb11.List([]string{}) != sb11 {
		t.Error("List() should return the same StrBuilder instance")
	}

	sb12 := NewStrBuilder()
	if sb12.NumberedList([]string{}) != sb12 {
		t.Error("NumberedList() should return the same StrBuilder instance")
	}
}

func TestRealtimeOutputDefault(t *testing.T) {
	sb := NewStrBuilder()
	if sb.IsRealtimeOutput() {
		t.Error("IsRealtimeOutput() should be false by default")
	}
}

func TestSetRealtimeOutput(t *testing.T) {
	sb := NewStrBuilder()

	sb.SetRealtimeOutput(true)
	if !sb.IsRealtimeOutput() {
		t.Error("IsRealtimeOutput() should be true after SetRealtimeOutput(true)")
	}

	sb.SetRealtimeOutput(false)
	if sb.IsRealtimeOutput() {
		t.Error("IsRealtimeOutput() should be false after SetRealtimeOutput(false)")
	}
}

func TestSetRealtimeOutputFluentInterface(t *testing.T) {
	sb := NewStrBuilder()
	if sb.SetRealtimeOutput(true) != sb {
		t.Error("SetRealtimeOutput() should return the same StrBuilder instance")
	}
}
