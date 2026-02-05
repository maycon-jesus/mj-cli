package tint

import (
	"strings"
	"testing"
)

func TestStyle_Bold(t *testing.T) {
	s := NewStyle()
	result := s.Bold()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Bold() deve retornar o mesmo Style")
	}

	// Verifica que bold foi ativado
	if !s.bold {
		t.Error("Bold() deve ativar o flag bold")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_BOLD)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Bold() deve adicionar código ANSI %d, got %q", STYLE_BOLD, rendered)
		}
	})
}

func TestStyle_Bold_Idempotent(t *testing.T) {
	// Aplicar múltiplas vezes não deve mudar o comportamento
	s := NewStyle().Bold().Bold().Bold()

	if !s.bold {
		t.Error("Bold() múltiplas vezes deve manter bold = true")
	}

	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		// Deve ter apenas um código de bold
		count := strings.Count(rendered, "\033[1m")
		if count != 1 {
			t.Errorf("Bold() múltiplas vezes deve ter apenas um código, got %d", count)
		}
	})
}

func TestStyle_Dim(t *testing.T) {
	s := NewStyle()
	result := s.Dim()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Dim() deve retornar o mesmo Style")
	}

	// Verifica que dim foi ativado
	if !s.dim {
		t.Error("Dim() deve ativar o flag dim")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_DIM)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Dim() deve adicionar código ANSI %d, got %q", STYLE_DIM, rendered)
		}
	})
}

func TestStyle_Italic(t *testing.T) {
	s := NewStyle()
	result := s.Italic()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Italic() deve retornar o mesmo Style")
	}

	// Verifica que italic foi ativado
	if !s.italic {
		t.Error("Italic() deve ativar o flag italic")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_ITALIC)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Italic() deve adicionar código ANSI %d, got %q", STYLE_ITALIC, rendered)
		}
	})
}

func TestStyle_Underline(t *testing.T) {
	s := NewStyle()
	result := s.Underline()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Underline() deve retornar o mesmo Style")
	}

	// Verifica que underline foi ativado
	if !s.underline {
		t.Error("Underline() deve ativar o flag underline")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_UNDERLINE)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Underline() deve adicionar código ANSI %d, got %q", STYLE_UNDERLINE, rendered)
		}
	})
}

func TestStyle_Reverse(t *testing.T) {
	s := NewStyle()
	result := s.Reverse()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Reverse() deve retornar o mesmo Style")
	}

	// Verifica que reverse foi ativado
	if !s.reverse {
		t.Error("Reverse() deve ativar o flag reverse")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_REVERSE)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Reverse() deve adicionar código ANSI %d, got %q", STYLE_REVERSE, rendered)
		}
	})
}

func TestStyle_Hidden(t *testing.T) {
	s := NewStyle()
	result := s.Hidden()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Hidden() deve retornar o mesmo Style")
	}

	// Verifica que hidden foi ativado
	if !s.hidden {
		t.Error("Hidden() deve ativar o flag hidden")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_HIDDEN)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Hidden() deve adicionar código ANSI %d, got %q", STYLE_HIDDEN, rendered)
		}
	})
}

func TestStyle_Strikethrough(t *testing.T) {
	s := NewStyle()
	result := s.Strikethrough()

	// Verifica retorno para encadeamento
	if result != s {
		t.Error("Strikethrough() deve retornar o mesmo Style")
	}

	// Verifica que strikethrough foi ativado
	if !s.strikethrough {
		t.Error("Strikethrough() deve ativar o flag strikethrough")
	}

	// Verifica que o código ANSI correto aparece no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")
		expectedCode := createAnsiCodeSequence(STYLE_STRIKETHROUGH)
		if !strings.Contains(rendered, expectedCode) {
			t.Errorf("Strikethrough() deve adicionar código ANSI %d, got %q", STYLE_STRIKETHROUGH, rendered)
		}
	})
}

func TestStyle_AllTextStylesCombination(t *testing.T) {
	s := NewStyle().
		Bold().
		Dim().
		Italic().
		Underline().
		Reverse().
		Hidden().
		Strikethrough()

	// Verifica todos os flags
	if !s.bold {
		t.Error("bold deveria estar ativo")
	}
	if !s.dim {
		t.Error("dim deveria estar ativo")
	}
	if !s.italic {
		t.Error("italic deveria estar ativo")
	}
	if !s.underline {
		t.Error("underline deveria estar ativo")
	}
	if !s.reverse {
		t.Error("reverse deveria estar ativo")
	}
	if !s.hidden {
		t.Error("hidden deveria estar ativo")
	}
	if !s.strikethrough {
		t.Error("strikethrough deveria estar ativo")
	}

	// Verifica que todos os códigos aparecem no Render
	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")

		styleChecks := []struct {
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

		for _, check := range styleChecks {
			expectedCode := createAnsiCodeSequence(check.code)
			if !strings.Contains(rendered, expectedCode) {
				t.Errorf("%s código (%q) não encontrado em %q", check.name, expectedCode, rendered)
			}
		}
	})
}

func TestStyle_TextStylesOrder(t *testing.T) {
	// Verifica que a ordem dos estilos no Render segue a ordem definida
	// (bold, dim, italic, underline, reverse, hidden, strikethrough)
	s := NewStyle().
		Bold().
		Dim().
		Italic().
		Underline().
		Reverse().
		Hidden().
		Strikethrough()

	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")

		// Encontra posições dos códigos
		positions := []struct {
			code int
			name string
			pos  int
		}{
			{STYLE_BOLD, "bold", -1},
			{STYLE_DIM, "dim", -1},
			{STYLE_ITALIC, "italic", -1},
			{STYLE_UNDERLINE, "underline", -1},
			{STYLE_REVERSE, "reverse", -1},
			{STYLE_HIDDEN, "hidden", -1},
			{STYLE_STRIKETHROUGH, "strikethrough", -1},
		}

		for i := range positions {
			codeStr := createAnsiCodeSequence(positions[i].code)
			positions[i].pos = strings.Index(rendered, codeStr)
			if positions[i].pos == -1 {
				t.Errorf("código de %s não encontrado", positions[i].name)
			}
		}

		// Verifica que cada código vem depois do anterior
		for i := 1; i < len(positions); i++ {
			if positions[i].pos != -1 && positions[i-1].pos != -1 {
				if positions[i].pos <= positions[i-1].pos {
					t.Errorf("%s deveria vir depois de %s na saída",
						positions[i].name, positions[i-1].name)
				}
			}
		}
	})
}

func TestStyle_TextStylesWithColors(t *testing.T) {
	// Testa combinação de estilos de texto com cores
	s := NewStyle().
		Bold().
		Italic().
		Foreground(Red).
		Background(Blue)

	withColorSupport(t, ColorSupport24bit, func() {
		rendered := s.Render("test")

		// Verifica estilos de texto
		if !strings.Contains(rendered, "\033[1m") {
			t.Error("Bold código não encontrado")
		}
		if !strings.Contains(rendered, "\033[3m") {
			t.Error("Italic código não encontrado")
		}

		// Verifica cores
		if !strings.Contains(rendered, "\033[31m") {
			t.Error("Red foreground código não encontrado")
		}
		if !strings.Contains(rendered, "\033[44m") {
			t.Error("Blue background código não encontrado")
		}

		// Verifica que a mensagem está presente
		if !strings.Contains(rendered, "test") {
			t.Error("mensagem não encontrada")
		}
	})
}

func TestStyle_TextStylesChainReturnsSamePointer(t *testing.T) {
	// Testa Bold
	s1 := NewStyle()
	if s1.Bold() != s1 {
		t.Error("Bold() não retornou o mesmo ponteiro")
	}

	// Testa Dim
	s2 := NewStyle()
	if s2.Dim() != s2 {
		t.Error("Dim() não retornou o mesmo ponteiro")
	}

	// Testa Italic
	s3 := NewStyle()
	if s3.Italic() != s3 {
		t.Error("Italic() não retornou o mesmo ponteiro")
	}

	// Testa Underline
	s4 := NewStyle()
	if s4.Underline() != s4 {
		t.Error("Underline() não retornou o mesmo ponteiro")
	}

	// Testa Reverse
	s5 := NewStyle()
	if s5.Reverse() != s5 {
		t.Error("Reverse() não retornou o mesmo ponteiro")
	}

	// Testa Hidden
	s6 := NewStyle()
	if s6.Hidden() != s6 {
		t.Error("Hidden() não retornou o mesmo ponteiro")
	}

	// Testa Strikethrough
	s7 := NewStyle()
	if s7.Strikethrough() != s7 {
		t.Error("Strikethrough() não retornou o mesmo ponteiro")
	}
}

func TestStyle_TextStylesIndependence(t *testing.T) {
	// Verifica que cada estilo é independente
	tests := []struct {
		name       string
		setup      func() *Style
		checkFlag  func(s *Style) bool
		checkOther func(s *Style) bool
	}{
		{
			name:       "apenas Bold",
			setup:      func() *Style { return NewStyle().Bold() },
			checkFlag:  func(s *Style) bool { return s.bold },
			checkOther: func(s *Style) bool { return !s.italic && !s.underline },
		},
		{
			name:       "apenas Italic",
			setup:      func() *Style { return NewStyle().Italic() },
			checkFlag:  func(s *Style) bool { return s.italic },
			checkOther: func(s *Style) bool { return !s.bold && !s.underline },
		},
		{
			name:       "apenas Underline",
			setup:      func() *Style { return NewStyle().Underline() },
			checkFlag:  func(s *Style) bool { return s.underline },
			checkOther: func(s *Style) bool { return !s.bold && !s.italic },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			if !tt.checkFlag(s) {
				t.Error("flag esperado não está ativo")
			}
			if !tt.checkOther(s) {
				t.Error("outros flags não deveriam estar ativos")
			}
		})
	}
}
