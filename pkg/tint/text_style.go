package tint

// Bold ativa o estilo negrito no texto.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Bold() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bold = true
	return s
}

// Dim ativa o estilo de brilho reduzido (texto mais escuro) no texto.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Dim() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dim = true
	return s
}

// Italic ativa o estilo itálico no texto.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Italic() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.italic = true
	return s
}

// Underline ativa o estilo sublinhado no texto.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Underline() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.underline = true
	return s
}

// Reverse ativa o estilo de inversão de cores, trocando as cores
// de primeiro plano e fundo entre si.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Reverse() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reverse = true
	return s
}

// Hidden ativa o estilo oculto, tornando o texto invisível.
// O texto ainda ocupa espaço e pode ser selecionado/copiado.
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Hidden() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hidden = true
	return s
}

// Strikethrough ativa o estilo de texto riscado (tachado).
// Retorna o próprio Style para permitir encadeamento de métodos.
func (s *Style) Strikethrough() *Style {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strikethrough = true
	return s
}
