package beautyoutput

import "fmt"

// Text adiciona o texto especificado ao builder.
func (sb *StrBuilder) Text(text string) *StrBuilder {
	sb.content += text
	return sb
}

func (sb *StrBuilder) Textf(format string, a ...interface{}) *StrBuilder {
	sb.content += fmt.Sprintf(format, a...)
	return sb
}

func (sb *StrBuilder) NewLine() *StrBuilder {
	sb.content += "\n"
	return sb
}

// String retorna a string construída com um reset automático ao final.
func (sb *StrBuilder) String() string {
	return sb.content + "\033[0m"
}
