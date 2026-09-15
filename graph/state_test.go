package graph

import "testing"

func TestNewState(t *testing.T) {
    state := NewState()

    if state == nil {
        t.Fatal("state is nil")
    }

    state.Set("name", "test")

    value, ok := state.Get("name")

    if !ok {
        t.Fatal("value not found")
    }

    if value != "test" {
        t.Fatalf("expected %q, got %q", "test", value)
    }
}

func TestStateGetMissing(t *testing.T) {
    state := NewState()

    _, ok := state.Get("missing")

    if ok {
        t.Fatal("expected value to be missing")
    }
}