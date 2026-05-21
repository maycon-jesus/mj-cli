# utils

Utilitários internos do projeto. Reúne dois grupos de funções independentes:
substituição de variáveis com sintaxe mustache duplo (`{{nome}}`) e construção
de mapas aninhados (usados para montar a estrutura do arquivo de configuração).

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

### Comportamentos

| Situação | Resultado |
|----------|-----------|
| Mapa vazio ou `nil` em `ReplaceVariables` | String retornada sem alteração |
| `ReplaceFn` é `nil` | String retornada sem alteração |
| Placeholder sem valor no mapa | Mantido literalmente na saída |
| Callback retorna string vazia | Placeholder mantido literalmente na saída |
| Placeholder sem formato válido (ex: `{nome}`) | Ignorado, não é substituído |

Apenas identificadores com caracteres alfanuméricos e underscore são reconhecidos como variáveis (`\w+`).

## Mapas aninhados

Funções para construir e combinar `map[string]interface{}` aninhados. São usadas
pelo comando `config generate` para transformar chaves de configuração com pontos
(ex.: `server.host`) na estrutura hierárquica do arquivo YAML.

### CreateDeepMap

Converte uma lista de chaves em um mapa aninhado, onde cada chave vira um nível e
o `value` é colocado na chave mais profunda:

```go
m := utils.CreateDeepMap([]string{"server", "host"}, "localhost")
// → map[string]interface{}{
//       "server": map[string]interface{}{"host": "localhost"},
//   }
```

| Situação | Resultado |
|----------|-----------|
| Lista de chaves vazia ou `nil` | Retorna `nil` |
| Uma única chave | `{chave: value}` |
| Múltiplas chaves | Um nível de mapa por chave, `value` na folha |
| `value` é `nil` | A chave folha recebe `nil` normalmente |

### MergeDeepMaps

Combina `src` dentro de `dest` recursivamente e retorna `dest`:

```go
a := utils.CreateDeepMap([]string{"server", "host"}, "localhost")
b := utils.CreateDeepMap([]string{"server", "port"}, 8080)

merged := utils.MergeDeepMaps(a, b)
// → map[string]interface{}{
//       "server": map[string]interface{}{
//           "host": "localhost",
//           "port": 8080,
//       },
//   }
```

| Situação | Resultado |
|----------|-----------|
| Chave existe só em `src` | Adicionada a `dest` |
| Ambos os valores são mapas | Merge recursivo dos submapas |
| Valores conflitantes (não-mapas) | `src` sobrescreve `dest` |
| Um lado é mapa e o outro escalar | `src` sobrescreve `dest` (sem merge) |

> **Atenção:** `MergeDeepMaps` **muta** o mapa `dest` — o valor retornado é o
> próprio `dest`. Submapas vindos de `src` também são referenciados diretamente,
> sem cópia profunda.

## Testes

```bash
go test ./pkg/utils/ -v
```

Os testes cobrem:

- Substituição com mapa estático e com callback, incluindo mapas/funções `nil` e placeholders sem valor.
- `CreateDeepMap` com listas vazias, uma chave, múltiplos níveis e `value` nil.
- `MergeDeepMaps` com chaves disjuntas, sobrescrita escalar, merge recursivo e conflitos mapa↔escalar.
- Garantia de que `MergeDeepMaps` muta e retorna o próprio `dest`.
- Uso combinado de `CreateDeepMap` + `MergeDeepMaps` para montar configuração aninhada.
