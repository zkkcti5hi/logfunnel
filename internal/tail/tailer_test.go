package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logfunnel/internal/tail"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logfunnel-tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	out := make(chan tail.Entry, 10)
	tr := tail.New(f.Name(), out, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- tr.Run(ctx)
	}()

	// Give the tailer time to seek to EOF before we write.
	time.Sleep(80 * time.Millisecond)

	writeLine(t, f, "hello world")
	writeLine(t, f, "second line")

	received := make([]string, 0, 2)
	timeout := time.After(2 * time.Second)
	for len(received) < 2 {
		select {
		case entry := <-out:
			received = append(received, entry.Line)
		case <-timeout:
			t.Fatalf("timed out waiting for log entries; got %d/2", len(received))
		}
	}

	cancel()
	<-errCh

	if received[0] != "hello world\n" {
		t.Errorf("expected 'hello world\\n', got %q", received[0])
	}
	if received[1] != "second line\n" {
		t.Errorf("expected 'second line\\n', got %q", received[1])
	}
}

func TestTailer_MissingFile(t *testing.T) {
	out := make(chan tail.Entry, 1)
	tr := tail.New("/nonexistent/path/to/file.log", out, 50*time.Millisecond)

	err := tr.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
