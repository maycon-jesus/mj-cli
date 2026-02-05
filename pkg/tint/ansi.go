package tint

import (
	"fmt"
	"os"
	"strings"
)

// ColorSupport representa o nível de suporte a cores do terminal.
// Os valores possíveis são ColorSupport4bit, ColorSupport8bit e ColorSupport24bit.
type ColorSupport int

// DetectedColorSupport contém o nível de suporte a cores detectado automaticamente
// no terminal atual. Esta variável é inicializada durante o carregamento do pacote.
//
// Será [ColorSupportNone] se a variável de ambiente NO_COLOR estiver definida,
// desabilitando toda saída de cores e estilos ANSI.
var DetectedColorSupport = detectColorSupport()

// detectColorSupport detecta o nível de suporte a cores do terminal.
//
// A ordem de verificação é:
//  1. NO_COLOR: se definida (qualquer valor), retorna [ColorSupportNone]
//  2. COLORTERM: verifica "truecolor" ou "24bit" para suporte 24-bit
//  3. TERM: verifica padrões como "256color", "xterm", "tmux", etc.
//  4. Fallback: [ColorSupport4bit] (16 cores básicas)
//
// Referência: https://no-color.org
func detectColorSupport() ColorSupport {
	noColor := os.Getenv("NO_COLOR")
	if noColor != "" {
		return ColorSupportNone
	}

	// Verifica COLORTERM primeiro (mais confiável para truecolor)
	colorTerm := os.Getenv("COLORTERM")
	if strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit") {
		return ColorSupport24bit
	}

	// Verifica TERM para outros níveis de suporte
	term := os.Getenv("TERM")

	// Truecolor (24-bit)
	if strings.Contains(term, "truecolor") || strings.Contains(term, "24bit") {
		return ColorSupport24bit
	}

	// 256 cores (8-bit)
	if strings.Contains(term, "256color") {
		return ColorSupport8bit
	}

	// Verifica se ao menos tem suporte básico a cores (assume 8-bit para terminais modernos)
	if strings.Contains(term, "color") ||
		strings.Contains(term, "ansi") ||
		strings.Contains(term, "xterm") ||
		strings.Contains(term, "screen") ||
		strings.Contains(term, "tmux") {
		return ColorSupport8bit
	}

	// Fallback para 16 cores básicas (4-bit)
	return ColorSupport4bit
}

// createAnsiCodeSequence cria uma sequência de escape ANSI a partir dos códigos fornecidos.
// Se nenhum código for passado, retorna a sequência de reset.
func createAnsiCodeSequence(codes ...int) string {
	if len(codes) == 0 {
		return "\033[m"
	}

	parts := make([]string, len(codes))
	for i, p := range codes {
		parts[i] = fmt.Sprintf("%d", p)
	}

	return fmt.Sprintf("\033[%sm", strings.Join(parts, ";"))
}
