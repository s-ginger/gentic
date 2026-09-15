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
}

func NewGraph() *Graph {
	return &Graph{
		nodes:            make(map[string]Node),
		edges:            make(map[string]string),
		conditionalEdges: make(map[string]ConditionalEdge),
	}
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

	current := g.entryPoint

	for current != End {
		node, ok := g.nodes[current]
		if !ok {
			return fmt.Errorf("node %q not found", current)
		}

		if err := node(ctx, state); err != nil {
			return err
		}

		if edge, ok := g.conditionalEdges[current]; ok {
			current = edge.Route(state)
			continue
		}

		current = g.edges[current]
	}

	return nil
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