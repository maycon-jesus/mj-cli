package beautyoutput

// StrBuilder é um builder para construção de strings estilizadas com códigos ANSI.
// Utiliza o padrão fluent interface, permitindo encadeamento de métodos.
type StrBuilder struct {
	content string
}

// NewStrBuilder cria e retorna uma nova instância de StrBuilder.
func NewStrBuilder() *StrBuilder {
	return &StrBuilder{content: ""}
}
