package cmd

import (
	"fmt"
	"strings"
	"unicode"
)

// parseCommand divide uma string de comando em argumentos, respeitando
// strings entre aspas e caracteres de escape. Aspas duplas permitem
// sequências de escape (\"), enquanto aspas simples tratam tudo
// literalmente (como no bash).
func parseCommand(command string) ([]string, error) {
	var args []string
	var current strings.Builder
	var inQuote rune
	escaped := false

	for _, char := range command {
		if escaped {
			if inQuote == '"' {
				// Dentro de aspas duplas, apenas certos escapes são válidos
				if char == '"' || char == '\\' || char == '$' || char == '`' {
					current.WriteRune(char)
				} else {
					// Preserva a barra invertida para escapes inválidos
					current.WriteRune('\\')
					current.WriteRune(char)
				}
			} else {
				// Fora de aspas, escape torna o próximo caractere literal
				current.WriteRune(char)
			}
			escaped = false
			continue
		}

		if char == '\\' {
			// Em aspas simples, barra invertida é literal
			if inQuote == '\'' {
				current.WriteRune(char)
			} else {
				escaped = true
			}
			continue
		}

		if char == '"' || char == '\'' {
			if inQuote == char {
				inQuote = 0
			} else if inQuote == 0 {
				inQuote = char
			} else {
				// Aspas diferentes dentro de uma string entre aspas
				current.WriteRune(char)
			}
			continue
		}

		if unicode.IsSpace(char) && inQuote == 0 {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}

	if escaped {
		return nil, fmt.Errorf("unexpected end of command: trailing backslash")
	}

	if inQuote != 0 {
		return nil, fmt.Errorf("unclosed quote: %c", inQuote)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args, nil
}
