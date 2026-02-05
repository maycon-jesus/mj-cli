# intl

Pacote de internacionalização (i18n). Fornece um `Translator` thread-safe para gerenciar e resolver traduções em múltiplos idiomas com suporte a substituição de variáveis.

## Uso básico

```go
t := intl.NewTranslator("pt-BR")

t.AddMessages("pt-BR", map[string]string{
    "greeting": "Olá, {{name}}!",
})
t.AddMessages("en", map[string]string{
    "greeting": "Hello, {{name}}!",
})

fmt.Println(t.T("greeting", map[string]string{"name": "mundo"}))
// Olá, mundo!
```

## Criação

```go
t := intl.NewTranslator("pt-BR") // idioma padrão: pt-BR
t := intl.NewTranslator("")      // string vazia → padrão "en"
```

## Adicionar mensagens

Para um único idioma:

```go
t.AddMessages("pt-BR", map[string]string{
    "app.start": "Iniciando...",
    "app.stop":  "Encerrando...",
})
```

Para múltiplos idiomas de uma só vez:

```go
t.AddMessagesBulk(intl.Translations{
    "pt-BR": {"app.start": "Iniciando...", "app.stop": "Encerrando..."},
    "en":    {"app.start": "Starting...",  "app.stop": "Stopping..."},
})
```

Chamadas repetidas para o mesmo idioma são aditivas — chaves existentes são sobrescritas, chaves novas são adicionadas.

## Tradução

### T

Resolve uma chave no idioma atual. Variáveis no formato `{{nome}}` são substituídas pelo mapa fornecido. Passa `nil` quando não há variáveis:

```go
msg := t.T("app.start", nil)
msg := t.T("greeting", map[string]string{"name": "mundo"})
```

### Println

Traduz e imprime diretamente na saída padrão:

```go
t.Println("app.start", nil)
```

### Errorf

Traduz e retorna como um `error`:

```go
err := t.Errorf("app.notfound", map[string]string{"id": "42"})
```

## Trocar idioma em tempo de execução

```go
t.SetLanguage("en")
fmt.Println(t.Language()) // "en"
```

## Fallback

A resolução de uma chave segue esta ordem:

1. Dicionário do idioma atual.
2. Dicionário `"en"` (inglês), se o idioma atual não for `"en"`.
3. A própria chave é retornada como último recurso.

```
chave "app.start", idioma "fr"
  → não existe "fr"   → tenta "en"
  → existe "en"       → retorna tradução em inglês
```

## Variáveis

Variáveis usam a sintaxe de mustache duplo `{{nome}}`. Qualquer placeholder sem valor correspondente no mapa é mantido literalmente na string de saída:

```go
t.AddMessages("en", map[string]string{
    "msg": "olá {{nome}}, você tem {{count}} mensagens",
})

t.T("msg", map[string]string{"nome": "Ana"})
// → "olá Ana, você tem {{count}} mensagens"
```

## Thread safety

Todos os métodos de `Translator` são seguros para uso concorrente. Leituras (`T`, `Language`) usam `RLock`; escritas (`AddMessages`, `SetLanguage`) usam `Lock`.
