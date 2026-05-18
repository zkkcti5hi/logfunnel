// Package metrics provides lightweight in-process counters for
// tracking log lines processed, routed, and dropped by logfunnel.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Counter is a thread-safe monotonically increasing counter.
type Counter struct {
	value atomic.Int64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { c.value.Add(1) }

// Add increments the counter by n.
func (c *Counter) Add(n int64) { c.value.Add(n) }

// Value returns the current counter value.
func (c *Counter) Value() int64 { return c.value.Load() }

// Registry holds a named set of counters.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{counters: make(map[string]*Counter)}
}

// Counter returns the named counter, creating it if it does not exist.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Snapshot returns a copy of all counter values keyed by name.
func (r *Registry) Snapshot() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]int64, len(r.counters))
	for name, c := range r.counters {
		out[name] = c.Value()
	}
	return out
}

// WriteTo writes a human-readable summary of all counters to w.
func (r *Registry) WriteTo(w io.Writer) (int64, error) {
	snap := r.Snapshot()
	var total int64
	for name, val := range snap {
		n, err := fmt.Fprintf(w, "%s: %d\n", name, val)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
