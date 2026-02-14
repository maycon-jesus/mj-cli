// Package cmd fornece funcionalidades para executar comandos de sistema
// de forma simplificada, com suporte para parsing de argumentos complexos,
// incluindo strings entre aspas e caracteres de escape.
//
// O pacote oferece três modos principais de execução:
//   - RunCommand: execução simples com entrada/saída padrão
//   - RunCommandWithOptions: execução com controle customizado de I/O
//   - GetCommandOutput: captura da saída do comando
//
// Exemplo de uso básico:
//
//	err := cmd.RunCommand("ls -la /tmp")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Exemplo com opções customizadas:
//
//	options := cmd.CommandOptions{
//		Stdin:  strings.NewReader("input"),
//		Stdout: &buf,
//		Stderr: os.Stderr,
//	}
//	err := cmd.RunCommandWithOptions("grep pattern", options)
package cmd

import (
	"io"
	"os"
	"os/exec"
)

// RunCommand executa um comando de sistema usando os streams padrão
// (stdin, stdout, stderr) do processo atual.
//
// O parâmetro cmdStr é parseado automaticamente para separar o comando
// e seus argumentos, com suporte para aspas e caracteres de escape.
//
// Exemplo:
//
//	err := RunCommand("echo 'Hello World'")
//	err := RunCommand("git commit -m \"mensagem\"")
//
// Retorna erro se o parsing falhar ou se o comando não executar com sucesso.
func RunCommand(cmdStr string) error {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// CommandOptions define as opções de I/O para execução de comandos.
//
// Permite customizar os streams de entrada e saída do comando,
// útil para testes, redirecionamento ou captura de dados.
//
// Exemplo:
//
//	var buf bytes.Buffer
//	options := CommandOptions{
//		Stdin:  strings.NewReader("input data\n"),
//		Stdout: &buf,
//		Stderr: os.Stderr,
//	}
type CommandOptions struct {
	// Stdin define o stream de entrada do comando
	Stdin io.Reader
	// Stdout define o stream de saída padrão do comando
	Stdout io.Writer
	// Stderr define o stream de saída de erro do comando
	Stderr io.Writer
}

// RunCommandWithOptions executa um comando com controle customizado
// dos streams de entrada e saída.
//
// Esta função é útil quando você precisa redirecionar stdin, stdout
// ou stderr para buffers, arquivos ou outros destinos customizados.
//
// O parâmetro cmdStr é parseado automaticamente para separar o comando
// e seus argumentos, com suporte para aspas e caracteres de escape.
//
// Exemplo:
//
//	var output bytes.Buffer
//	options := CommandOptions{
//		Stdin:  os.Stdin,
//		Stdout: &output,
//		Stderr: os.Stderr,
//	}
//	err := RunCommandWithOptions("ls -la", options)
//	fmt.Println(output.String())
//
// Retorna erro se o parsing falhar ou se o comando não executar com sucesso.
func RunCommandWithOptions(cmdStr string, options CommandOptions) error {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = options.Stdin
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// GetCommandOutput executa um comando e retorna sua saída combinada
// (stdout + stderr) como uma string.
//
// Esta função é útil quando você precisa capturar e processar a saída
// de um comando sem interação com o usuário.
//
// O parâmetro cmdStr é parseado automaticamente para separar o comando
// e seus argumentos, com suporte para aspas e caracteres de escape.
//
// Exemplo:
//
//	output, err := GetCommandOutput("git status")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(output)
//
// Retorna a saída do comando como string e um erro se o parsing falhar
// ou se o comando não executar com sucesso.
func GetCommandOutput(cmdStr string) (string, error) {
	args, err := parseCommand(cmdStr)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
