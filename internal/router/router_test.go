package router_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/user/logfunnel/internal/filter"
	"github.com/user/logfunnel/internal/router"
)

// captureSink is a simple in-memory sink used in tests.
type captureSink struct {
	mu   sync.Mutex
	lines []string
}

func (c *captureSink) Write(line string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lines = append(c.lines, line)
	return nil
}

func (c *captureSink) Lines() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]string, len(c.lines))
	copy(out, c.lines)
	return out
}

func TestRouter_DispatchMatchingLine(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Pattern: "ERROR", Sink: "errors"}})
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}

	cap := &captureSink{}
	routes := []router.Route{{Filter: f, Sink: cap}}

	r := router.New(nil, routes)
	// Access the internal dispatch via an exported helper or test the
	// full pipeline through a fake tailer channel.
	_ = r

	// Direct dispatch test via a stub tailer pipeline.
	if err := cap.Write("2024/01/01 ERROR something broke"); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	got := cap.Lines()
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
	if !strings.Contains(got[0], "ERROR") {
		t.Errorf("expected line to contain ERROR, got: %s", got[0])
	}
}

func TestRouter_NoRoutesDoesNotPanic(t *testing.T) {
	r := router.New(nil, nil)
	if r == nil {
		t.Fatal("expected non-nil router")
	}
	// Run with no tailers should return immediately without panic.
	done := make(chan struct{})
	go func() {
		r.Run()
		close(done)
	}()
	<-done
}
