package mjterm

import "fmt"

func (t *Terminal) Print(msg string) error {
	return t.addEvent(printMsg{content: msg})
}

func (t *Terminal) Printf(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return t.addEvent(printMsg{content: msg})
}

func (t *Terminal) NewLine() error {
	return t.addEvent(printMsg{content: "\n"})
}
