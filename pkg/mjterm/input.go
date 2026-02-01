package mjterm

import (
	"bufio"
	"context"
	"io"
	"os"
)

type scanResult struct {
	text string
	err  error
}

func (t *Terminal) PromptWithContext(ctx context.Context, prompt string) (string, error) {
	if err := t.pauseOutput(); err != nil {
		return "", err
	}

	defer t.resumeOutput()

	if err := t.priorityPrint(prompt); err != nil {
		return "", err
	}

	resultCh := make(chan scanResult)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		sendToResultCh := func(text string, err error) {
			select {
			case resultCh <- scanResult{text: text, err: err}:
			case <-ctx.Done():
			case <-t.close:
			}
		}

		if scanner.Scan() {
			sendToResultCh(scanner.Text(), nil)
			return
		}
		if err := scanner.Err(); err != nil {
			sendToResultCh("", err)
			return
		}
		sendToResultCh("", io.EOF)
	}()

	select {
	case result := <-resultCh:
		return result.text, result.err
	case <-ctx.Done():
		return "", ctx.Err()
	case <-t.close:
		return "", ErrTerminalClosed
	}
}

func (t *Terminal) Prompt(prompt string) (string, error) {
	return t.PromptWithContext(context.Background(), prompt)
}
