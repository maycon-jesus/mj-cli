# Package cmd

O pacote `cmd` fornece funcionalidades para executar comandos de sistema de forma simplificada, com suporte para parsing de argumentos complexos, incluindo strings entre aspas e caracteres de escape.

## Características

- **Parser Inteligente**: Processa automaticamente comandos com aspas duplas, aspas simples e caracteres de escape
- **Modos de Execução Flexíveis**: Três formas diferentes de executar comandos
- **Controle de I/O**: Redirecionamento customizado de stdin, stdout e stderr
- **Compatível com Bash**: Comportamento similar ao shell bash para parsing de comandos
- **Testado**: Cobertura completa de testes unitários e benchmarks

## Instalação

```go
import "github.com/maycon-jesus/mj-cli/pkg/cmd"
```

## Uso Básico

### RunCommand - Execução Simples

Executa um comando usando os streams padrão (stdin, stdout, stderr):

```go
err := cmd.RunCommand("ls -la /tmp")
if err != nil {
    log.Fatal(err)
}
```

### RunCommandWithOptions - Execução Customizada

Executa um comando com controle total sobre os streams de I/O:

```go
var output bytes.Buffer
options := cmd.CommandOptions{
    Stdin:  strings.NewReader("input data"),
    Stdout: &output,
    Stderr: os.Stderr,
}

err := cmd.RunCommandWithOptions("grep pattern", options)
if err != nil {
    log.Fatal(err)
}

fmt.Println(output.String())
```

### GetCommandOutput - Captura de Saída

Executa um comando e captura sua saída (stdout + stderr):

```go
output, err := cmd.GetCommandOutput("git status")
if err != nil {
    log.Fatal(err)
}
fmt.Println(output)
```

## API Reference

### Funções

#### RunCommand

```go
func RunCommand(cmdStr string) error
```

Executa um comando de sistema usando os streams padrão do processo atual.

**Parâmetros:**
- `cmdStr`: String contendo o comando e seus argumentos

**Retorna:**
- `error`: Erro se o parsing falhar ou se o comando não executar com sucesso

**Exemplo:**
```go
err := cmd.RunCommand("echo 'Hello World'")
err := cmd.RunCommand("git commit -m \"mensagem\"")
```

#### RunCommandWithOptions

```go
func RunCommandWithOptions(cmdStr string, options CommandOptions) error
```

Executa um comando com controle customizado dos streams de entrada e saída.

**Parâmetros:**
- `cmdStr`: String contendo o comando e seus argumentos
- `options`: Estrutura CommandOptions com configurações de I/O

**Retorna:**
- `error`: Erro se o parsing falhar ou se o comando não executar com sucesso

**Exemplo:**
```go
var output bytes.Buffer
options := cmd.CommandOptions{
    Stdin:  os.Stdin,
    Stdout: &output,
    Stderr: os.Stderr,
}
err := cmd.RunCommandWithOptions("ls -la", options)
```

#### GetCommandOutput

```go
func GetCommandOutput(cmdStr string) (string, error)
```

Executa um comando e retorna sua saída combinada (stdout + stderr) como string.

**Parâmetros:**
- `cmdStr`: String contendo o comando e seus argumentos

**Retorna:**
- `string`: Saída combinada do comando
- `error`: Erro se o parsing falhar ou se o comando não executar com sucesso

**Exemplo:**
```go
output, err := cmd.GetCommandOutput("git log --oneline -n 5")
if err != nil {
    log.Fatal(err)
}
fmt.Println(output)
```

### Tipos

#### CommandOptions

```go
type CommandOptions struct {
    Stdin  io.Reader
    Stdout io.Writer
    Stderr io.Writer
}
```

Define as opções de I/O para execução de comandos.

**Campos:**
- `Stdin`: Stream de entrada do comando
- `Stdout`: Stream de saída padrão do comando
- `Stderr`: Stream de saída de erro do comando

## Parser de Comandos

O pacote inclui um parser interno que processa a string do comando de forma inteligente:

### Aspas Duplas (`"`)

Permitem agrupar argumentos e suportam sequências de escape:

```go
cmd.RunCommand(`echo "Hello World"`)           // Argumento: Hello World
cmd.RunCommand(`echo "Line 1\nLine 2"`)        // Preserva \n literalmente
cmd.RunCommand(`git commit -m "fix: bug"`)     // Argumento: fix: bug
```

**Escapes válidos dentro de aspas duplas:**
- `\"` - Aspas duplas
- `\\` - Barra invertida
- `\$` - Cifrão
- `` \` `` - Crase

### Aspas Simples (`'`)

Tratam tudo literalmente, sem processamento de escape (como no bash):

```go
cmd.RunCommand(`echo 'Hello World'`)           // Argumento: Hello World
cmd.RunCommand(`echo '$HOME'`)                 // Argumento: $HOME (literal)
cmd.RunCommand(`echo 'test\ntest'`)           // Argumento: test\ntest (literal)
```

### Caracteres de Escape

Fora de aspas, a barra invertida (`\`) torna o próximo caractere literal:

```go
cmd.RunCommand(`echo hello\ world`)            // Argumento: hello world
cmd.RunCommand(`echo test\*`)                  // Argumento: test*
```

### Exemplos de Parsing

```go
// Comando simples
"ls -la"                    → ["ls", "-la"]

// Com aspas duplas
"echo \"hello world\""      → ["echo", "hello world"]

// Com aspas simples
"echo 'hello world'"        → ["echo", "hello world"]

// Caracteres especiais
"echo '$HOME'"              → ["echo", "$HOME"]

// Escape fora de aspas
"echo hello\\ world"        → ["echo", "hello world"]

// Aspas aninhadas
"echo \"it's ok\""          → ["echo", "it's ok"]
```

### Erros de Parsing

O parser retorna erro nas seguintes situações:

```go
// Aspas não fechadas
cmd.RunCommand(`echo "unclosed`)  // Error: unclosed quote: "

// Barra invertida no final
cmd.RunCommand(`echo test\`)      // Error: unexpected end of command: trailing backslash
```

## Casos de Uso

### Executar comandos Git

```go
// Commit com mensagem
err := cmd.RunCommand(`git commit -m "feat: new feature"`)

// Ver status
output, err := cmd.GetCommandOutput("git status --short")
```

### Processar saída de comandos

```go
var output bytes.Buffer
options := cmd.CommandOptions{
    Stdout: &output,
    Stderr: &output,
}

err := cmd.RunCommandWithOptions("npm install", options)
if err != nil {
    log.Printf("Erro na instalação: %s", output.String())
}
```

### Redirecionar entrada

```go
input := strings.NewReader("search pattern\n")
var output bytes.Buffer

options := cmd.CommandOptions{
    Stdin:  input,
    Stdout: &output,
    Stderr: os.Stderr,
}

err := cmd.RunCommandWithOptions("grep -i pattern", options)
```

### Capturar logs de aplicação

```go
logs, err := cmd.GetCommandOutput("docker logs my-container --tail 100")
if err != nil {
    log.Fatal(err)
}

// Processar logs
for _, line := range strings.Split(logs, "\n") {
    fmt.Println(line)
}
```

## Testes

O pacote inclui testes abrangentes cobrindo:

- Execução de comandos simples e complexos
- Parsing de aspas e caracteres de escape
- Redirecionamento de I/O
- Captura de stdout e stderr
- Tratamento de erros
- Comandos inexistentes
- Caracteres especiais
- Compatibilidade multi-plataforma

Para executar os testes:

```bash
go test ./pkg/cmd
```

Para executar com cobertura:

```bash
go test ./pkg/cmd -cover
```

Para executar benchmarks:

```bash
go test ./pkg/cmd -bench=.
```

## Observações de Plataforma

Alguns testes são pulados automaticamente no Windows devido a diferenças nos comandos disponíveis. O pacote funciona em todas as plataformas, mas os comandos executados devem estar disponíveis no sistema operacional.

## Tratamento de Erros

O pacote retorna erros em três situações principais:

1. **Erro de Parsing**: Quando a string do comando contém sintaxe inválida
   ```go
   err := cmd.RunCommand(`echo "unclosed`)
   // Error: unclosed quote: "
   ```

2. **Comando Não Encontrado**: Quando o executável não existe
   ```go
   err := cmd.RunCommand("nonexistent-command")
   // Error: executable file not found
   ```

3. **Erro de Execução**: Quando o comando retorna código de saída não-zero
   ```go
   err := cmd.RunCommand("ls /nonexistent-path")
   // Error: exit status 2
   ```

## Boas Práticas

1. **Sempre verifique erros**: Nunca ignore o erro retornado pelas funções
   ```go
   if err := cmd.RunCommand("git push"); err != nil {
       log.Printf("Erro ao fazer push: %v", err)
       return err
   }
   ```

2. **Use aspas para argumentos com espaços**:
   ```go
   cmd.RunCommand(`git commit -m "mensagem com espaços"`)
   ```

3. **Prefira GetCommandOutput para capturar saída**:
   ```go
   output, err := cmd.GetCommandOutput("git status")
   // Melhor que criar buffers manualmente
   ```

4. **Use CommandOptions quando precisar de controle fino**:
   ```go
   // Para testes, redirecionamento específico, etc.
   options := cmd.CommandOptions{
       Stdout: customWriter,
       Stderr: customErrorHandler,
   }
   ```

## Licença

Este pacote faz parte do projeto mj-cli.
