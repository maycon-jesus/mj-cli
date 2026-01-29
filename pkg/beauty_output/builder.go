package beautyoutput

import (
	"strings"
	"sync"
)

// StrBuilder é um builder para construção de strings estilizadas com códigos ANSI.
// Utiliza o padrão fluent interface, permitindo encadeamento de métodos.
type StrBuilder struct {
	content strings.Builder
	mu      sync.RWMutex

	// Indica se a saída em tempo real está habilitada.
	RealtimeOutput bool
}

// NewStrBuilder cria e retorna uma nova instância de StrBuilder.
func NewStrBuilder() *StrBuilder {
	return &StrBuilder{content: strings.Builder{}}
}

func (sb *StrBuilder) write(text string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.content.WriteString(text)
	sb.printRealtimeOutputUnsafe()
}
