package filesystem

import (
	"context"
	"fmt"
	"os"

	"github.com/s-ginger/gentic/tools"
)

type ReadFile struct{}

func (ReadFile) Name() string {
	return "read_file"
}

func (ReadFile) Description() string {
	return "Reads the contents of a file."
}

func (ReadFile) InputSchema() tools.Schema {
	return tools.Schema{
		Type: "object",
		Properties: map[string]tools.Property{
			"path": {
				Type:        "string",
				Description: "Path to the file to read.",
			},
		},
		Required: []string{"path"},
	}
}

func (ReadFile) Call(ctx context.Context, input map[string]any) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	path, ok := input["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path must be a string")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	return string(data), nil
}