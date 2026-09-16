package memory

import (
	"context"

	"github.com/s-ginger/gentic/graph"
	"github.com/s-ginger/gentic/llm"
)

// Memory provides persistent storage for graph state and conversation
// messages.
//
// Implementations may store data in memory, a database, a file, or another
// storage backend.
//
// Memory is independent of Graph. The application is responsible for
// loading state before running a graph and saving it afterwards.
//
// Example:
//
//	state, err := memory.Load(ctx)
//	if err != nil {
//		return err
//	}
//
//	if err := graph.Run(ctx, state); err != nil {
//		return err
//	}
//
//	return memory.Save(ctx, state)
type Memory interface {
	// Load returns the previously saved state.
	//
	// If no state has been saved yet, implementations should return
	// an empty state rather than nil when possible.
	Load(ctx context.Context) (graph.State, error)

	// Save persists the provided state.
	//
	// Implementations should only persist the data they are responsible
	// for storing.
	Save(ctx context.Context, state graph.State) error

	// AddMessage adds a message to the conversation history.
	AddMessage(ctx context.Context, message llm.Message) error

	// Messages returns the conversation history.
	//
	// Implementations should return a copy or otherwise prevent callers
	// from modifying the internal message collection directly.
	Messages(ctx context.Context) ([]llm.Message, error)
}