package sink_test

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logfunnel/internal/sink"
)

func TestStdoutSink_Write(t *testing.T) {
	s := sink.NewStdout()
	if err := s.Write("hello stdout"); err != nil {
		t.Fatalf("unexpected error writing to stdout sink: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error closing stdout sink: %v", err)
	}
}

func TestFileSink_WriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	s, err := sink.NewFile(path)
	if err != nil {
		t.Fatalf("NewFile: %v", err)
	}

	lines := []string{"first line", "second line", "third line"}
	for _, l := range lines {
		if err := s.Write(l); err != nil {
			t.Fatalf("Write(%q): %v", l, err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open result file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var got []string
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: want %q, got %q", i, want, got[i])
		}
	}
}

func TestFileSink_InvalidPath(t *testing.T) {
	_, err := sink.NewFile("/nonexistent-dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
