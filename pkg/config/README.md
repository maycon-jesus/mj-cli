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

cm.SetEntry("debug", true)
cm.SetEntry("general.lang", "en")

value, _ := cm.GetEntry("general.lang") // "en"

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

O valor "efetivo" devolvido por `GetEntry` é `Value` quando não-nulo, caindo
para `DefaultValue` caso contrário.

### ConfigManager

Mantém o registro de entradas e delega persistência ao `PersistModule`
opcional.

| Método | Descrição |
|--------|-----------|
| `NewConfigManager()` | Cria um manager sem persister |
| `WithPersister(p)` | Associa um `PersistModule`; retorna o próprio manager para encadear |
| `AddEntry(key, desc, default)` | Registra uma nova chave; carrega o valor do persister se já existir |
| `GetEntry(key)` | Retorna o valor efetivo (ou `DefaultValue`) e se a chave existe |
| `SetEntry(key, value)` | Atualiza o valor de uma chave **já registrada**; retorna `false` se não existir |
| `Save()` | Persiste todas as entradas via o persister configurado |
| `Load()` | Recarrega o estado do persister |

`Save` e `Load` retornam erro se nenhum persister foi configurado.

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
2. Se a chave existir no backend, seu valor vira o `Value` inicial da entrada.
3. Caso contrário, `Value` fica `nil` e `GetEntry` devolverá o `DefaultValue`.

Isso permite registrar entradas em código e ainda assim respeitar valores já
gravados no arquivo de configuração, desde que `Load()` tenha sido chamado no
persister antes do `AddEntry`.

## Testes

Atualmente o pacote não possui testes automatizados.
