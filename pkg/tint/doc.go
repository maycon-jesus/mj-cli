// Package tint fornece funcionalidades para estilização de texto em terminais
// utilizando códigos de escape ANSI.
//
// O pacote detecta automaticamente o nível de suporte a cores do terminal,
// podendo trabalhar com 16 cores (4-bit), 256 cores (8-bit) ou 16.7 milhões
// de cores (24-bit/truecolor).
//
// # Uso Básico
//
// Para estilizar texto, crie um novo Style e encadeie os métodos desejados:
//
//	style := tint.NewStyle().Bold().Foreground(tint.Red)
//	fmt.Println(style.Render("Texto em vermelho e negrito"))
//
// # Cores Disponíveis
//
// O pacote oferece cores ANSI pré-definidas:
//
//   - Cores básicas: Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
//   - Cores brilhantes: BrightBlack, BrightRed, BrightGreen, BrightYellow,
//     BrightBlue, BrightMagenta, BrightCyan, BrightWhite
//
// # Cores Customizadas
//
// Para cores customizadas, utilize os tipos disponíveis:
//
//	// TrueColor (RGB)
//	cor := tint.TrueColor{R: 255, G: 128, B: 0}
//
//	// Cor da paleta de 256 cores
//	cor256 := tint.ColorAnsi256{Code: 208}
//
//	// Cor ANSI básica (código direto)
//	corAnsi := tint.ColorAnsi{Code: 31}
//
// # Estilos de Texto
//
// Os seguintes estilos de texto estão disponíveis:
//
//   - Bold: texto em negrito
//   - Dim: texto com brilho reduzido
//   - Italic: texto em itálico
//   - Underline: texto sublinhado
//   - Reverse: inverte cores de fundo e primeiro plano
//   - Hidden: texto oculto
//   - Strikethrough: texto riscado
//
// # Thread Safety
//
// Todas as operações em [Style] são seguras para uso concorrente,
// protegidas por mutex interno.
//
// # Suporte a NO_COLOR
//
// O pacote respeita o padrão NO_COLOR (https://no-color.org).
// Quando a variável de ambiente NO_COLOR está definida (com qualquer valor),
// todas as saídas de cores e estilos ANSI são desabilitadas automaticamente.
//
// Isso é útil para:
//   - Acessibilidade (usuários com leitores de tela)
//   - Redirecionamento de saída para arquivos
//   - Ambientes onde códigos ANSI não são suportados
//   - Preferência pessoal do usuário
//
// Exemplo de uso com NO_COLOR:
//
//	# No terminal:
//	NO_COLOR=1 ./meu-programa
//
//	# No código, o comportamento é automático:
//	style := tint.NewStyle().Bold().Foreground(tint.Red)
//	fmt.Println(style.Render("Este texto será sem formatação se NO_COLOR estiver definido"))
//
// # Detecção de Suporte a Cores
//
// O pacote detecta automaticamente o nível de suporte do terminal através
// das variáveis de ambiente COLORTERM e TERM:
//
//   - [ColorSupportNone]: NO_COLOR definido, nenhum código ANSI emitido
//   - [ColorSupport4bit]: 16 cores básicas (fallback padrão)
//   - [ColorSupport8bit]: 256 cores (terminais com "256color" ou modernos)
//   - [ColorSupport24bit]: 16.7M cores RGB (terminais com truecolor)
//
// O nível detectado está disponível em [DetectedColorSupport].
//
// # CompleteColor para Degradação Graceful
//
// Use [CompleteColor] para definir cores com fallbacks automáticos:
//
//	cor := tint.CompleteColor{
//		TrueColor: tint.TrueColor{R: 255, G: 128, B: 0},  // Terminal 24-bit
//		Ansi256:   tint.ColorAnsi256{Code: 208},          // Terminal 256 cores
//		Ansi:      tint.ColorAnsi{Code: 33},              // Terminal básico
//	}
//	// O pacote escolhe automaticamente o formato adequado ao terminal
package tint
