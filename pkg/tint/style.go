package tint

import (
	"strings"
	"sync"
)

// Style representa um conjunto de estilos de texto que podem ser aplicados
// a uma string. Inclui cores de primeiro plano e fundo, além de formatações
// como negrito, itálico, sublinhado, entre outros.
//
// Style é seguro para uso concorrente através de mutex interno.
//
// Exemplo de uso:
//
//	style := tint.NewStyle().Bold().Foreground(tint.Red)
//	fmt.Println(style.Render("Texto estilizado"))
type Style struct {
	fgColor Color
	bgColor Color

	bold          bool
	dim           bool
	italic        bool
	underline     bool
	reverse       bool
	hidden        bool
	strikethrough bool

	mu sync.RWMutex
}

// NewStyle cria e retorna uma nova instância de Style sem nenhum estilo aplicado.
func NewStyle() *Style {
	return &Style{}
}

// Foreground define a cor do texto (primeiro plano).
// Retorna o próprio Style para permitir encadeamento de métodos.
func (sb *Style) Foreground(color Color) *Style {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.fgColor = color
	return sb
}

// Background define a cor de fundo do texto.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (sb *Style) Background(color Color) *Style {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.bgColor = color
	return sb
}

// Render aplica todos os estilos configurados à mensagem fornecida
// e retorna a string resultante com os códigos de escape ANSI.
//
// Se NO_COLOR estiver definido ([DetectedColorSupport] == [ColorSupportNone]),
// retorna a mensagem original sem nenhum código ANSI.
func (sb *Style) Render(message string) string {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	if DetectedColorSupport == ColorSupportNone {
		return message
	}

	var builder strings.Builder
	builder.WriteString(createAnsiCodeSequence(STYLE_RESET))

	if sb.bold {
		builder.WriteString(createAnsiCodeSequence(STYLE_BOLD))
	}
	if sb.dim {
		builder.WriteString(createAnsiCodeSequence(STYLE_DIM))
	}
	if sb.italic {
		builder.WriteString(createAnsiCodeSequence(STYLE_ITALIC))
	}
	if sb.underline {
		builder.WriteString(createAnsiCodeSequence(STYLE_UNDERLINE))
	}
	if sb.reverse {
		builder.WriteString(createAnsiCodeSequence(STYLE_REVERSE))
	}
	if sb.hidden {
		builder.WriteString(createAnsiCodeSequence(STYLE_HIDDEN))
	}
	if sb.strikethrough {
		builder.WriteString(createAnsiCodeSequence(STYLE_STRIKETHROUGH))
	}

	if sb.fgColor != nil {
		builder.WriteString(sb.fgColor.GetForegroundCode())
	}
	if sb.bgColor != nil {
		builder.WriteString(sb.bgColor.GetBackgroundCode())
	}

	builder.WriteString(message)
	builder.WriteString(createAnsiCodeSequence(STYLE_RESET))

	return builder.String()
}
