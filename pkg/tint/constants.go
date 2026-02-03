package tint

// Constantes de nível de suporte a cores do terminal.
const (
	ColorSupportNone  ColorSupport = iota // Sem cores (NO_COLOR definido)
	ColorSupport4bit                      // 16 cores (cores ANSI básicas)
	ColorSupport8bit                      // 256 cores (paleta extendida)
	ColorSupport24bit                     // 16.7M cores (RGB/truecolor)
)

// Códigos de estilo ANSI SGR (Select Graphic Rendition).
// Estes códigos são usados para aplicar formatação ao texto.
const (
	STYLE_RESET         = 0 + iota // Reseta todos os atributos
	STYLE_BOLD                     // Texto em negrito
	STYLE_DIM                      // Texto com brilho reduzido
	STYLE_ITALIC                   // Texto em itálico
	STYLE_UNDERLINE                // Texto sublinhado
	_                              // Reservado (blink lento)
	_                              // Reservado (blink rápido)
	STYLE_REVERSE                  // Inverte cores de fundo e texto
	STYLE_HIDDEN                   // Texto oculto
	STYLE_STRIKETHROUGH            // Texto riscado
)

// Códigos de cores ANSI básicas (primeiro plano).
// Para cores de fundo, adicione 10 a estes valores.
const (
	ANSI_COLOR_BLACK   = 30 + iota // Preto
	ANSI_COLOR_RED                 // Vermelho
	ANSI_COLOR_GREEN               // Verde
	ANSI_COLOR_YELLOW              // Amarelo
	ANSI_COLOR_BLUE                // Azul
	ANSI_COLOR_MAGENTA             // Magenta
	ANSI_COLOR_CYAN                // Ciano
	ANSI_COLOR_WHITE               // Branco
)

// ANSI_COLOR_BRIGHT_OFFSET é o offset adicionado aos códigos de cor
// para obter as versões brilhantes (90-97 para primeiro plano).
const (
	ANSI_COLOR_BRIGHT_OFFSET = 60
)
