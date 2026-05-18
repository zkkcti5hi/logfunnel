package metrics_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"logfunnel/internal/metrics"
)

func TestCounter_IncAndValue(t *testing.T) {
	reg := metrics.NewRegistry()
	c := reg.Counter("lines_in")
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestCounter_Add(t *testing.T) {
	reg := metrics.NewRegistry()
	c := reg.Counter("bytes")
	c.Add(100)
	c.Add(50)
	if c.Value() != 150 {
		t.Fatalf("expected 150, got %d", c.Value())
	}
}

func TestRegistry_SameNameReturnsSameCounter(t *testing.T) {
	reg := metrics.NewRegistry()
	a := reg.Counter("foo")
	b := reg.Counter("foo")
	a.Inc()
	if b.Value() != 1 {
		t.Fatal("expected same counter instance")
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	reg := metrics.NewRegistry()
	reg.Counter("a").Add(3)
	reg.Counter("b").Add(7)
	snap := reg.Snapshot()
	if snap["a"] != 3 || snap["b"] != 7 {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestCounter_ConcurrentInc(t *testing.T) {
	reg := metrics.NewRegistry()
	c := reg.Counter("concurrent")
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if c.Value() != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, c.Value())
	}
}

func TestRegistry_WriteTo(t *testing.T) {
	reg := metrics.NewRegistry()
	reg.Counter("lines_routed").Add(42)
	var buf bytes.Buffer
	_, err := reg.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if !strings.Contains(buf.String(), "lines_routed: 42") {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}
