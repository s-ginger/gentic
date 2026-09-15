package llm

import "context"

type MockLLM struct {
	Response Response
}

func (m *MockLLM) Generate(
	ctx context.Context,
	request Request,
) (Response, error) {
	return m.Response, nil
}