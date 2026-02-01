package mjterm

import "fmt"

var ErrTerminalClosed = fmt.Errorf("terminal is closed")

func (t *Terminal) Print(msg string) error {
	t.lockClose.RLock()
	closed := t.isClosed()
	t.lockClose.RUnlock()

	if closed {
		return ErrTerminalClosed
	}

	select {
	case t.output <- msg:
		return nil
	case <-t.close:
		return ErrTerminalClosed
	}
}

func (t *Terminal) Printf(format string, a ...any) error {
	t.lockClose.RLock()
	closed := t.isClosed()
	t.lockClose.RUnlock()

	if closed {
		return ErrTerminalClosed
	}

	msg := fmt.Sprintf(format, a...)
	select {
	case t.output <- msg:
		return nil
	case <-t.close:
		return ErrTerminalClosed
	}
}

func (t *Terminal) priorityPrint(msg string) error {
	t.lockClose.RLock()
	closed := t.isClosed()
	t.lockClose.RUnlock()

	if closed {
		return ErrTerminalClosed
	}

	select {
	case t.priorityOutput <- msg:
		return nil
	case <-t.close:
		return ErrTerminalClosed
	}
}

func (t *Terminal) NewLine() error {
	t.lockClose.RLock()
	closed := t.isClosed()
	t.lockClose.RUnlock()

	if closed {
		return ErrTerminalClosed
	}

	select {
	case t.output <- "\n":
		return nil
	case <-t.close:
		return ErrTerminalClosed
	}
}
