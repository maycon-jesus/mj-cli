# logger

Pacote de logging JSON estruturado construído sobre o [`log/slog`](https://pkg.go.dev/log/slog) da biblioteca padrão do Go.

Fornece duas implementações customizadas de `slog.Handler` — `FileHandler` e `MultiHandler` — que podem ser compostas para distribuir registros de log para múltiplos destinos.

## Uso rápido

```go
lgr, err := logger.NewLoggerComplete("mj-cli", "1.0.0")
if err != nil {
    log.Fatal(err)
}
defer lgr.FileHandler.Close()

lgr.Log.Info("aplicação iniciada")
lgr.Log.Warn("algo inesperado", "detalhe", "valor")
lgr.Log.Error("falha ao processar", "err", err)
```

O arquivo de log é criado no diretório temporário do sistema operacional com o nome `mj-cli-log-*`. Para descobrir o caminho:

```go
fmt.Println(lgr.FileHandler.Name()) // ex: /tmp/mj-cli-log-1234567890
```

## Componentes

### Logger

Struct principal que agrupa um `*slog.Logger` e o `*FileHandler` subjacente:

```go
type Logger struct {
    Log         *slog.Logger
    FileHandler *FileHandler
}
```

Criado via `NewLoggerComplete(appName, appVersion)`, que:

1. Cria um arquivo temporário de log
2. Define o nível padrão como `Debug`
3. Adiciona os atributos `app` e `version` a todos os registros
4. Verifica a variável de ambiente `DEBUG` para override de nível

### FileHandler

Handler thread-safe que escreve registros como objetos JSON (um por linha) em um `*os.File`.

```go
// A partir de um arquivo existente
handler := logger.NewFileHandler(file)

// Ou criando um arquivo temporário automaticamente
handler, err := logger.NewTemporaryFileHandler("minha-app")
```

**Formato de saída** (JSON Lines):

```json
{"app":"mj-cli","version":"1.0.0","level":"INFO","time":"2025-01-15T10:30:00Z","message":"aplicação iniciada","key":"value"}
```

**Métodos principais:**

| Método | Descrição |
|--------|-----------|
| `Enabled(ctx, level)` | Verifica se o nível está habilitado |
| `Handle(ctx, record)` | Serializa e escreve o registro no arquivo |
| `WithAttrs(attrs)` | Retorna um handler clone com atributos adicionais |
| `WithGroup(name)` | Retorna um handler clone com prefixo de grupo (`name.`) |
| `SetLevel(level)` | Altera o nível mínimo de log |
| `Name()` | Retorna o caminho do arquivo de log |
| `Close()` | Fecha o arquivo subjacente |

Handlers clonados (via `WithAttrs`/`WithGroup`) compartilham o mesmo arquivo e mutex, garantindo segurança em escritas concorrentes.

### MultiHandler

Handler que distribui (fan-out) cada registro para múltiplos handlers nomeados. Um registro só é encaminhado para um handler se ele estiver habilitado para o nível do registro.

```go
multi := logger.NewMultiHandler()
multi.AddHandler("file", fileHandler)
multi.AddHandler("outro", outroHandler)

log := slog.New(multi)
log.Info("enviado para todos os handlers habilitados")
```

**Métodos principais:**

| Método | Descrição |
|--------|-----------|
| `AddHandler(name, handler)` | Registra um handler pelo nome (substitui se já existir) |
| `Enabled(ctx, level)` | `true` se pelo menos um handler aceita o nível |
| `Handle(ctx, record)` | Encaminha para todos os handlers habilitados |
| `WithAttrs(attrs)` | Aplica atributos a todos os handlers registrados |
| `WithGroup(name)` | Aplica grupo a todos os handlers registrados |

## Variável de ambiente DEBUG

A variável `DEBUG` controla o nível de log do `FileHandler`:

| Valor | Comportamento |
|-------|---------------|
| *(não definida ou vazia)* | Nível padrão: `Debug` |
| `debug` | Nível: `Debug` |
| `warn` | Nível: `Warn` |
| `error` | Nível: `Error` |
| *(valor inválido)* | Fallback para `Debug` |

```bash
# Habilitar logs de debug
DEBUG=debug mj-cli

# Apenas warnings e erros
DEBUG=warn mj-cli
```

## Thread Safety

Todos os componentes são thread-safe:

- `FileHandler` usa `*sync.RWMutex` para proteger escritas no arquivo e acesso ao nível; handlers clonados (via `WithAttrs`/`WithGroup`) compartilham o mesmo mutex, garantindo serialização de escritas
- `MultiHandler` usa `*sync.RWMutex` para proteger o mapa de handlers; handlers clonados recebem um mutex próprio, pois cada clone opera sobre sua cópia independente do mapa

## Testes

```bash
go test ./pkg/logger/ -v
```

Os testes cobrem:

- Criação de handlers e logger completo
- Serialização JSON e atributos
- Níveis de log e filtragem
- `WithAttrs` e `WithGroup` (incluindo aninhamento)
- Override via `DEBUG`
- Acesso concorrente (50 goroutines simultâneas)
- Integração com `slog.New()`
