package tool

import "context"


type Tool interface {
	Name() string
	Description() string
	InputSchema() Schema

	Call(
		ctx context.Context,
		input map[string]any,
	) (any, error)
}

type Schema struct {
	Type       string
	Properties map[string]Property
	Required   []string
}

type Property struct {
	Type        string
	Description string
}