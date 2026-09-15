package graph

import "context"

type Node func(
    ctx context.Context,
    state State,
) error

type Route func(State) string