package mjterm

func (t *Terminal) StartSpinner(id string, message string) error {
	return t.addEvent(startSpinnerMsg{
		id:      id,
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	})
}

func (t *Terminal) StopSpinner(id string, message string) error {
	return t.addEvent(stopSpinnerMsg{
		id:      id,
		message: message,
	})
}

func (t *Terminal) StopSpinnerWithError(id string, message string) error {
	return t.addEvent(stopSpinnerMsg{
		id:      id,
		message: message,
		err:     true,
	})
}
