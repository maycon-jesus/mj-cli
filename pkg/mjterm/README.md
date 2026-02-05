# mjterm

Gerenciador de I/O para terminais baseado em loop de eventos. Todas as operações de impressão, entrada e spinner são serializadas em uma única goroutine, eliminando *race conditions* sem exigir sincronização manual por parte do chamador.

## Ciclo de vida

```go
term := mjterm.New()   // inicia o loop de eventos
defer term.Close()     // aguarda eventos pendentes e encerra
```

`Close` é idempotente — pode ser chamado mais de uma vez sem efeitos adversos. Após o fechamento, qualquer operação retorna `mjterm.ErrTerminalClosed`.

> **Nota:** o flag `closed` é marcado apenas depois que todos os eventos pendentes forem processados e o WaitGroup drenado, evitando uma *race condition* entre operações em andamento e a leitura de `closed` no `addEvent`.

## Impressão

```go
term.Print("olá mundo")
term.Println("olá mundo")          // equivalente a Print + "\n"
term.Printf("porta: %d\n", 8080)
term.NewLine()
```

### io.Writer

`Terminal` implementa `io.Writer`, permitindo que seja passado diretamente a qualquer função que aceite esse interface (ex: `log.New`, `json.NewEncoder`).

```go
enc := json.NewEncoder(term)   // saída vai para o loop de eventos
enc.Encode(dados)
```

Se um spinner estiver ativo no momento da impressão, a linha do spinner é limpa automaticamente antes de escrever a mensagem.

## Entrada de texto

```go
// sem timeout
nome, err := term.Prompt("Nome: ")

// com cancelamento via context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
nome, err := term.PromptWithContext(ctx, "Nome: ")
```

Durante a espera de entrada o spinner ativo é pausado automaticamente e retoma após a leitura.

## Spinner

O spinner exibe uma animação em Braille (`⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏`) atualizada a cada 100 ms.

```go
term.StartSpinner("deploy", "enviando pacote...")

// operação longa ...

// sucesso — exibe ✓ em verde
term.StopSpinner("deploy", "pacote enviado")

// ou erro — exibe ✗ em vermelho
term.StopSpinnerWithError("deploy", "falha no envio")
```

Comportamentos automáticos do spinner:

- A linha é limpa antes de qualquer `Print` ou `Prompt`.
- A animação é pausada enquanto há entrada pendente.
- `Close` para o spinner antes de encerrar o loop.

## Arquitetura

```
chamador          canal (buffer 10)       goroutine do loop
   │                     │                       │
   ├─ Print ──────────►  │  ── printMsg ───────► │ fmt.Print
   ├─ Prompt ─────────►  │  ── promptMsg ──────► │ bufio.Scan → result ch
   ├─ StartSpinner ────► │  ── startSpinnerMsg ► │ inicia ticker goroutine
   ├─ StopSpinner ─────► │  ── stopSpinnerMsg ─► │ para ticker
   └─ Close ───────────► │  ── closeMsg ──────►  │ sinaliza encerramento
```

Todas as mensagens são processadas na ordem em que chegam ao canal, garantindo que a saída do terminal seja sempre coerente.
