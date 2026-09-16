package tool

// Registry stores tools by their unique names.
//
// A Registry can be used by an application or a custom graph node
// to discover and execute tools.
//
// Example:
//
//	registry := tool.NewRegistry()
//	registry.Register(MyTool{})
//
//	t, ok := registry.Get("my_tool")
//	if !ok {
//		// Tool is not registered.
//	}
//
//	result, err := t.Call(ctx, input)
type Registry map[string]Tool

// NewRegistry creates an empty tool registry.
func NewRegistry() Registry {
	return make(Registry)
}

// Register adds a tool to the registry.
//
// If a tool with the same name is already registered, it is replaced.
func (r Registry) Register(t Tool) {
	r[t.Name()] = t
}

// Get returns the tool registered under name.
//
// The second return value is false when no tool with the given name exists.
func (r Registry) Get(name string) (Tool, bool) {
	t, ok := r[name]
	return t, ok
}

// List returns all tools currently registered.
//
// The returned slice is a new slice and can be modified by the caller
// without changing the registry.
func (r Registry) List() []Tool {
	tools := make([]Tool, 0, len(r))

	for _, t := range r {
		tools = append(tools, t)
	}

	return tools
}