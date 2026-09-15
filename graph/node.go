package graph

import "context"

type Node func(
    ctx context.Context,
    state State,
) error

type ConditionalEdge struct {
	Targets []string
	Route   func(State) string
}