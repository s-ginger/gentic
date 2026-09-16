package llm

import (
	"context"
	"fmt"
	"github.com/s-ginger/gentic/graph"
)

type Request struct {
	Messages []Message
}

type Response struct {
	Message Message
}

type LLM interface {
	Generate(
		ctx context.Context,
		request Request,
	) (Response, error)
}

func NewLLMNode(
	model LLM,
	config Config,
) graph.Node {
	return func(
		ctx context.Context,
		state graph.State,
	) error {
		value, ok := state.Get(config.MessagesKey)
		if !ok {
			return fmt.Errorf(
				"messages not found in state key %q",
				config.MessagesKey,
			)
		}

		messages, ok := value.([]Message)
		if !ok {
			return fmt.Errorf(
				"invalid messages type",
			)
		}

		response, err := model.Generate(
			ctx,
			Request{
				Messages: messages,
			},
		)
		if err != nil {
			return err
		}

		state.Set(
			config.ResponseKey,
			response.Message,
		)

		return nil
	}
}
