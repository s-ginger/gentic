package memory

import (
	"context"
	"sync"

	"github.com/s-ginger/gentic/graph"
	"github.com/s-ginger/gentic/llm"
)

type InMemory struct {
	mu sync.RWMutex

	config Config

	state    graph.State
	messages []llm.Message
}

func NewInMemory(config Config) *InMemory {
	return &InMemory{
		config: config,
		state:  graph.NewState(),
	}
}

func (m *InMemory) Load(
	ctx context.Context,
) (graph.State, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	state := graph.NewState()

	for _, key := range m.config.PersistentKeys {
		value, ok := m.state.Get(key)
		if !ok {
			continue
		}

		state.Set(key, value)
	}

	return state, nil
}

func (m *InMemory) Save(
	ctx context.Context,
	state graph.State,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range m.config.PersistentKeys {
		value, ok := state.Get(key)
		if !ok {
			continue
		}

		m.state.Set(key, value)
	}

	return nil
}

func (m *InMemory) AddMessage(
	ctx context.Context,
	message llm.Message,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, existing := range m.messages {
		if existing == message {
			return nil
		}
	}

	m.messages = append(m.messages, message)

	return nil
}

func (m *InMemory) Messages(
	ctx context.Context,
) ([]llm.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	messages := make([]llm.Message, len(m.messages))

	copy(messages, m.messages)

	return messages, nil
}