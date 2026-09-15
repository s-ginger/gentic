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
	nodes  map[string]Node
	edges  map[string]string
	routes map[string]Route

	entryPoint string
}

func NewGraph() *Graph {
	return &Graph{
		nodes:  make(map[string]Node),
		edges:  make(map[string]string),
		routes: make(map[string]Route),
	}
}

func (g *Graph) AddNode(name string, node Node) {
    g.nodes[name] = node
}

func (g *Graph) AddEdge(from string, to string) {
    g.edges[from] = to
}

func (g *Graph) AddConditionalEdge(from string, route Route) {
	g.routes[from] = route
}

func (g *Graph) SetEntryPoint(name string) {
	g.entryPoint = name
}

func (g *Graph) Run(
	ctx context.Context,
	state State,
) error {
	current := g.entryPoint

	for current != "" && current != End {
		node, ok := g.nodes[current]
		if !ok {
			return fmt.Errorf("node %q not found", current)
		}

		if err := node(ctx, state); err != nil {
			return err
		}

		if route, ok := g.routes[current]; ok {
			current = route(state)
			continue
		}

		current = g.edges[current]
	}

	return nil
}

