package memory

import (
	"context"

	"github.com/s-ginger/gentic/graph"
	"github.com/s-ginger/gentic/llm"
)

type Memory interface {
	Load(ctx context.Context) (graph.State, error)

	Save(ctx context.Context, state graph.State) error

	AddMessage(ctx context.Context, message llm.Message) error

	Messages(ctx context.Context) ([]llm.Message, error)
}

