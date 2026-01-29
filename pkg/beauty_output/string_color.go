// Package beautyoutput fornece um builder fluente para construção de strings
// com estilos e cores ANSI para saída em terminal.
package beautyoutput

import (
	"strconv"
)

// Bold aplica o estilo negrito ao texto subsequente.
func (sb *StrBuilder) Bold() *StrBuilder {
	sb.write("\033[1m")
	return sb
}

// Italic aplica o estilo itálico ao texto subsequente.
func (sb *StrBuilder) Italic() *StrBuilder {
	sb.write("\033[3m")
	return sb
}

// Underline aplica o estilo sublinhado ao texto subsequente.
func (sb *StrBuilder) Underline() *StrBuilder {
	sb.write("\033[4m")
	return sb
}

// Dim aplica o estilo de baixa intensidade (opaco) ao texto subsequente.
func (sb *StrBuilder) Dim() *StrBuilder {
	sb.write("\033[2m")
	return sb
}

// Black aplica a cor preta ao texto subsequente.
func (sb *StrBuilder) Black() *StrBuilder {
	sb.write("\033[30m")
	return sb
}

// Red aplica a cor vermelha ao texto subsequente.
func (sb *StrBuilder) Red() *StrBuilder {
	sb.write("\033[31m")
	return sb
}

// Green aplica a cor verde ao texto subsequente.
func (sb *StrBuilder) Green() *StrBuilder {
	sb.write("\033[32m")
	return sb
}

// Yellow aplica a cor amarela ao texto subsequente.
func (sb *StrBuilder) Yellow() *StrBuilder {
	sb.write("\033[33m")
	return sb
}

// Blue aplica a cor azul ao texto subsequente.
func (sb *StrBuilder) Blue() *StrBuilder {
	sb.write("\033[34m")
	return sb
}

// Magenta aplica a cor magenta ao texto subsequente.
func (sb *StrBuilder) Magenta() *StrBuilder {
	sb.write("\033[35m")
	return sb
}

// Cyan aplica a cor ciano ao texto subsequente.
func (sb *StrBuilder) Cyan() *StrBuilder {
	sb.write("\033[36m")
	return sb
}

// White aplica a cor branca ao texto subsequente.
func (sb *StrBuilder) White() *StrBuilder {
	sb.write("\033[37m")
	return sb
}

// BrightBlack aplica a cor preta brilhante (cinza) ao texto subsequente.
func (sb *StrBuilder) BrightBlack() *StrBuilder {
	sb.write("\033[90m")
	return sb
}

// BrightRed aplica a cor vermelha brilhante ao texto subsequente.
func (sb *StrBuilder) BrightRed() *StrBuilder {
	sb.write("\033[91m")
	return sb
}

// BrightGreen aplica a cor verde brilhante ao texto subsequente.
func (sb *StrBuilder) BrightGreen() *StrBuilder {
	sb.write("\033[92m")
	return sb
}

// BrightYellow aplica a cor amarela brilhante ao texto subsequente.
func (sb *StrBuilder) BrightYellow() *StrBuilder {
	sb.write("\033[93m")
	return sb
}

// BrightBlue aplica a cor azul brilhante ao texto subsequente.
func (sb *StrBuilder) BrightBlue() *StrBuilder {
	sb.write("\033[94m")
	return sb
}

// BrightMagenta aplica a cor magenta brilhante ao texto subsequente.
func (sb *StrBuilder) BrightMagenta() *StrBuilder {
	sb.write("\033[95m")
	return sb
}

// BrightCyan aplica a cor ciano brilhante ao texto subsequente.
func (sb *StrBuilder) BrightCyan() *StrBuilder {
	sb.write("\033[96m")
	return sb
}

// BrightWhite aplica a cor branca brilhante ao texto subsequente.
func (sb *StrBuilder) BrightWhite() *StrBuilder {
	sb.write("\033[97m")
	return sb
}

// RGB aplica uma cor customizada usando valores RGB (0-255) ao texto subsequente.
func (sb *StrBuilder) RGB(r, g, b int) *StrBuilder {
	r = max(0, min(255, r))
	g = max(0, min(255, g))
	b = max(0, min(255, b))
	sb.write("\033[38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m")
	return sb
}

// Reset adiciona o código de reset ANSI, removendo todos os estilos e cores aplicados.
func (sb *StrBuilder) Reset() *StrBuilder {
	sb.write("\033[0m")
	return sb
}
