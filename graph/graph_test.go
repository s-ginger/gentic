package graph

import (
	"context"
	"errors"
	"testing"
)

func TestGraphRun(t *testing.T) {
	graph := NewGraph()

	graph.AddNode("first", func(
		ctx context.Context,
		state State,
	) error {
		state.Set("value", 10)
		return nil
	})

	graph.AddNode("second", func(
		ctx context.Context,
		state State,
	) error {
		value, _ := state.Get("value")

		state.Set("value", value.(int)*2)

		return nil
	})

	graph.AddEdge("first", "second")

	state := NewState()

	graph.SetEntryPoint("first")

	err := graph.Run(
		context.Background(),
		state,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok := state.Get("value")

	if !ok {
		t.Fatal("value not found")
	}

	if value != 20 {
		t.Fatalf("expected 20, got %v", value)
	}
}

func TestGraphRunNodeError(t *testing.T) {
	graph := NewGraph()

	expectedErr := errors.New("node failed")

	graph.AddNode("broken", func(
		ctx context.Context,
		state State,
	) error {
		return expectedErr
	})

	state := NewState()

	graph.SetEntryPoint("broken")

	err := graph.Run(
		context.Background(),
		state,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestGraphRunUnknownNode(t *testing.T) {
	graph := NewGraph()
	state := NewState()

	graph.SetEntryPoint("unknown")

	err := graph.Run(
		context.Background(),
		state,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGraphConditionalEdge(t *testing.T) {
	graph := NewGraph()

	graph.AddNode("start", func(
		ctx context.Context,
		state State,
	) error {
		state.Set("need_search", true)
		return nil
	})

	graph.AddNode("search", func(
		ctx context.Context,
		state State,
	) error {
		state.Set("result", "searched")
		return nil
	})

	graph.AddNode("answer", func(
		ctx context.Context,
		state State,
	) error {
		state.Set("result", "answered")
		return nil
	})

	graph.AddConditionalEdge("start", func(state State) string {
		value, _ := state.Get("need_search")

		if value.(bool) {
			return "search"
		}

		return "answer"
	})

	state := NewState()

	graph.SetEntryPoint("start")

	err := graph.Run(
		context.Background(),
		state,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok := state.Get("result")

	if !ok {
		t.Fatal("result not found")
	}

	if value != "searched" {
		t.Fatalf("expected searched, got %v", value)
	}
}
