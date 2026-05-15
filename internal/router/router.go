// Package router wires together tailers, filters, and sinks so that
// log lines from each source are evaluated against every rule and
// dispatched to the matching sink.
package router

import (
	"log"
	"sync"

	"github.com/user/logfunnel/internal/filter"
	"github.com/user/logfunnel/internal/sink"
	"github.com/user/logfunnel/internal/tail"
)

// Route describes a single routing rule: a compiled filter and the
// sink that should receive matching lines.
type Route struct {
	Filter *filter.Filter
	Sink   sink.Sink
}

// Router fans out lines from one or more tailers through a set of
// routes and forwards matches to the appropriate sink.
type Router struct {
	tailers []*tail.Tailer
	routes  []Route
}

// New creates a Router that will read from the supplied tailers and
// apply the given routes to every line that arrives.
func New(tailers []*tail.Tailer, routes []Route) *Router {
	return &Router{
		tailers: tailers,
		routes:  routes,
	}
}

// Run starts one goroutine per tailer and blocks until all of them
// have finished (i.e. their line channels are closed).
func (r *Router) Run() {
	var wg sync.WaitGroup
	for _, t := range r.tailers {
		wg.Add(1)
		go func(t *tail.Tailer) {
			defer wg.Done()
			for line := range t.Lines() {
				r.dispatch(line)
			}
		}(t)
	}
	wg.Wait()
}

// dispatch evaluates a single log line against every route and writes
// it to the sink of each matching route.
func (r *Router) dispatch(line string) {
	matched := false
	for _, route := range r.routes {
		if sinkName, ok := route.Filter.Match(line); ok {
			_ = sinkName
			if err := route.Sink.Write(line); err != nil {
				log.Printf("router: sink write error: %v", err)
			}
			matched = true
		}
	}
	if !matched {
		log.Printf("router: no route matched line: %s", line)
	}
}
