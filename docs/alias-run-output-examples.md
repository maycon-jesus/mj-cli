# Exemplos de Output Bonito para `mj-cli alias run`

Este documento mostra exemplos de como melhorar o output do comando `alias run` usando o pacote `beautyoutput`.

## Uso Básico do StrBuilder

```go
import beautyoutput "github.com/maycon-jesus/mj-cli/pkg/beauty_output"

// Criar um novo builder
sb := beautyoutput.NewStrBuilder()
```

## Exemplos de Output

### 1. Header com Nome do Alias

```go
func printAliasHeader(aliasName string) {
    sb := beautyoutput.NewStrBuilder()
    output := sb.Bold().Cyan().Text("▶ Running alias: ").
        Reset().Bold().White().Text(aliasName).
        String()
    fmt.Println(output)
}
```

**Resultado visual:**
```
▶ Running alias: my-alias
```

---

### 2. Mostrar o Comando Sendo Executado

```go
func printCommand(command string) {
    sb := beautyoutput.NewStrBuilder()

    // Linha de separação
    fmt.Println(sb.Dim().Text("─────────────────────────────────").String())

    sb = beautyoutput.NewStrBuilder()
    output := sb.BrightBlack().Text("$ ").
        Reset().Yellow().Text(command).
        String()
    fmt.Println(output)

    sb = beautyoutput.NewStrBuilder()
    fmt.Println(sb.Dim().Text("─────────────────────────────────").String())
}
```

**Resultado visual:**
```
─────────────────────────────────
$ git status
─────────────────────────────────
```

---

### 3. Status de Sucesso/Erro

```go
func printSuccess(aliasName string) {
    sb := beautyoutput.NewStrBuilder()
    output := sb.Bold().Green().Text("✓ ").
        Reset().Green().Text("Alias '").
        Bold().Text(aliasName).
        Reset().Green().Text("' executed successfully").
        String()
    fmt.Println(output)
}

func printError(aliasName string, err error) {
    sb := beautyoutput.NewStrBuilder()
    output := sb.Bold().Red().Text("✗ ").
        Reset().Red().Text("Failed to run alias '").
        Bold().Text(aliasName).
        Reset().Red().Text("': ").
        Dim().Text(err.Error()).
        String()
    fmt.Println(output)
}
```

**Resultado visual:**
```
✓ Alias 'my-alias' executed successfully
✗ Failed to run alias 'my-alias': command not found
```

---

### 4. Exibir Variáveis Substituídas

```go
func printVariables(variables map[string]string) {
    if len(variables) == 0 {
        return
    }

    sb := beautyoutput.NewStrBuilder()
    fmt.Println(sb.Dim().Italic().Text("Variables:").String())

    for key, value := range variables {
        sb = beautyoutput.NewStrBuilder()
        output := sb.BrightBlack().Text("  $").
            Cyan().Text(key).
            BrightBlack().Text(" → ").
            White().Text(value).
            String()
        fmt.Println(output)
    }
    fmt.Println()
}
```

**Resultado visual:**
```
Variables:
  $1 → main
  $2 → feature-branch
```

---

### 5. Exemplo Completo Integrado

```go
func runAliasWithBeautifulOutput(aliasName, command string, variables map[string]string) error {
    // Header
    sb := beautyoutput.NewStrBuilder()
    fmt.Println()
    fmt.Println(sb.Bold().Cyan().Text("▶ Running alias: ").
        Reset().Bold().White().Text(aliasName).String())

    // Mostrar variáveis se houver
    if len(variables) > 0 {
        sb = beautyoutput.NewStrBuilder()
        fmt.Println(sb.Dim().Italic().Text("  Variables:").String())
        for key, value := range variables {
            sb = beautyoutput.NewStrBuilder()
            fmt.Println(sb.BrightBlack().Text("    $").
                Cyan().Text(key).
                BrightBlack().Text(" → ").
                White().Text(value).String())
        }
    }

    // Comando
    sb = beautyoutput.NewStrBuilder()
    fmt.Println()
    fmt.Println(sb.Dim().Text("  Command: ").
        Reset().Yellow().Text(command).String())

    // Separador
    sb = beautyoutput.NewStrBuilder()
    fmt.Println(sb.Dim().Text("  ─────────────────────────────────").String())
    fmt.Println()

    // Executar comando
    err := cmd.RunCommand(command)

    fmt.Println()

    // Status final
    if err != nil {
        sb = beautyoutput.NewStrBuilder()
        fmt.Println(sb.Bold().Red().Text("  ✗ ").
            Reset().Red().Text("Command failed").String())
        return err
    }

    sb = beautyoutput.NewStrBuilder()
    fmt.Println(sb.Bold().Green().Text("  ✓ ").
        Reset().Green().Text("Done").String())
    fmt.Println()

    return nil
}
```

**Resultado visual:**
```
▶ Running alias: deploy
  Variables:
    $1 → production
    $2 → v1.2.3

  Command: ./deploy.sh production v1.2.3
  ─────────────────────────────────

  [output do comando aqui]

  ✓ Done
```

---

### 6. Estilo Minimalista

```go
func runAliasMinimal(aliasName, command string) {
    sb := beautyoutput.NewStrBuilder()
    fmt.Println(sb.Dim().Text("$ ").Reset().Text(command).String())
}
```

**Resultado visual:**
```
$ git push origin main
```

---

### 7. Box Style com Bordas

```go
func printAliasBox(aliasName, command string) {
    sb := beautyoutput.NewStrBuilder()

    // Borda superior
    fmt.Println(sb.Cyan().Text("╭─────────────────────────────────╮").String())

    // Título
    sb = beautyoutput.NewStrBuilder()
    fmt.Printf("%s  %s%s\n",
        sb.Cyan().Text("│").String(),
        beautyoutput.NewStrBuilder().Bold().White().Text(aliasName).String(),
        beautyoutput.NewStrBuilder().Cyan().Text("                          │").String())

    // Separador
    sb = beautyoutput.NewStrBuilder()
    fmt.Println(sb.Cyan().Text("├─────────────────────────────────┤").String())

    // Comando
    sb = beautyoutput.NewStrBuilder()
    fmt.Printf("%s  %s\n",
        sb.Cyan().Text("│").String(),
        beautyoutput.NewStrBuilder().Yellow().Text(command).String())

    // Borda inferior
    sb = beautyoutput.NewStrBuilder()
    fmt.Println(sb.Cyan().Text("╰─────────────────────────────────╯").String())
}
```

**Resultado visual:**
```
╭─────────────────────────────────╮
│  deploy                          │
├─────────────────────────────────┤
│  ./deploy.sh $1 $2
╰─────────────────────────────────╯
```

---

### 8. Usando Cores RGB Customizadas

```go
func printWithCustomColors(text string) {
    sb := beautyoutput.NewStrBuilder()

    // Cor laranja personalizada
    output := sb.RGB(255, 165, 0).Bold().Text("⚡ ").
        Reset().RGB(255, 200, 100).Text(text).
        String()

    fmt.Println(output)
}
```

---

## Símbolos Úteis para Output

| Símbolo | Uso |
|---------|-----|
| `▶` | Executando |
| `✓` | Sucesso |
| `✗` | Erro |
| `⚡` | Rápido/Ação |
| `⚠` | Aviso |
| `ℹ` | Informação |
| `→` | Seta |
| `•` | Item de lista |
| `─` | Linha horizontal |
| `│` | Linha vertical |
| `╭╮╰╯` | Cantos de box |

## Dicas

1. **Sempre use `Reset()`** antes de mudar de estilo para evitar que estilos anteriores afetem o novo texto.

2. **Combine estilos**: Você pode encadear `.Bold().Cyan().Underline()` para combinar múltiplos estilos.

3. **Mantenha consistência**: Use o mesmo padrão de cores em todo o CLI:
   - Cyan: Headers/títulos
   - Yellow: Comandos
   - Green: Sucesso
   - Red: Erros
   - Dim/BrightBlack: Informações secundárias

4. **Menos é mais**: Não exagere nas cores. Use-as para destacar informações importantes.
