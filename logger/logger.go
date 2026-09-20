package logger

import (
	"context"
	"log/slog"

	"github.com/s-ginger/gentic/graph"
)

type Hook struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Hook {
	return &Hook{
		logger: logger,
	}
}

func (h *Hook) BeforeRun(
	ctx context.Context,
	run graph.RunContext,
) error {
	h.logger.InfoContext(ctx, "graph started")

	return nil
}

func (h *Hook) AfterRun(
	ctx context.Context,
	run graph.RunContext,
) error {
	h.logger.InfoContext(ctx, "graph finished")

	return nil
}

func (h *Hook) BeforeNode(
	ctx context.Context,
	node graph.NodeContext,
) error {
	h.logger.InfoContext(
		ctx,
		"node started",
		slog.String("node", node.Name),
	)

	return nil
}

func (h *Hook) AfterNode(
	ctx context.Context,
	node graph.NodeContext,
) error {
	h.logger.InfoContext(
		ctx,
		"node finished",
		slog.String("node", node.Name),
	)

	return nil
}

func (h *Hook) OnError(
	ctx context.Context,
	errContext graph.ErrorContext,
) {
	h.logger.ErrorContext(
		ctx,
		"graph error",
		slog.String("node", errContext.Node),
		slog.Any("error", errContext.Err),
	)
}
