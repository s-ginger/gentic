package llm

import (
	"context"
	"fmt"

	"github.com/s-ginger/gentic/graph"
)

type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	CachedTokens int
}

type Request struct {
	Messages []Message
}

type Response struct {
	Message Message
	Usage   Usage
}

type StreamChunk struct {
	Delta   string
	Message *Message
	Usage   *Usage
	Done    bool
}

type LLM interface {
	Generate(
		ctx context.Context,
		request Request,
	) (Response, error)

	Stream(
		ctx context.Context,
		request Request,
	) (<-chan StreamChunk, error)
}


func NewLLMNode(model LLM, cfg Config) graph.Node {
	return func(ctx context.Context, state graph.State) error {
		messages, err := getMessages(state, cfg.MessagesKey)
		if err != nil {
			return err
		}

		response, err := model.Generate(ctx, Request{
			Messages: messages,
		})
		if err != nil {
			return fmt.Errorf("generate LLM response: %w", err)
		}

		state.Set(cfg.ResponseKey, response.Message)

		return nil
	}
}

func getMessages(state graph.State, key string) ([]Message, error) {
	value, ok := state.Get(key)
	if !ok {
		return nil, fmt.Errorf(
			"messages not found in state key %q",
			key,
		)
	}

	messages, ok := value.([]Message)
	if !ok {
		return nil, fmt.Errorf(
			"state key %q contains %T, expected []Message",
			key,
			value,
		)
	}

	return messages, nil
}
