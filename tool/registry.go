package tool

// Registry stores tools by name.
type Registry map[string]Tool

// NewRegistry creates an empty tool registry.
func NewRegistry() *Registry {
	r := Registry{}

	return &r
}

// Register adds a tool to the registry.
func (r *Registry) Register(t Tool) {
	(*r)[t.Name()] = t
}

// Get returns a tool by name.
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := (*r)[name]

	return t, ok
}

// List returns all registered tools.
func (r *Registry) List() []Tool {
	tools := make([]Tool, 0, len(*r))

	for _, t := range *r {
		tools = append(tools, t)
	}

	return tools
}