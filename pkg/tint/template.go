package tint

import "fmt"

// PresetTitle formata a mensagem como título, com delimitadores "===" em negrito e ciano.
func PresetTitle(message string) string {
	str := fmt.Sprintf("=== %s ===", message)
	return NewStyle().
		Bold().
		Foreground(Cyan).
		Render(str)
}

// PresetSubtitle formata a mensagem como subtítulo, com delimitadores "---" em negrito e branco.
func PresetSubtitle(message string) string {
	str := fmt.Sprintf("--- %s ---", message)
	return NewStyle().
		Bold().
		Foreground(White).
		Render(str)
}

// PresetSuccess formata a mensagem como indicador de sucesso, com prefixo "✔" em negrito e verde.
func PresetSuccess(message string) string {
	str := fmt.Sprintf("✔ %s", message)
	return NewStyle().
		Bold().
		Foreground(Green).
		Render(str)
}

// PresetError formata a mensagem como indicador de erro, com prefixo "✖" em negrito e vermelho.
func PresetError(message string) string {
	str := fmt.Sprintf("✖ %s", message)
	return NewStyle().
		Bold().
		Foreground(Red).
		Render(str)
}

// PresetWarning formata a mensagem como indicador de aviso, com prefixo "⚠" em negrito e amarelo.
func PresetWarning(message string) string {
	str := fmt.Sprintf("⚠ %s", message)
	return NewStyle().
		Bold().
		Foreground(Yellow).
		Render(str)
}

// PresetInfo formata a mensagem como indicador informativo, com prefixo "ℹ" em negrito e azul.
func PresetInfo(message string) string {
	str := fmt.Sprintf("ℹ %s", message)
	return NewStyle().
		Bold().
		Foreground(Blue).
		Render(str)
}

// PresetLambda formata a mensagem com prefixo "λ" em negrito e magenta.
func PresetLambda(message string) string {
	str := fmt.Sprintf("λ %s", message)
	return NewStyle().
		Bold().
		Foreground(Magenta).
		Render(str)
}
