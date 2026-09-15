package graph

type State map[string]any

func NewState() State {
    return make(State)
}

func (s State) Set(key string, value any) {
    s[key] = value
}

func (s State) Get(key string) (any, bool) {
    value, ok := s[key]
    return value, ok
}