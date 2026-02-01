package mjterm

import (
	"time"
)

func (t *Terminal) StartSpinner() {
	go func() {
		t.status.Store(STATUS_SPINNER_RUNNING)
		i := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				frame := t.frames[i%len(t.frames)]
				t.Printf("\033[36m%s\033[0m %s", frame, t.message)
				i++
			case <-t.closeSpinner:
				t.Printf("") // clear line
				t.closeSpinnerAck <- struct{}{}
				return
			}
		}
	}()
}

func (t *Terminal) StopSpinner() {
	t.closeSpinner <- struct{}{}
	<-t.closeSpinnerAck
	t.status.Store(STATUS_RUNNING)
	t.Printf("\r\033[K\033[32m✓\033[0m %s\n", t.message)
}

func (t *Terminal) StopSpinnerWithError() {
	t.closeSpinner <- struct{}{}
	<-t.closeSpinnerAck
	t.status.Store(STATUS_RUNNING)
	t.Printf("\r\033[K\033[31m✗\033[0m %s\n", t.message)
}
