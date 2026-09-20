package tools

import (
	"context"
	"fmt"
)

// Execute validates the input and then executes the tool.
func Execute(
	ctx context.Context,
	t Tool,
	input map[string]any,
) (any, error) {
	if t == nil {
		return nil, fmt.Errorf("tool must not be nil")
	}

	if err := t.InputSchema().Validate(input); err != nil {
		return nil, fmt.Errorf(
			"invalid input for tool %q: %w",
			t.Name(),
			err,
		)
	}

	return t.Call(ctx, input)
}
