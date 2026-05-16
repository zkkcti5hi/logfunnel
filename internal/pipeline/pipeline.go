// Package pipeline wires together tailers, filters, and sinks into a
// running log-processing pipeline driven by the loaded configuration.
package pipeline

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/user/logfunnel/internal/config"
	"github.com/user/logfunnel/internal/filter"
	"github.com/user/logfunnel/internal/router"
	"github.com/user/logfunnel/internal/sink"
	"github.com/user/logfunnel/internal/tail"
)

// Pipeline holds all running components for a logfunnel instance.
type Pipeline struct {
	cfg     *config.Config
	sinks   map[string]sink.Sink
	routers []*router.Router
	wg      sync.WaitGroup
}

// Build constructs a Pipeline from cfg but does not start it.
func Build(cfg *config.Config) (*Pipeline, error) {
	p := &Pipeline{
		cfg:   cfg,
		sinks: make(map[string]sink.Sink),
	}

	for _, sc := range cfg.Sinks {
		var s sink.Sink
		var err error
		switch sc.Type {
		case "stdout":
			s = sink.NewStdout()
		case "file":
			s, err = sink.NewFile(sc.Path)
			if err != nil {
				return nil, fmt.Errorf("sink %q: %w", sc.Name, err)
			}
		default:
			return nil, fmt.Errorf("sink %q: unknown type %q", sc.Name, sc.Type)
		}
		p.sinks[sc.Name] = s
	}

	for _, src := range cfg.Sources {
		f, err := filter.New(cfg.RulesForSource(src.Name))
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", src.Name, err)
		}
		sinkMap := make(map[string]sink.Sink)
		for _, rule := range cfg.RulesForSource(src.Name) {
			sinkMap[rule.Sink] = p.sinks[rule.Sink]
		}
		r := router.New(f, sinkMap)
		p.routers = append(p.routers, r)
	}

	return p, nil
}

// Run starts all tailers and blocks until ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context) error {
	for i, src := range p.cfg.Sources {
		t, err := tail.New(src.Path)
		if err != nil {
			return fmt.Errorf("source %q: %w", src.Name, err)
		}
		r := p.routers[i]
		p.wg.Add(1)
		go func(t *tail.Tailer, r *router.Router, name string) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					t.Stop()
					return
				case line, ok := <-t.Lines():
					if !ok {
						return
					}
					if err := r.Dispatch(line); err != nil {
						log.Printf("[pipeline] source %q dispatch error: %v", name, err)
					}
				}
			}
		}(t, r, src.Name)
	}
	p.wg.Wait()
	return nil
}
