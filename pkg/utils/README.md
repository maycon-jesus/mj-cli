# utils

Utilitários internos do projeto. Atualmente contém funções de substituição de variáveis com sintaxe mustache duplo (`{{nome}}`).

## Substituição de variáveis

### ReplaceVariables

Substitui placeholders `{{nome}}` utilizando um mapa estático de valores:

```go
result := utils.ReplaceVariables("olá {{nome}}, porta {{port}}", map[string]string{
    "nome": "Ana",
    "port": "8080",
})
// → "olá Ana, porta 8080"
```

### ReplaceVariablesWithFn

Substitui placeholders utilizando uma função callback. Útil quando os valores são resolvidos dinamicamente:

```go
result := utils.ReplaceVariablesWithFn("olá {{nome}}", func(variableName string) string {
    return db.GetValue(variableName) // resolução dinâmica
})
// → "olá <valor retornado pela fn>"
```

O tipo da callback é `utils.ReplaceFn`:

```go
type ReplaceFn func(variableName string) string
```

## Comportamentos

| Situação | Resultado |
|----------|-----------|
| Mapa vazio ou `nil` em `ReplaceVariables` | String retornada sem alteração |
| `ReplaceFn` é `nil` | String retornada sem alteração |
| Placeholder sem valor no mapa | Mantido literalmente na saída |
| Callback retorna string vazia | Placeholder mantido literalmente na saída |
| Placeholder sem formato válido (ex: `{nome}`) | Ignorado, não é substituído |

Apenas identificadores com caracteres alfanuméricos e underscore são reconhecidos como variáveis (`\w+`).
