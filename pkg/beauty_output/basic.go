package beautyoutput

import "fmt"

// Text adiciona o texto especificado ao builder.
func (sb *StrBuilder) Text(text string) *StrBuilder {
	sb.write(text)
	return sb
}

func (sb *StrBuilder) Textf(format string, a ...interface{}) *StrBuilder {
	sb.write(fmt.Sprintf(format, a...))
	return sb
}

func (sb *StrBuilder) NewLine() *StrBuilder {
	sb.write("\n")
	return sb
}

// String retorna a string construída com um reset automático ao final.
func (sb *StrBuilder) String() string {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.content.String() + "\033[0m"
}
