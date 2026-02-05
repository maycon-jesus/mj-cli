# tint

Pacote para estilização de texto em terminais utilizando códigos de escape ANSI. Detecta automaticamente o nível de suporte a cores do terminal e respeita o padrão [NO_COLOR](https://no-color.org).

## Uso básico

```go
style := tint.NewStyle().Bold().Foreground(tint.Red)
fmt.Println(style.Render("Texto em vermelho e negrito"))
```

## Cores pré-definidas

| Básicas | Brilhantes        |
|---------|-------------------|
| Black   | BrightBlack       |
| Red     | BrightRed         |
| Green   | BrightGreen       |
| Yellow  | BrightYellow      |
| Blue    | BrightBlue        |
| Magenta | BrightMagenta     |
| Cyan    | BrightCyan        |
| White   | BrightWhite       |

## Cores customizadas

Três tipos de cor são disponíveis, correspondentes aos níveis de suporte do terminal:

```go
// TrueColor — 24-bit RGB (16.7 milhões de cores)
cor := tint.TrueColor{R: 255, G: 128, B: 0}

// ColorAnsi256 — paleta de 256 cores
cor := tint.ColorAnsi256{Code: 208}

// ColorAnsi — 16 cores básicas (4-bit)
cor := tint.ColorAnsi{Code: 31} // vermelho
```

### CompleteColor — degradação graceful

`CompleteColor` aceita as três variantes e seleciona automaticamente a mais capaz que o terminal suporta:

```go
cor := tint.CompleteColor{
    TrueColor: tint.TrueColor{R: 255, G: 128, B: 0},
    Ansi256:   tint.ColorAnsi256{Code: 208},
    Ansi:      tint.ColorAnsi{Code: 33},
}

style := tint.NewStyle().Foreground(cor)
fmt.Println(style.Render("cor adaptada ao terminal"))
```

## Estilos de texto

Todos os métodos retornam `*Style` e podem ser encadeados livemente:

| Método          | Efeito                                  |
|-----------------|-----------------------------------------|
| `Bold()`        | Negrito                                 |
| `Dim()`         | Brilho reduzido                         |
| `Italic()`      | Itálico                                 |
| `Underline()`   | Sublinhado                              |
| `Reverse()`     | Inverte cores de fundo e primeiro plano |
| `Hidden()`      | Texto oculto                            |
| `Strikethrough()`| Texto riscado                          |

```go
tint.NewStyle().
    Bold().
    Underline().
    Foreground(tint.Cyan).
    Background(tint.BrightBlack).
    Render("texto estilizado")
```

## Templates pré-definidos

Funções de atalho que combinam formato e cor em um único passo:

| Função             | Formato            | Cor     |
|--------------------|--------------------|---------|
| `PresetTitle`      | `=== mensagem ===` | Cyan    |
| `PresetSubtitle`   | `--- mensagem ---` | White   |
| `PresetSuccess`    | `✔ mensagem`       | Green   |
| `PresetError`      | `✖ mensagem`       | Red     |
| `PresetWarning`    | `⚠ mensagem`       | Yellow  |
| `PresetInfo`       | `ℹ mensagem`       | Blue    |
| `PresetLambda`     | `λ mensagem`       | Magenta |

```go
fmt.Println(tint.PresetSuccess("operação concluída"))
fmt.Println(tint.PresetError("algo deu errado"))
```

## Detecção de suporte a cores

O nível de suporte é detectado automaticamente no arranque do programa, com base nas variáveis de ambiente `NO_COLOR`, `COLORTERM` e `TERM`. O valor detectado fica exposto em `tint.DetectedColorSupport`:

| Nível              | Cores               | Condição de detecção                          |
|--------------------|---------------------|-----------------------------------------------|
| `ColorSupportNone` | nenhuma             | `NO_COLOR` definido                           |
| `ColorSupport4bit` | 16                  | fallback padrão                               |
| `ColorSupport8bit` | 256                 | `TERM` contém `256color` ou similar           |
| `ColorSupport24bit`| 16.7 milhões (RGB)  | `COLORTERM=truecolor` ou `TERM` com `24bit`   |

## NO_COLOR

Quando a variável de ambiente `NO_COLOR` está definida (com qualquer valor), todas as saídas de ANSI são suprimidas e `Render` retorna o texto plano:

```bash
NO_COLOR=1 ./meu-programa
```

## Thread safety

Todas as operações em `Style` são seguras para uso concorrente, protegidas por um `sync.RWMutex` interno.
