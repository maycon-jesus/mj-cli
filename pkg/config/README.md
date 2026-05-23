# config

Pacote de gerenciamento de configuração com chaves planas e separador `.` para
hierarquia. Cada entrada é declarada com uma descrição e um valor padrão; a
persistência é plugável via uma interface (`PersistModule`), e a implementação
incluída grava em YAML preservando a descrição como comentário.

## Uso rápido

```go
yamlPersister := config.NewYamlModule("./config.yaml")
yamlPersister.Load()

cm := config.NewConfigManager().WithPersister(yamlPersister)

cm.AddEntry("app_name", "Nome da aplicação", "MyApp")
cm.AddEntry("debug", "Habilita modo de debug", false)
cm.AddEntry("general.lang", "Idioma da UI", "pt-BR")
cm.AddEntry("git.signing.enabled", "Assinar commits", false)

cm.Set("debug", true)
cm.Set("general.lang", "en")

value, _ := cm.Get("general.lang") // "en"

if err := cm.Save(); err != nil {
    log.Fatal(err)
}
```

## Componentes

### Entry

Cada entrada registrada tem descrição, valor atual e valor padrão:

```go
type Entry struct {
    Description  string
    Value        any
    DefaultValue any
}
```

O valor "efetivo" devolvido por `Get` é `Value` quando não-nulo, caindo
para `DefaultValue` caso contrário.

### ConfigManager

Mantém o registro de entradas e delega persistência ao `PersistModule`
opcional.

| Método | Descrição |
|--------|-----------|
| `NewConfigManager()` | Cria um manager sem persister |
| `WithPersister(p)` | Associa um `PersistModule`; retorna o próprio manager para encadear |
| `AddEntry(key, desc, default)` | Registra uma nova chave; carrega o valor do persister se já existir |
| `Get(key)` | Retorna o valor efetivo (ou `DefaultValue`) e se a chave existe |
| `Set(key, value)` | Atualiza o valor de uma chave **já registrada**; retorna `false` se não existir |
| `Has(key)` | Indica se uma chave está registrada (independente do valor armazenado) |
| `Save()` | Persiste todas as entradas via o persister configurado |
| `Load()` | Recarrega o estado do persister |

`Save` e `Load` retornam erro se nenhum persister foi configurado.

`Has` checa existência exata da chave registrada — não considera ancestrais.
Por exemplo, com `git.branch` registrado, `Has("git.branch")` retorna `true`
mas `Has("git")` retorna `false`.

### ReadOnlyConfig

Interface mínima exposta para consumidores que só precisam **ler** a
configuração, sem poder mutá-la ou persistir:

```go
type ReadOnlyConfig interface {
    Get(key string) (any, bool)
    Has(key string) bool
}
```

`*ConfigManager` satisfaz `ReadOnlyConfig`, então basta tipar o parâmetro da
função consumidora com essa interface para impedir chamadas a `Set`, `Save`,
`AddEntry` etc.

### PersistModule

Interface que abstrai o backend de persistência:

```go
type PersistModule interface {
    Save(configs map[string]Entry) error
    Load() error
    Get(key string) (any, bool)
}
```

- `Get` é consultada por `AddEntry` para popular o `Value` inicial a partir do
  estado já persistido.
- `Save` recebe o mapa completo de entradas.
- `Load` recarrega o estado interno do backend.

### YamlModule

Implementação de `PersistModule` que persiste em um arquivo YAML.

```go
yaml := config.NewYamlModule("./config.yaml")
yaml.Load()
```

Características:

- Chaves com `.` viram níveis aninhados no YAML (ex.: `git.signing.enabled` →
  `git: { signing: { enabled: ... } }`).
- As chaves são serializadas em ordem alfabética para gerar saída
  determinística.
- A `Description` de cada entrada é escrita como **head comment** acima da
  chave folha.
- O arquivo é gravado com permissão `0644`.
- `Load` em um arquivo inexistente não é erro: o estado interno fica vazio.

## Regras de chaves

`AddEntry` valida hierarquia para impedir conflitos folha/nó:

| Situação | Resultado |
|----------|-----------|
| Chave já registrada | Erro: `entry with key %q already exists` |
| Um ancestral da chave já é folha (ex.: `git` registrado, tenta-se `git.branch`) | Erro: `ancestor %q is a leaf; cannot register %q` |
| A chave proposta tem descendentes registrados (ex.: `git.branch` existe, tenta-se `git`) | Erro: `descendant %q exists; cannot register leaf %q` |

Essas regras garantem que toda chave registrada possa ser representada como
folha na árvore YAML sem ambiguidade.

## Carregamento de valores existentes

Ao chamar `AddEntry` com um persister associado:

1. O manager consulta `persister.Get(key)`.
2. Se a chave existir no backend, seu valor vira o `Value` inicial da entrada
   (mesmo que esse valor seja `nil` — nesse caso `Get` ainda cai no
   `DefaultValue` graças ao fallback de `effective()`).
3. Se a chave não existir no backend, `Value` fica `nil` e `Get` devolverá
   o `DefaultValue`.

Isso permite registrar entradas em código e ainda assim respeitar valores já
gravados no arquivo de configuração, desde que `Load()` tenha sido chamado no
persister antes do `AddEntry`.

## Testes

```bash
go test ./pkg/config/ -v
```

Os testes cobrem:

- **`Entry.effective()`** — prioridade entre `Value` e `DefaultValue`,
  incluindo casos delicados (`false`, `0`, ambos `nil`).
- **`ConfigManager`** — construção, encadeamento de `WithPersister`,
  validação de hierarquia em `AddEntry` (duplicata, ancestral folha,
  descendente existente, irmãos aninhados), carga inicial via persister,
  `Get`/`Set`/`Has` (incluindo o caso em que `Has` não confunde chave com
  ancestral) e propagação/erros de `Save`/`Load`. Um `fakePersister` em
  memória evita acesso a disco.
- **`YamlModule`** — gravação plana e aninhada, descrição como head
  comment, ordem alfabética determinística, uso do `DefaultValue` quando
  `Value` é `nil`, erros em paths inválidos, e roundtrip Save→Load para
  tipos variados (string, int, float, bool, slices, nil).
