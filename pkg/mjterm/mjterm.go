package mjterm

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	STATUS_CREATED = 1 + iota
	STATUS_RUNNING
	STATUS_SPINNER_RUNNING
	STATUS_SPINNER_PAUSED
	STATUS_SPINNER_STOPPED
	STATUS_PAUSED
	STATUS_STOPPING
	STATUS_STOPPED
)

type Terminal struct {
	priorityOutput  chan string
	output          chan string
	pause           chan struct{}
	pauseAck        chan struct{}
	resume          chan struct{}
	resumeAck       chan struct{}
	close           chan struct{}
	closeAck        chan struct{}
	closeSpinner    chan struct{}
	closeSpinnerAck chan struct{}
	stopOnce        sync.Once
	status          atomic.Int32
	lockClose       sync.RWMutex
	lockInput       sync.Mutex

	// spinner related
	frames  []string
	message string
}

func New() *Terminal {
	term := &Terminal{
		priorityOutput:  make(chan string),
		output:          make(chan string),
		pause:           make(chan struct{}),
		pauseAck:        make(chan struct{}),
		resume:          make(chan struct{}),
		resumeAck:       make(chan struct{}),
		close:           make(chan struct{}),
		closeAck:        make(chan struct{}),
		closeSpinner:    make(chan struct{}),
		closeSpinnerAck: make(chan struct{}),
		status:          atomic.Int32{},
		lockClose:       sync.RWMutex{},
		lockInput:       sync.Mutex{},

		// spinner related
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		message: "Teste loading",
	}
	term.status.Store(STATUS_CREATED)

	go func() {
		term.status.Store(STATUS_RUNNING)
		for {
			select {
			case <-term.pause:
				term.status.Store(STATUS_PAUSED)
				term.pauseAck <- struct{}{}

				closedDuringPause := false
				exitLoop := false
				for !exitLoop {
					select {
					case <-term.resume:
						exitLoop = true
					case msg := <-term.priorityOutput:
						fmt.Print(msg)
					case <-term.close:
						closedDuringPause = true
						exitLoop = true
					}
				}

				if closedDuringPause {
					goto drain
				}
				term.status.Store(STATUS_RUNNING)
				term.resumeAck <- struct{}{}

			case msg := <-term.priorityOutput:
				fmt.Print(msg)
			case msg := <-term.output:
				if term.status.Load() == STATUS_SPINNER_RUNNING {
					fmt.Print("\r\033[K")
				}
				fmt.Print(msg)
			case <-term.close:
				goto drain
			}
		}

	drain:
		term.status.Store(STATUS_STOPPING)
		for {
			select {
			case msg := <-term.priorityOutput:
				fmt.Print(msg)
			case msg := <-term.output:
				fmt.Print(msg)
			default:
				term.closeAck <- struct{}{}
				term.status.Store(STATUS_STOPPED)
				return
			}
		}
	}()

	return term
}

func (t *Terminal) Stop() {
	t.lockClose.Lock()
	defer t.lockClose.Unlock()
	t.stopOnce.Do(func() {
		close(t.close)
		<-t.closeAck
		close(t.priorityOutput)
		close(t.output)
		close(t.pause)
		close(t.pauseAck)
		close(t.resume)
		close(t.resumeAck)
		close(t.closeAck)
	})
}

func (t *Terminal) Status() int32 {
	return t.status.Load()
}

func (t *Terminal) isClosed() bool {
	status := t.status.Load()
	return status == STATUS_STOPPED || status == STATUS_STOPPING
}

func (t *Terminal) pauseOutput() error {
	select {
	case t.pause <- struct{}{}:
		select {
		case <-t.pauseAck:
			return nil
		case <-t.close:
			return ErrTerminalClosed
		}
	case <-t.close:
		return ErrTerminalClosed
	}
}

func (t *Terminal) resumeOutput() error {
	select {
	case t.resume <- struct{}{}:
		select {
		case <-t.resumeAck:
			return nil
		case <-t.close:
			return ErrTerminalClosed
		}
	case <-t.close:
		return ErrTerminalClosed
	}
}
