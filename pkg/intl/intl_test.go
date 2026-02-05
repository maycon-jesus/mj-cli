package intl

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewTranslator(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		wantLang string
	}{
		{
			name:     "idioma fornecido",
			lang:     "pt-BR",
			wantLang: "pt-BR",
		},
		{
			name:     "string vazia — padrão en",
			lang:     "",
			wantLang: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTranslator(tt.lang)
			if got := tr.Language(); got != tt.wantLang {
				t.Errorf("NewTranslator(%q).Language() = %q, want %q", tt.lang, got, tt.wantLang)
			}
		})
	}
}

func TestSetLanguage(t *testing.T) {
	tr := NewTranslator("en")
	tr.SetLanguage("pt-BR")
	if got := tr.Language(); got != "pt-BR" {
		t.Errorf("Language() após SetLanguage(\"pt-BR\") = %q, want %q", got, "pt-BR")
	}
}

func TestAddMessages(t *testing.T) {
	t.Run("adiciona mensagens para idioma novo", func(t *testing.T) {
		tr := NewTranslator("pt-BR")
		tr.AddMessages("pt-BR", map[string]string{"saudação": "Olá"})
		if got := tr.T("saudação", nil); got != "Olá" {
			t.Errorf("T(\"saudação\") = %q, want %q", got, "Olá")
		}
	})

	t.Run("chaves novas coexistem com as antigas", func(t *testing.T) {
		tr := NewTranslator("en")
		tr.AddMessages("en", map[string]string{"a": "A"})
		tr.AddMessages("en", map[string]string{"b": "B"})
		if got := tr.T("a", nil); got != "A" {
			t.Errorf("T(\"a\") = %q, want %q", got, "A")
		}
		if got := tr.T("b", nil); got != "B" {
			t.Errorf("T(\"b\") = %q, want %q", got, "B")
		}
	})

	t.Run("sobrescreve chave existente", func(t *testing.T) {
		tr := NewTranslator("en")
		tr.AddMessages("en", map[string]string{"a": "antigo"})
		tr.AddMessages("en", map[string]string{"a": "novo"})
		if got := tr.T("a", nil); got != "novo" {
			t.Errorf("T(\"a\") = %q, want %q", got, "novo")
		}
	})

	t.Run("msgs nil não panicia", func(t *testing.T) {
		tr := NewTranslator("en")
		tr.AddMessages("en", nil)
	})
}

func TestAddMessagesBulk(t *testing.T) {
	t.Run("adiciona múltiplos idiomas", func(t *testing.T) {
		tr := NewTranslator("pt-BR")
		tr.AddMessagesBulk(Translations{
			"pt-BR": {"saudação": "Olá"},
			"en":    {"saudação": "Hello"},
		})
		if got := tr.T("saudação", nil); got != "Olá" {
			t.Errorf("T(\"saudação\") com lang pt-BR = %q, want %q", got, "Olá")
		}
		tr.SetLanguage("en")
		if got := tr.T("saudação", nil); got != "Hello" {
			t.Errorf("T(\"saudação\") com lang en = %q, want %q", got, "Hello")
		}
	})

	t.Run("bulk é aditivo", func(t *testing.T) {
		tr := NewTranslator("en")
		tr.AddMessages("en", map[string]string{"a": "A"})
		tr.AddMessagesBulk(Translations{
			"en": {"b": "B"},
		})
		if got := tr.T("a", nil); got != "A" {
			t.Errorf("T(\"a\") = %q, want %q", got, "A")
		}
		if got := tr.T("b", nil); got != "B" {
			t.Errorf("T(\"b\") = %q, want %q", got, "B")
		}
	})

	t.Run("sobrescreve chaves existentes", func(t *testing.T) {
		tr := NewTranslator("en")
		tr.AddMessages("en", map[string]string{"a": "antigo"})
		tr.AddMessagesBulk(Translations{
			"en": {"a": "novo"},
		})
		if got := tr.T("a", nil); got != "novo" {
			t.Errorf("T(\"a\") = %q, want %q", got, "novo")
		}
	})

	t.Run("Translations nil não panicia", func(t *testing.T) {
		tr := NewTranslator("en")
		tr.AddMessagesBulk(nil)
	})
}

func TestT(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *Translator
		key       string
		variables map[string]string
		want      string
	}{
		{
			name: "tradução básica",
			setup: func() *Translator {
				tr := NewTranslator("pt-BR")
				tr.AddMessages("pt-BR", map[string]string{"saudação": "Olá"})
				return tr
			},
			key:  "saudação",
			want: "Olá",
		},
		{
			name: "tradução com variáveis",
			setup: func() *Translator {
				tr := NewTranslator("pt-BR")
				tr.AddMessages("pt-BR", map[string]string{"saudação": "Olá, {{nome}}!"})
				return tr
			},
			key:       "saudação",
			variables: map[string]string{"nome": "Ana"},
			want:      "Olá, Ana!",
		},
		{
			name: "variáveis nil — placeholder mantido",
			setup: func() *Translator {
				tr := NewTranslator("en")
				tr.AddMessages("en", map[string]string{"msg": "olá {{nome}}"})
				return tr
			},
			key:  "msg",
			want: "olá {{nome}}",
		},
		{
			name: "fallback para en quando chave não existe no idioma atual",
			setup: func() *Translator {
				tr := NewTranslator("pt-BR")
				tr.AddMessages("en", map[string]string{"chave": "english value"})
				return tr
			},
			key:  "chave",
			want: "english value",
		},
		{
			name: "chave existe no idioma atual e em en — usa o do idioma atual",
			setup: func() *Translator {
				tr := NewTranslator("pt-BR")
				tr.AddMessages("pt-BR", map[string]string{"chave": "valor pt"})
				tr.AddMessages("en", map[string]string{"chave": "valor en"})
				return tr
			},
			key:  "chave",
			want: "valor pt",
		},
		{
			name: "chave não existe em nenhum idioma — retorna a própria chave",
			setup: func() *Translator {
				tr := NewTranslator("pt-BR")
				tr.AddMessages("pt-BR", map[string]string{"outra": "valor"})
				return tr
			},
			key:  "inexistente",
			want: "inexistente",
		},
		{
			name: "idioma atual é en e chave não existe — retorna a chave",
			setup: func() *Translator {
				return NewTranslator("en")
			},
			key:  "chave_ausente",
			want: "chave_ausente",
		},
		{
			name: "idioma atual não existe como dicionário mas en existe",
			setup: func() *Translator {
				tr := NewTranslator("fr")
				tr.AddMessages("en", map[string]string{"chave": "english"})
				return tr
			},
			key:  "chave",
			want: "english",
		},
		{
			name: "fallback para en também substitui variáveis",
			setup: func() *Translator {
				tr := NewTranslator("pt-BR")
				tr.AddMessages("en", map[string]string{"msg": "hello {{nome}}"})
				return tr
			},
			key:       "msg",
			variables: map[string]string{"nome": "Bob"},
			want:      "hello Bob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.setup()
			if got := tr.T(tt.key, tt.variables); got != tt.want {
				t.Errorf("T(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *Translator
		key       string
		variables map[string]string
		wantMsg   string
	}{
		{
			name: "error com mensagem traduzida",
			setup: func() *Translator {
				tr := NewTranslator("en")
				tr.AddMessages("en", map[string]string{"err.notfound": "not found"})
				return tr
			},
			key:     "err.notfound",
			wantMsg: "not found",
		},
		{
			name: "error com variáveis substituídas",
			setup: func() *Translator {
				tr := NewTranslator("en")
				tr.AddMessages("en", map[string]string{"err.id": "id {{id}} não encontrado"})
				return tr
			},
			key:       "err.id",
			variables: map[string]string{"id": "42"},
			wantMsg:   "id 42 não encontrado",
		},
		{
			name: "chave não encontrada — error com a própria chave",
			setup: func() *Translator {
				return NewTranslator("en")
			},
			key:     "chave.ausente",
			wantMsg: "chave.ausente",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := tt.setup()
			err := tr.Errorf(tt.key, tt.variables)
			if err == nil {
				t.Fatal("Errorf() retornou nil, esperava error")
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("Errorf().Error() = %q, want %q", err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestPrintln(t *testing.T) {
	// Println usa println (builtin → stderr). Verifica apenas que não panicia.
	tr := NewTranslator("en")
	tr.AddMessages("en", map[string]string{"msg": "hello {{nome}}"})
	tr.Println("msg", map[string]string{"nome": "teste"})
	tr.Println("chave_ausente", nil)
}

// --- concorrência (executar com go test -race) ---

func TestT_concorrente(t *testing.T) {
	tr := NewTranslator("pt-BR")
	tr.AddMessagesBulk(Translations{
		"pt-BR": {"msg": "olá {{nome}}"},
		"en":    {"msg": "hello {{nome}}"},
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr.T("msg", map[string]string{"nome": "goroutine"})
		}()
	}
	wg.Wait()
}

func TestAddMessages_concorrente(t *testing.T) {
	tr := NewTranslator("en")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tr.AddMessages("en", map[string]string{
				fmt.Sprintf("chave_%d", i): fmt.Sprintf("valor_%d", i),
			})
		}(i)
	}
	wg.Wait()
}

func TestSetLanguage_concorrente(t *testing.T) {
	tr := NewTranslator("en")
	langs := []string{"en", "pt-BR", "fr", "es"}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tr.SetLanguage(langs[i%len(langs)])
			tr.Language()
		}(i)
	}
	wg.Wait()
}

func TestAddMessagesBulk_concorrente(t *testing.T) {
	tr := NewTranslator("en")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tr.AddMessagesBulk(Translations{
				"en": {fmt.Sprintf("chave_%d", i): fmt.Sprintf("valor_%d", i)},
			})
		}(i)
	}
	wg.Wait()
}
