package tool

type Registry map[string]Tool

func NewRegistry() *Registry {
	r := Registry{}
	return &r
}

func (r *Registry) Register(t Tool) {
	(*r)[t.Name()] = t
}

func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := (*r)[name]
	return t, ok
}

func (r *Registry) List() []Tool {
	tools := make([]Tool, 0, len(*r))

	for _, t := range *r {
		tools = append(tools, t)
	}

	return tools
}
