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

	for from := range g.routes {
		if _, ok := g.nodes[from]; !ok {
			return fmt.Errorf(
				"route source node %q does not exist",
				from,
			)
		}
	}

	return nil
}
