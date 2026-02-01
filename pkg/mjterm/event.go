package mjterm

import "fmt"

type msg interface{}

var ErrTerminalClosed = fmt.Errorf("terminal is closed")

func (t *Terminal) addEvent(event msg) error {
	if t.closed.Load() {
		return ErrTerminalClosed
	}
	t.wg.Add(1)
	t.events <- event
	return nil
}

type printMsg struct {
	content string
}

type promptResult struct {
	response string
	err      error
}

type promptMsg struct {
	content string
	result  chan<- promptResult
}

type startSpinnerMsg struct {
	id      string
	message string
	frames  []string
	line    int
}

type stopSpinnerMsg struct {
	id      string
	message string
	err     bool
}

type closeMsg struct {
	closed chan struct{}
}
