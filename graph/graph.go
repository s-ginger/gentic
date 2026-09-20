package graph

import (
	"context"
	"fmt"
)

const (
	Start = "__start__"
	End   = "__end__"
)

type Graph struct {
	nodes            map[string]Node
	edges            map[string]string
	conditionalEdges map[string]ConditionalEdge
	entryPoint       string

	hooks []Hook
}

func NewGraph() *Graph {
	return &Graph{
		nodes:            make(map[string]Node),
		edges:            make(map[string]string),
		conditionalEdges: make(map[string]ConditionalEdge),
	}
}

// Use adds one or more hooks to the graph.
//
// Hooks execute in the order they were added.
func (g *Graph) Use(hooks ...Hook) {
	g.hooks = append(g.hooks, hooks...)
}

func (g *Graph) AddNode(name string, node Node) {
	g.nodes[name] = node
}

func (g *Graph) AddEdge(from string, to string) {
	g.edges[from] = to
}

func (g *Graph) AddConditionalEdge(
	from string,
	targets []string,
	route func(State) string,
) {
	g.conditionalEdges[from] = ConditionalEdge{
		Targets: targets,
		Route:   route,
	}
}

func (g *Graph) SetEntryPoint(name string) {
	g.entryPoint = name
}

func (g *Graph) Run(
	ctx context.Context,
	state State,
) error {
	if err := g.Validate(); err != nil {
		return err
	}

	run := RunContext{
		State: state,
	}

	if err := g.runBeforeRun(ctx, run); err != nil {
		return err
	}

	current := g.entryPoint

	for current != End {
		node, ok := g.nodes[current]
		if !ok {
			err := fmt.Errorf("node %q not found", current)

			g.runError(ctx, ErrorContext{
				Node:  current,
				State: state,
				Err:   err,
			})

			return err
		}

		nodeContext := NodeContext{
			Name:  current,
			State: state,
		}

		if err := g.runBeforeNode(ctx, nodeContext); err != nil {
			g.runError(ctx, ErrorContext{
				Node:  current,
				State: state,
				Err:   err,
			})

			return err
		}

		if err := node(ctx, state); err != nil {
			g.runError(ctx, ErrorContext{
				Node:  current,
				State: state,
				Err:   err,
			})

			return err
		}

		if err := g.runAfterNode(ctx, nodeContext); err != nil {
			g.runError(ctx, ErrorContext{
				Node:  current,
				State: state,
				Err:   err,
			})

			return err
		}

		if edge, ok := g.conditionalEdges[current]; ok {
			current = edge.Route(state)
			continue
		}

		current = g.edges[current]
	}

	if err := g.runAfterRun(ctx, run); err != nil {
		g.runError(ctx, ErrorContext{
			State: state,
			Err:   err,
		})

		return err
	}

	return nil
}

func (g *Graph) runBeforeRun(
	ctx context.Context,
	run RunContext,
) error {
	for _, hook := range g.hooks {
		if h, ok := hook.(BeforeRunHook); ok {
			if err := h.BeforeRun(ctx, run); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Graph) runAfterRun(
	ctx context.Context,
	run RunContext,
) error {
	for _, hook := range g.hooks {
		if h, ok := hook.(AfterRunHook); ok {
			if err := h.AfterRun(ctx, run); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Graph) runBeforeNode(
	ctx context.Context,
	node NodeContext,
) error {
	for _, hook := range g.hooks {
		if h, ok := hook.(BeforeNodeHook); ok {
			if err := h.BeforeNode(ctx, node); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Graph) runAfterNode(
	ctx context.Context,
	node NodeContext,
) error {
	for _, hook := range g.hooks {
		if h, ok := hook.(AfterNodeHook); ok {
			if err := h.AfterNode(ctx, node); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Graph) runError(
	ctx context.Context,
	errContext ErrorContext,
) {
	for _, hook := range g.hooks {
		if h, ok := hook.(ErrorHook); ok {
			h.OnError(ctx, errContext)
		}
	}
}

func (g *Graph) Validate() error {
	if g.entryPoint == "" {
		return fmt.Errorf("entry point is not set")
	}

	if _, ok := g.nodes[g.entryPoint]; !ok {
		return fmt.Errorf(
			"entry point %q does not exist",
			g.entryPoint,
		)
	}

	for from, to := range g.edges {
		if _, ok := g.nodes[from]; !ok {
			return fmt.Errorf(
				"edge source node %q does not exist",
				from,
			)
		}

		if to == End {
			continue
		}

		if _, ok := g.nodes[to]; !ok {
			return fmt.Errorf(
				"edge from %q points to unknown node %q",
				from,
				to,
			)
		}
	}

	for from, edge := range g.conditionalEdges {
		if _, ok := g.nodes[from]; !ok {
			return fmt.Errorf(
				"conditional edge source node %q does not exist",
				from,
			)
		}

		for _, target := range edge.Targets {
			if target == End {
				continue
			}

			if _, ok := g.nodes[target]; !ok {
				return fmt.Errorf(
					"conditional edge from %q points to unknown node %q",
					from,
					target,
				)
			}
		}
	}

	return nil
}
