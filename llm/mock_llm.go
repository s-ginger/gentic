package llm

import (
	"context"
)

type MockLLM struct {
	Response Response
}

func (m *MockLLM) Generate(
	ctx context.Context,
	request Request,
) (Response, error) {
	return m.Response, nil
}

func (m *MockLLM) Stream(
	ctx context.Context,
	request Request,
) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk)

	go func() {
		defer close(ch)

		message := m.Response.Message.Content

		for _, r := range message {
			select {
			case <-ctx.Done():
				return
			case ch <- StreamChunk{
				Delta: string(r),
			}:
			}
		}

		ch <- StreamChunk{
			Done: true,
			Usage: &m.Response.Usage,
		}
	}()

	return ch, nil
}