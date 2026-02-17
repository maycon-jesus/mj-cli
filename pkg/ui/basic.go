package ui

import (
	"fmt"

	"github.com/maycon-jesus/mj-cli/pkg/tint"
)

// Title formata a mensagem como título, com delimitadores "===" em negrito e ciano.
func Title(message string) string {
	str := fmt.Sprintf("=== %s ===", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Cyan).
		Render(str)
}

// Subtitle formata a mensagem como subtítulo, com delimitadores "---" em negrito e branco.
func Subtitle(message string) string {
	str := fmt.Sprintf("--- %s ---", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.White).
		Render(str)
}

// Success formata a mensagem como indicador de sucesso, com prefixo "✔" em negrito e verde.
func Success(message string) string {
	str := fmt.Sprintf("✔ %s", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Green).
		Render(str)
}

// Error formata a mensagem como indicador de erro, com prefixo "✖" em negrito e vermelho.
func Error(message string) string {
	str := fmt.Sprintf("✖ %s", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Red).
		Render(str)
}

// Warning formata a mensagem como indicador de aviso, com prefixo "⚠" em negrito e amarelo.
func Warning(message string) string {
	str := fmt.Sprintf("⚠ %s", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Yellow).
		Render(str)
}

// Info formata a mensagem como indicador informativo, com prefixo "ℹ" em negrito e azul.
func Info(message string) string {
	str := fmt.Sprintf("ℹ %s", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Blue).
		Render(str)
}

// Lambda formata a mensagem com prefixo "λ" em negrito e magenta.
func Lambda(message string) string {
	str := fmt.Sprintf("λ %s", message)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Magenta).
		Render(str)
}

// Lambdaf é a versão formatada de Lambda que aceita formatação de string com argumentos, similar a fmt.Sprintf.
func Lambdaf(message string, args ...any) string {
	str := fmt.Sprintf("λ "+message, args...)
	return tint.NewStyle().
		Bold().
		Foreground(tint.Magenta).
		Render(str)
}
