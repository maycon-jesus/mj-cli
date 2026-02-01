package mjterm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const defaultEventBufferSize = 10

type Spinner struct {
	id      string
	message string
	frames  []string
	line    int
	ticker  *time.Ticker
	stop    chan struct{}
}

type Terminal struct {
	events       chan msg
	spinner      *Spinner
	closed       atomic.Bool
	wg           sync.WaitGroup
	closeOnce    sync.Once
	inputPending atomic.Bool
}

func New() *Terminal {
	term := &Terminal{
		events:       make(chan msg, defaultEventBufferSize),
		closed:       atomic.Bool{},
		wg:           sync.WaitGroup{},
		closeOnce:    sync.Once{},
		inputPending: atomic.Bool{},
	}

	go func() {
		for event := range term.events {
			switch e := event.(type) {
			case printMsg:
				if term.spinner != nil {
					fmt.Print("\r\033[K") // clear line
				}
				fmt.Print(e.content)
			case promptMsg:
				term.inputPending.Store(true)

				if term.spinner != nil {
					fmt.Print("\r\033[K") // clear line
				}
				fmt.Print(e.content)

				sendToResultCh := func(text string, err error) {
					select {
					case e.result <- promptResult{response: text, err: err}:
					default:
					}
					term.inputPending.Store(false)
				}

				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					sendToResultCh(scanner.Text(), nil)
					break
				}
				if err := scanner.Err(); err != nil {
					sendToResultCh("", err)
					break
				}
				sendToResultCh("", io.EOF)
			case startSpinnerMsg:
				spinner := &Spinner{
					id:      e.id,
					message: e.message,
					frames:  e.frames,
					line:    e.line,
					ticker:  time.NewTicker(100 * time.Millisecond),
					stop:    make(chan struct{}),
				}
				term.spinner = spinner

				go func(s *Spinner) {
					i := 0
					for {
						select {
						case <-s.ticker.C:
							if term.inputPending.Load() {
								continue
							}
							msg := fmt.Sprintf("\r\033[K\033[36m%s\033[0m %s", s.frames[i%len(s.frames)], s.message)
							term.addEvent(printMsg{content: msg})
							i++
						case <-s.stop:
							s.ticker.Stop()
							term.addEvent(printMsg{content: "\r\033[K"})
							return
						}
					}
				}(spinner)
			case stopSpinnerMsg:
				if term.spinner == nil {
					break
				}
				close(term.spinner.stop)
				term.spinner = nil
				finalMsg := ""
				if e.err {
					finalMsg = fmt.Sprintf("\r\033[K\033[31m✗\033[0m %s\n", e.message)
				} else {
					finalMsg = fmt.Sprintf("\r\033[K\033[32m✓\033[0m %s\n", e.message)
				}
				term.addEvent(printMsg{content: finalMsg})
			case closeMsg:
				if term.spinner != nil {
					close(term.spinner.stop)
					term.spinner = nil
				}
				close(e.closed)
			}
			term.wg.Done()
		}
	}()

	return term
}

func (t *Terminal) Close() {
	t.closeOnce.Do(func() {
		closed := make(chan struct{})
		t.addEvent(closeMsg{closed})
		t.closed.Store(true)
		<-closed
		t.wg.Wait()
		close(t.events)
	})
}
