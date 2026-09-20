package graph

import "context"

type RunContext struct {
	State State
}

type NodeContext struct {
	Name  string
	State State
}

type ErrorContext struct {
	Node  string
	State State
	Err   error
}

type BeforeRunHook interface {
	BeforeRun(context.Context, RunContext) error
}

type AfterRunHook interface {
	AfterRun(context.Context, RunContext) error
}

type BeforeNodeHook interface {
	BeforeNode(context.Context, NodeContext) error
}

type AfterNodeHook interface {
	AfterNode(context.Context, NodeContext) error
}

type ErrorHook interface {
	OnError(context.Context, ErrorContext) error
}

type Hook interface{}