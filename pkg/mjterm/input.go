package mjterm

import "context"

func (t *Terminal) Prompt(prompt string) (string, error) {
	return t.PromptWithContext(context.Background(), prompt)
}

func (t *Terminal) PromptWithContext(ctx context.Context, prompt string) (string, error) {
	resultCh := make(chan promptResult)

	t.addEvent(promptMsg{
		content: prompt,
		result:  resultCh,
	})

	select {
	case result := <-resultCh:
		return result.response, result.err
	case <-ctx.Done():
		close(resultCh)
		return "", ctx.Err()
	}
}
