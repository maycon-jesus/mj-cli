package tint

import (
	"strconv"
)

// Color é a interface que define os métodos para obter códigos de escape ANSI
// de cores para primeiro plano e plano de fundo.
type Color interface {
	// GetForegroundCode retorna a sequência de escape ANSI para definir a cor do texto.
	GetForegroundCode() string
	// GetBackgroundCode retorna a sequência de escape ANSI para definir a cor de fundo.
	GetBackgroundCode() string
}

// TrueColor representa uma cor no formato RGB (24-bit/truecolor).
// Cada componente (R, G, B) deve estar no intervalo de 0 a 255.
type TrueColor struct {
	R int
	G int
	B int
}

// GetForegroundCode retorna a sequência de escape ANSI para definir
// a cor do texto usando o formato RGB truecolor.
func (c TrueColor) GetForegroundCode() string {
	return "\033[38;2;" + strconv.Itoa(c.R) + ";" + strconv.Itoa(c.G) + ";" + strconv.Itoa(c.B) + "m"
}

// GetBackgroundCode retorna a sequência de escape ANSI para definir
// a cor de fundo usando o formato RGB truecolor.
func (c TrueColor) GetBackgroundCode() string {
	return "\033[48;2;" + strconv.Itoa(c.R) + ";" + strconv.Itoa(c.G) + ";" + strconv.Itoa(c.B) + "m"
}

// ColorAnsi256 representa uma cor da paleta de 256 cores (8-bit).
// O código deve estar no intervalo de 0 a 255.
type ColorAnsi256 struct {
	Code int
}

// GetForegroundCode retorna a sequência de escape ANSI para definir
// a cor do texto usando a paleta de 256 cores.
func (c ColorAnsi256) GetForegroundCode() string {
	return "\033[38;5;" + strconv.Itoa(c.Code) + "m"
}

// GetBackgroundCode retorna a sequência de escape ANSI para definir
// a cor de fundo usando a paleta de 256 cores.
func (c ColorAnsi256) GetBackgroundCode() string {
	return "\033[48;5;" + strconv.Itoa(c.Code) + "m"
}

// ColorAnsi representa uma cor ANSI básica (4-bit/16 cores).
// O código corresponde aos valores padrão ANSI (30-37 para cores normais,
// 90-97 para cores brilhantes).
type ColorAnsi struct {
	Code int
}

// GetForegroundCode retorna a sequência de escape ANSI para definir
// a cor do texto usando cores ANSI básicas.
func (c ColorAnsi) GetForegroundCode() string {
	return "\033[" + strconv.Itoa(c.Code) + "m"
}

// GetBackgroundCode retorna a sequência de escape ANSI para definir
// a cor de fundo usando cores ANSI básicas.
func (c ColorAnsi) GetBackgroundCode() string {
	return "\033[" + strconv.Itoa(c.Code+10) + "m"
}

// CompleteColor representa uma cor com fallbacks para diferentes níveis
// de suporte a cores. Permite definir uma cor em todos os formatos (TrueColor,
// 256 cores e ANSI básico), selecionando automaticamente o formato apropriado
// baseado no suporte detectado do terminal.
type CompleteColor struct {
	TrueColor TrueColor
	Ansi256   ColorAnsi256
	Ansi      ColorAnsi
}

// GetForegroundCode retorna a sequência de escape ANSI para definir a cor do texto,
// selecionando automaticamente o formato baseado no suporte detectado do terminal.
func (c CompleteColor) GetForegroundCode() string {
	switch DetectedColorSupport {
	case ColorSupportNone:
		return ""
	case ColorSupport4bit:
		return c.Ansi.GetForegroundCode()
	case ColorSupport8bit:
		return c.Ansi256.GetForegroundCode()
	case ColorSupport24bit:
		return c.TrueColor.GetForegroundCode()
	default:
		return c.TrueColor.GetForegroundCode()
	}
}

// GetBackgroundCode retorna a sequência de escape ANSI para definir a cor de fundo,
// selecionando automaticamente o formato baseado no suporte detectado do terminal.
func (c CompleteColor) GetBackgroundCode() string {
	switch DetectedColorSupport {
	case ColorSupportNone:
		return ""
	case ColorSupport4bit:
		return c.Ansi.GetBackgroundCode()
	case ColorSupport8bit:
		return c.Ansi256.GetBackgroundCode()
	case ColorSupport24bit:
		return c.TrueColor.GetBackgroundCode()
	default:
		return c.TrueColor.GetBackgroundCode()
	}
}

// Cores ANSI pré-definidas para uso conveniente.
// Inclui as 8 cores básicas e suas versões brilhantes.
var (
	Black         = ColorAnsi{Code: ANSI_COLOR_BLACK}
	Red           = ColorAnsi{Code: ANSI_COLOR_RED}
	Green         = ColorAnsi{Code: ANSI_COLOR_GREEN}
	Yellow        = ColorAnsi{Code: ANSI_COLOR_YELLOW}
	Blue          = ColorAnsi{Code: ANSI_COLOR_BLUE}
	Magenta       = ColorAnsi{Code: ANSI_COLOR_MAGENTA}
	Cyan          = ColorAnsi{Code: ANSI_COLOR_CYAN}
	White         = ColorAnsi{Code: ANSI_COLOR_WHITE}
	BrightBlack   = ColorAnsi{Code: ANSI_COLOR_BLACK + ANSI_COLOR_BRIGHT_OFFSET}
	BrightRed     = ColorAnsi{Code: ANSI_COLOR_RED + ANSI_COLOR_BRIGHT_OFFSET}
	BrightGreen   = ColorAnsi{Code: ANSI_COLOR_GREEN + ANSI_COLOR_BRIGHT_OFFSET}
	BrightYellow  = ColorAnsi{Code: ANSI_COLOR_YELLOW + ANSI_COLOR_BRIGHT_OFFSET}
	BrightBlue    = ColorAnsi{Code: ANSI_COLOR_BLUE + ANSI_COLOR_BRIGHT_OFFSET}
	BrightMagenta = ColorAnsi{Code: ANSI_COLOR_MAGENTA + ANSI_COLOR_BRIGHT_OFFSET}
	BrightCyan    = ColorAnsi{Code: ANSI_COLOR_CYAN + ANSI_COLOR_BRIGHT_OFFSET}
	BrightWhite   = ColorAnsi{Code: ANSI_COLOR_WHITE + ANSI_COLOR_BRIGHT_OFFSET}
)
