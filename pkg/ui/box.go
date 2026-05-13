package ui

import "strings"

func Box(content string) string {
	lines := strings.Split(content, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	var sb strings.Builder
	sb.WriteString("┌" + strings.Repeat("─", maxWidth+2) + "┐\n")
	for _, line := range lines {
		sb.WriteString("│ " + line + strings.Repeat(" ", maxWidth-len(line)) + " │\n")
	}
	sb.WriteString("└" + strings.Repeat("─", maxWidth+2) + "┘")

	return sb.String()
}
