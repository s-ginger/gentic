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
	graph.AddEdge("second", End)

	graph.SetEntryPoint("first")

	state := NewState()

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

	graph.AddNode("router", func(
		ctx context.Context,
		state State,
	) error {
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

	graph.AddEdge("start", "router")

	graph.AddConditionalEdge(
		"router",
		[]string{
			"search",
			"answer",
		},
		func(state State) string {
			value, _ := state.Get("need_search")

			if value.(bool) {
				return "search"
			}

			return "answer"
		},
	)

	graph.AddEdge("search", End)
	graph.AddEdge("answer", End)

	graph.SetEntryPoint("start")

	state := NewState()

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
func TestGraphValidate(t *testing.T) {
	graph := NewGraph()

	graph.AddNode("start", func(
		ctx context.Context,
		state State,
	) error {
		return nil
	})

	graph.AddNode("finish", func(
		ctx context.Context,
		state State,
	) error {
		return nil
	})

	graph.SetEntryPoint("start")
	graph.AddEdge("start", "finish")
	graph.AddEdge("finish", End)

	err := graph.Validate()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGraphValidateWithoutEntryPoint(t *testing.T) {
	graph := NewGraph()

	err := graph.Validate()

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGraphValidateUnknownEdgeTarget(t *testing.T) {
	graph := NewGraph()

	graph.AddNode("start", func(
		ctx context.Context,
		state State,
	) error {
		return nil
	})

	graph.SetEntryPoint("start")
	graph.AddEdge("start", "unknown")

	err := graph.Validate()

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGraphConditionalEdgeInvalidTarget(t *testing.T) {
	graph := NewGraph()

	graph.AddNode("router", func(
		ctx context.Context,
		state State,
	) error {
		return nil
	})

	graph.AddConditionalEdge(
		"router",
		[]string{"search", "answer"},
		func(state State) string {
			return "unknown"
		},
	)

	graph.SetEntryPoint("router")

	err := graph.Run(
		context.Background(),
		NewState(),
	)

	if err == nil {
		t.Fatal("expected error")
	}
}
