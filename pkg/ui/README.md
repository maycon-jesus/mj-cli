# ui

Pacote de componentes de alto nível para renderização de mensagens no terminal. Construído sobre o pacote [`tint`](../tint/), fornece funções prontas para formatação consistente de saídas da CLI.

## Uso básico

```go
fmt.Println(ui.Title("Iniciando processo"))
fmt.Println(ui.Success("Operação concluída"))
fmt.Println(ui.Error("Algo deu errado"))
```

## Componentes disponíveis

| Função     | Formato            | Cor     | Exemplo de saída     |
|------------|--------------------|---------|----------------------|
| `Title`    | `=== mensagem ===` | Cyan    | `=== Deploy ===`     |
| `Subtitle` | `--- mensagem ---` | White   | `--- Etapa 1 ---`    |
| `Success`  | `✔ mensagem`       | Green   | `✔ Tudo certo`       |
| `Error`    | `✖ mensagem`       | Red     | `✖ Falha na build`   |
| `Warning`  | `⚠ mensagem`       | Yellow  | `⚠ Disco quase cheio`|
| `Info`     | `ℹ mensagem`       | Blue    | `ℹ 3 arquivos`       |
| `Lambda`   | `λ mensagem`       | Magenta | `λ git pull`         |

Todas as funções aplicam **negrito** automaticamente.

## NO_COLOR

Quando a variável de ambiente `NO_COLOR` está definida, as funções retornam o texto formatado sem códigos ANSI:

```go
// NO_COLOR=1
ui.Success("ok") // "✔ ok" (sem cores)
```
