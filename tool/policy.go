package tool

import "context"

type Policy interface {
    Check(ctx context.Context, tool Tool, input map[string]any) error
}