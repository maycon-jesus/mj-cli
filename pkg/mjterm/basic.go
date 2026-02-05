package mjterm

import "fmt"

func (t *Terminal) Print(msg string) error {
	return t.addEvent(printMsg{content: msg})
}

func (t *Terminal) Println(msg string) error {
	return t.addEvent(printMsg{content: msg + "\n"})
}

func (t *Terminal) Printf(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return t.addEvent(printMsg{content: msg})
}

// Write implementa io.Writer para o Terminal, permitindo que ele seja usado
// como destino de saída em funções que aceitam um Writer.
func (t *Terminal) Write(data []byte) (int, error) {
	err := t.addEvent(printMsg{content: string(data)})
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (t *Terminal) NewLine() error {
	return t.addEvent(printMsg{content: "\n"})
}
