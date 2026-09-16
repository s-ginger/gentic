package llm

import (
	"context"
	"github.com/s-ginger/gentic/graph"
	"testing"
)

func TestLLMNode(t *testing.T) {
	llm := &MockLLM{
		Response: Response{
			Message: Message{
				Role:    RoleAssistant,
				Content: "Hello from LLM",
			},
		},
	}

	node := NewLLMNode(
		llm,
		Config{
			MessagesKey: "messages",
			ResponseKey: "response",
		},
	)

	state := graph.NewState()

	state.Set("messages", []Message{
		{
			Role:    RoleUser,
			Content: "Hello",
		},
	})

	err := node(
		context.Background(),
		state,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok := state.Get("response")

	if !ok {
		t.Fatal("response not found")
	}

	response, ok := value.(Message)

	if !ok {
		t.Fatalf("expected Message, got %T", value)
	}

	if response.Content != "Hello from LLM" {
		t.Fatalf(
			"expected %q, got %q",
			"Hello from LLM",
			response.Content,
		)
	}
}