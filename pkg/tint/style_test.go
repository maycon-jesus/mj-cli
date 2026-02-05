package tint

import (
	"strings"
	"sync"
	"testing"
)

func TestNewStyle(t *testing.T) {
	s := NewStyle()

	if s == nil {
		t.Fatal("NewStyle() retornou nil")
	}

	// Verifica estado inicial - sem cores
	if s.fgColor != nil {
		t.Error("NewStyle() deveria ter fgColor nil")
	}
	if s.bgColor != nil {
		t.Error("NewStyle() deveria ter bgColor nil")
	}

	// Verifica estado inicial - sem estilos
	if s.bold {
		t.Error("NewStyle() deveria ter bold = false")
	}
	if s.dim {
		t.Error("NewStyle() deveria ter dim = false")
	}
	if s.italic {
		t.Error("NewStyle() deveria ter italic = false")
	}
	if s.underline {
		t.Error("NewStyle() deveria ter underline = false")
	}
	if s.reverse {
		t.Error("NewStyle() deveria ter reverse = false")
	}
	if s.hidden {
		t.Error("NewStyle() deveria ter hidden = false")
	}
	if s.strikethrough {
		t.Error("NewStyle() deveria ter strikethrough = false")
	}
}

func TestNewStyle_MultipleInstances(t *testing.T) {
	// Verifica que múltiplas instâncias são independentes
	s1 := NewStyle().Bold()
	s2 := NewStyle().Italic()

	if s1.bold != true {
		t.Error("s1 deveria ter bold = true")
	}
	if s2.bold != false {
		t.Error("s2 não deveria ser afetado por s1")
	}
	if s1.italic != false {
		t.Error("s1 não deveria ser afetado por s2")
	}
	if s2.italic != true {
		t.Error("s2 deveria ter italic = true")
	}
}

func TestStyle_Foreground(t *testing.T) {
	s := NewStyle()
	result := s.Foreground(Red)

	// Verifica encadeamento
	if result != s {
		t.Error("Foreground() deve retornar o mesmo ponteiro para encadeamento")
	}

	// Verifica que a cor foi definida
	if s.fgColor != Red {
		t.Errorf("Foreground() não definiu a cor corretamente, got %v, want %v", s.fgColor, Red)
	}
}

func TestStyle_Foreground_Override(t *testing.T) {
	s := NewStyle().Foreground(Red)
	s.Foreground(Blue)

	if s.fgColor != Blue {
		t.Errorf("Foreground() deveria sobrescrever a cor anterior, got %v, want %v", s.fgColor, Blue)
	}
}

func TestStyle_Background(t *testing.T) {
	s := NewStyle()
	result := s.Background(Blue)

	// Verifica encadeamento
	if result != s {
		t.Error("Background() deve retornar o mesmo ponteiro para encadeamento")
	}

	// Verifica que a cor foi definida
	if s.bgColor != Blue {
		t.Errorf("Background() não definiu a cor corretamente, got %v, want %v", s.bgColor, Blue)
	}
}

func TestStyle_Background_Override(t *testing.T) {
	s := NewStyle().Background(Blue)
	s.Background(Green)

	if s.bgColor != Green {
		t.Errorf("Background() deveria sobrescrever a cor anterior, got %v, want %v", s.bgColor, Green)
	}
}

func TestStyle_MethodChaining(t *testing.T) {
	s := NewStyle()
	result := s.Bold().Italic().Underline().Foreground(Red).Background(Blue)

	if result != s {
		t.Error("encadeamento de métodos deve sempre retornar o mesmo ponteiro")
	}

	// Verifica que todos os estilos foram aplicados
	if !s.bold {
		t.Error("Bold() não foi aplicado no encadeamento")
	}
	if !s.italic {
		t.Error("Italic() não foi aplicado no encadeamento")
	}
	if !s.underline {
		t.Error("Underline() não foi aplicado no encadeamento")
	}
	if s.fgColor != Red {
		t.Error("Foreground() não foi aplicado no encadeamento")
	}
	if s.bgColor != Blue {
		t.Error("Background() não foi aplicado no encadeamento")
	}
}

func TestStyle_Render_NoColor(t *testing.T) {
	s := NewStyle().Bold().Foreground(Red)
	message := "texto de teste"

	withColorSupport(t, ColorSupportNone, func() {
		got := s.Render(message)
		if got != message {
			t.Errorf("Render() com ColorSupportNone deveria retornar mensagem original, got %q, want %q", got, message)
		}
	})
}

func TestStyle_Render_NoStyles(t *testing.T) {
	s := NewStyle()
	message := "texto de teste"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Deve começar com reset
		if !strings.HasPrefix(got, "\033[0m") {
			t.Errorf("Render() deveria começar com reset, got %q", got)
		}

		// Deve terminar com reset
		if !strings.HasSuffix(got, "\033[0m") {
			t.Errorf("Render() deveria terminar com reset, got %q", got)
		}

		// Deve conter a mensagem
		if !strings.Contains(got, message) {
			t.Errorf("Render() deveria conter a mensagem, got %q", got)
		}

		// Formato esperado: \033[0m + message + \033[0m
		expected := "\033[0m" + message + "\033[0m"
		if got != expected {
			t.Errorf("Render() sem estilos = %q, want %q", got, expected)
		}
	})
}

func TestStyle_Render_Bold(t *testing.T) {
	s := NewStyle().Bold()
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Deve conter código de bold
		if !strings.Contains(got, "\033[1m") {
			t.Errorf("Render() com Bold deveria conter código 1, got %q", got)
		}

		// Verifica estrutura: reset + bold + message + reset
		expected := "\033[0m\033[1m" + message + "\033[0m"
		if got != expected {
			t.Errorf("Render() com Bold = %q, want %q", got, expected)
		}
	})
}

func TestStyle_Render_Foreground(t *testing.T) {
	s := NewStyle().Foreground(Red)
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Deve conter código de cor vermelha
		redCode := Red.GetForegroundCode()
		if !strings.Contains(got, redCode) {
			t.Errorf("Render() com Foreground(Red) deveria conter %q, got %q", redCode, got)
		}
	})
}

func TestStyle_Render_Background(t *testing.T) {
	s := NewStyle().Background(Blue)
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Deve conter código de cor de fundo azul
		blueCode := Blue.GetBackgroundCode()
		if !strings.Contains(got, blueCode) {
			t.Errorf("Render() com Background(Blue) deveria conter %q, got %q", blueCode, got)
		}
	})
}

func TestStyle_Render_MultipleStyles(t *testing.T) {
	s := NewStyle().Bold().Italic().Foreground(Red)
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Deve conter todos os códigos
		if !strings.Contains(got, "\033[1m") {
			t.Error("Render() deveria conter código bold")
		}
		if !strings.Contains(got, "\033[3m") {
			t.Error("Render() deveria conter código italic")
		}
		if !strings.Contains(got, "\033[31m") {
			t.Error("Render() deveria conter código de cor vermelha")
		}
	})
}

func TestStyle_Render_AllTextStyles(t *testing.T) {
	s := NewStyle().
		Bold().
		Dim().
		Italic().
		Underline().
		Reverse().
		Hidden().
		Strikethrough()
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Verifica presença de todos os códigos de estilo
		expectedCodes := []struct {
			code int
			name string
		}{
			{STYLE_BOLD, "bold"},
			{STYLE_DIM, "dim"},
			{STYLE_ITALIC, "italic"},
			{STYLE_UNDERLINE, "underline"},
			{STYLE_REVERSE, "reverse"},
			{STYLE_HIDDEN, "hidden"},
			{STYLE_STRIKETHROUGH, "strikethrough"},
		}

		for _, ec := range expectedCodes {
			codeStr := createAnsiCodeSequence(ec.code)
			if !strings.Contains(got, codeStr) {
				t.Errorf("Render() deveria conter código de %s (%q), got %q", ec.name, codeStr, got)
			}
		}
	})
}

func TestStyle_Render_EmptyMessage(t *testing.T) {
	s := NewStyle().Bold()

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render("")

		// Ainda deve ter os códigos ANSI, mesmo com mensagem vazia
		expected := "\033[0m\033[1m\033[0m"
		if got != expected {
			t.Errorf("Render() com mensagem vazia = %q, want %q", got, expected)
		}
	})
}

func TestStyle_Render_SpecialCharacters(t *testing.T) {
	s := NewStyle().Bold()
	testCases := []string{
		"texto com espaços",
		"texto\ncom\nnewlines",
		"texto\tcom\ttabs",
		"unicode: 日本語",
		"emojis: 🎉🚀",
	}

	withColorSupport(t, ColorSupport24bit, func() {
		for _, message := range testCases {
			got := s.Render(message)

			if !strings.Contains(got, message) {
				t.Errorf("Render() deveria preservar mensagem %q, got %q", message, got)
			}
		}
	})
}

func TestStyle_Render_WithTrueColor(t *testing.T) {
	trueColor := TrueColor{R: 255, G: 128, B: 64}
	s := NewStyle().Foreground(trueColor)
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		expectedCode := "\033[38;2;255;128;64m"
		if !strings.Contains(got, expectedCode) {
			t.Errorf("Render() com TrueColor deveria conter %q, got %q", expectedCode, got)
		}
	})
}

func TestStyle_Render_WithColorAnsi256(t *testing.T) {
	color256 := ColorAnsi256{Code: 208}
	s := NewStyle().Foreground(color256)
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		expectedCode := "\033[38;5;208m"
		if !strings.Contains(got, expectedCode) {
			t.Errorf("Render() com ColorAnsi256 deveria conter %q, got %q", expectedCode, got)
		}
	})
}

func TestStyle_Render_ForegroundAndBackground(t *testing.T) {
	s := NewStyle().Foreground(Red).Background(Blue)
	message := "texto"

	withColorSupport(t, ColorSupport24bit, func() {
		got := s.Render(message)

		// Deve conter código de cor vermelha (foreground)
		if !strings.Contains(got, "\033[31m") {
			t.Error("Render() deveria conter código de foreground vermelho")
		}

		// Deve conter código de cor azul (background)
		if !strings.Contains(got, "\033[44m") {
			t.Error("Render() deveria conter código de background azul")
		}
	})
}

// Testes de concorrência

func TestStyle_ConcurrentAccess(t *testing.T) {
	s := NewStyle()

	var wg sync.WaitGroup
	iterations := 100

	withColorSupport(t, ColorSupport24bit, func() {
		// Goroutines escrevendo (modificando estilos)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					switch id % 4 {
					case 0:
						s.Bold()
					case 1:
						s.Foreground(Red)
					case 2:
						s.Background(Blue)
					case 3:
						s.Italic()
					}
				}
			}(i)
		}

		// Goroutines lendo (renderizando)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					_ = s.Render("test")
				}
			}()
		}

		wg.Wait()
	})

	// Se chegou aqui sem panic/race condition, o teste passou
}

func TestStyle_ConcurrentRender(t *testing.T) {
	// Teste específico para leituras concorrentes
	s := NewStyle().Bold().Foreground(Red)

	var wg sync.WaitGroup
	results := make([]string, 100)

	withColorSupport(t, ColorSupport24bit, func() {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = s.Render("test")
			}(i)
		}

		wg.Wait()

		// Todos os resultados devem ser iguais
		expected := results[0]
		for i, r := range results {
			if r != expected {
				t.Errorf("resultado inconsistente no índice %d: %q != %q", i, r, expected)
			}
		}
	})
}

func TestStyle_ConcurrentModifyAndRender(t *testing.T) {
	// Testa modificação e leitura concorrentes
	s := NewStyle()

	var wg sync.WaitGroup

	withColorSupport(t, ColorSupport24bit, func() {
		// Uma goroutine modifica
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				s.Bold()
				s.Italic()
				s.Underline()
			}
		}()

		// Outra goroutine lê
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				_ = s.Render("concurrent test")
			}
		}()

		wg.Wait()
	})
}
