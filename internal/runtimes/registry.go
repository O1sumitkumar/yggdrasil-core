package runtimes

import "fmt"

// Registry holds registered runtime adapters.
type Registry struct {
	runtimes map[string]Runtime
	order    []string
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{runtimes: make(map[string]Runtime)}
}

// Register adds a runtime adapter.
func (r *Registry) Register(rt Runtime) {
	id := rt.ID()
	if _, exists := r.runtimes[id]; !exists {
		r.order = append(r.order, id)
	}
	r.runtimes[id] = rt
}

// Get returns a runtime by ID.
func (r *Registry) Get(id string) (Runtime, error) {
	rt, ok := r.runtimes[id]
	if !ok {
		return nil, fmt.Errorf("runtime %q not registered", id)
	}
	return rt, nil
}

// List returns registered runtimes in registration order.
func (r *Registry) List() []Runtime {
	out := make([]Runtime, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.runtimes[id])
	}
	return out
}
