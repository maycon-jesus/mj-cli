// Package intl fornece suporte à internacionalização com gerenciamento de traduções thread-safe.
//
// O tipo Translator é seguro para uso concorrente por múltiplas goroutines.
// Todos os métodos públicos sincronizam corretamente o acesso ao estado interno.
package intl

import (
	"errors"
	"maps"
	"sync"

	"github.com/maycon-jesus/mj-cli/pkg/utils"
)

// Translator gerencia traduções para múltiplos idiomas.
// É seguro para uso concorrente por múltiplas goroutines.
type Translator struct {
	lang     string
	mu       sync.RWMutex
	messages Translations
}

// Translations mapeia códigos de idioma para seus dicionários de mensagens.
// Exemplo: {"en": {"greeting": "Hello"}, "pt-BR": {"greeting": "Olá"}}
type Translations map[string]map[string]string

// NewTranslator cria um novo Translator com o idioma especificado.
// Se lang for uma string vazia, o padrão será "en" (inglês).
func NewTranslator(lang string) *Translator {
	if lang == "" {
		lang = "en"
	}
	return &Translator{lang: lang, messages: make(Translations)}
}

// AddMessages adiciona mensagens para um idioma específico.
// Se o idioma ainda não existir, ele será criado automaticamente.
func (t *Translator) AddMessages(lang string, msgs map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.messages[lang]; !exists {
		t.messages[lang] = make(map[string]string)
	}
	maps.Copy(t.messages[lang], msgs)
}

// AddMessagesBulk adiciona mensagens para múltiplos idiomas de uma só vez.
// Aceita um mapa de Translations contendo todos os idiomas e suas mensagens.
func (t *Translator) AddMessagesBulk(langs Translations) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for lang, langMsgs := range langs {
		if _, exists := t.messages[lang]; !exists {
			t.messages[lang] = make(map[string]string)
		}
		maps.Copy(t.messages[lang], langMsgs)
	}
}

// SetLanguage define o idioma atual do Translator.
func (t *Translator) SetLanguage(lang string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lang = lang
}

// Language retorna o idioma atual configurado no Translator.
func (t *Translator) Language() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lang
}

// tLocked traduz uma chave para o idioma especificado.
// IMPORTANTE: O chamador deve manter pelo menos um lock de leitura (t.mu.RLock).
func (t *Translator) tLocked(lang string, key string, variables map[string]string) string {
	if msgs, ok := t.messages[lang]; ok {
		if msg, ok := msgs[key]; ok {
			return utils.ReplaceVariables(msg, variables)
		}
	}

	// Fallback para inglês se o idioma especificado não for encontrado
	if lang != "en" {
		return t.tLocked("en", key, variables)
	}

	return key
}

// T traduz uma chave para o idioma atual, substituindo variáveis se fornecidas.
// Retorna a própria chave se nenhuma tradução for encontrada.
func (t *Translator) T(key string, variables map[string]string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tLocked(t.lang, key, variables)
}

// Println traduz uma chave e imprime o resultado na saída padrão.
func (t *Translator) Println(key string, variables map[string]string) {
	msg := t.T(key, variables)
	println(msg)
}

// Errorf traduz uma chave e retorna como um error.
func (t *Translator) Errorf(key string, variables map[string]string) error {
	msg := t.T(key, variables)
	return errors.New(msg)
}
