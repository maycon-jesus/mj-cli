package beautyoutput

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	frames       []string
	message      string
	stop         chan struct{}
	done         chan struct{}
	messageQueue chan string
	mu           sync.Mutex
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:       []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		message:      message,
		stop:         make(chan struct{}),
		done:         make(chan struct{}),
		messageQueue: make(chan string, 100),
	}
}

func (s *Spinner) Start() {
	fmt.Print("\033[?25l") // esconde cursor
	go func() {
		i := 0
		stopping := false
		for {
			select {
			case msg := <-s.messageQueue:
				s.mu.Lock()
				fmt.Print("\r\033[K")
				str := NewStrBuilder().Dim().Textf("%s", msg).Reset()
				fmt.Print(str)
				s.mu.Unlock()

				if stopping && len(s.messageQueue) == 0 {
					close(s.done)
					return
				}
			case <-s.stop:
				s.mu.Lock()
				stopping = true
				if len(s.messageQueue) == 0 {
					close(s.done)
					return
				}
				s.mu.Unlock()
			default:
				if stopping {
					continue
				}
				s.mu.Lock()
				msg := s.message
				frame := s.frames[i%len(s.frames)]
				fmt.Printf("\r\033[K\033[36m%s\033[0m %s", frame, msg)
				i++
				time.Sleep(100 * time.Millisecond)
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Spinner) Log(message string) {
	s.mu.Lock()
	s.messageQueue <- message
	s.mu.Unlock()
}

func (s *Spinner) Stop(finalMessage string) {
	close(s.stop)
	<-s.done
	fmt.Printf("\r\033[K\033[32m✓\033[0m %s\n", finalMessage)
	fmt.Print("\033[?25h") // mostra cursor
}

func (s *Spinner) StopWithError(finalMessage string) {
	close(s.stop)
	<-s.done
	fmt.Printf("\r\033[K\033[31m✗\033[0m %s\n", finalMessage)
	fmt.Print("\033[?25h") // mostra cursor
}

func (s *Spinner) Write(p []byte) (n int, err error) {
	s.Log(string(p))
	return len(p), nil
}
