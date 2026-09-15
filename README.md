# AI Graph

A graph-based workflow engine for Go, inspired by [LangGraph](https://github.com/langchain-ai/langgraph).

AI Graph lets you build stateful workflows from nodes, edges, and conditional routes.

The project is currently under development.

## Features

* Stateful graph execution
* Nodes with `context.Context`
* Regular edges
* Conditional edges
* `START` and `END` points
* Shared state between nodes
* Error propagation
* Configurable graph execution

## Installation

```bash
go get github.com/s-ginger/aigraph
```

## Basic Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/s-ginger/gentic/graph"
)

func main() {
	g := graph.NewGraph()

	g.AddNode("start", func(
		ctx context.Context,
		state graph.State,
	) error {
		state.Set("value", 10)
		return nil
	})

	g.AddNode("double", func(
		ctx context.Context,
		state graph.State,
	) error {
		value, _ := state.Get("value")
		state.Set("value", value.(int)*2)

		return nil
	})

	g.AddEdge("start", "double")
	g.AddEdge("double", graph.End)

	g.SetEntryPoint("start")

	state := graph.NewState()

	if err := g.Run(context.Background(), state); err != nil {
		panic(err)
	}

	value, _ := state.Get("value")

	fmt.Println(value)
}
```

Output:

```text
20
```

## Concepts

### State

`State` is a shared data container passed between nodes.

```go
state := graph.NewState()

state.Set("name", "John")

value, ok := state.Get("name")
```

Each node can read and modify the same state.

### Node

A node is a function that receives a context and the current state.

```go
type Node func(
	ctx context.Context,
	state State,
) error
```

Example:

```go
g.AddNode("process", func(
	ctx context.Context,
	state graph.State,
) error {
	state.Set("status", "processed")
	return nil
})
```

### Edge

An edge defines the next node.

```go
g.AddEdge("start", "process")
```

The execution becomes:

```text
start → process
```

### Conditional Edge

A conditional edge selects the next node based on the current state.

```go
g.AddConditionalEdge("router", func(state graph.State) string {
	value, _ := state.Get("need_search")

	if value.(bool) {
		return "search"
	}

	return "answer"
})
```

This allows workflows such as:

```text
             ┌──→ search ──┐
             │              ↓
START → router            answer → END
             │              ↑
             └──────────────┘
```

### START and END

`START` and `END` represent graph boundaries.

```go
g.SetEntryPoint("agent")

g.AddEdge("agent", graph.End)
```

`START` is represented internally by the configured entry point.

`END` terminates graph execution.

## Execution Model

A graph executes nodes sequentially.

```text
State
  ↓
START
  ↓
Node
  ↓
Edge / Conditional Edge
  ↓
Node
  ↓
END
```

Every node receives the same `State` instance.

This makes it possible to build stateful workflows without coupling nodes directly.

## Example: Agent Workflow

A typical AI workflow can be represented as:

```text
START
  ↓
agent
  ↓
router
 ├──→ tool
 │     ↓
 │   agent
 │
 └──→ answer
       ↓
      END
```

The agent can update the state, the router can inspect it, and the workflow can continue based on the result.

## Project Structure

```text
aigraph/
├── graph/
│   ├── state.go
│   ├── state_test.go
│   ├── node.go
│   ├── graph.go
│   ├── graph_test.go
│   └── route.go
├── llm/
├── memory/
├── tool/
├── examples/
├── go.mod
└── README.md
```

## Roadmap

### Core

* [x] State
* [x] Node
* [x] Graph
* [x] Regular edges
* [x] Conditional edges
* [x] START / END
* [ ] Graph validation
* [ ] Maximum execution steps
* [ ] Context cancellation
* [ ] Execution result

### Runtime

* [ ] Node lifecycle hooks
* [ ] Middleware
* [ ] Logging
* [ ] Metrics
* [ ] Tracing
* [ ] Parallel execution
* [ ] Retry policies

### AI

* [ ] LLM interface
* [ ] Message abstraction
* [ ] Tool interface
* [ ] Tool calling
* [ ] Agent nodes
* [ ] Streaming

### Persistence

* [ ] Checkpoints
* [ ] State persistence
* [ ] Resume execution
* [ ] Conversation history

## Design Goals

The project aims to provide:

* Simple Go APIs
* Explicit workflow structure
* Type-safe interfaces where practical
* Context-aware execution
* Extensible runtime
* Minimal dependencies
* Easy testing

The core graph engine should remain independent from specific LLM providers.

## Status

This project is experimental and under active development.

The API may change before the first stable release.

## License

MIT
