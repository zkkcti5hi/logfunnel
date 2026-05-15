// Package sink provides writers that receive filtered log entries
// and forward them to configured destinations (stdout, file, etc.).
package sink

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Sink represents a destination for log entries.
type Sink interface {
	Write(entry string) error
	Close() error
}

// StdoutSink writes log entries to standard output.
type StdoutSink struct{}

// NewStdout creates a Sink that writes to stdout.
func NewStdout() Sink {
	return &StdoutSink{}
}

func (s *StdoutSink) Write(entry string) error {
	_, err := fmt.Fprintln(os.Stdout, entry)
	return err
}

func (s *StdoutSink) Close() error { return nil }

// FileSink writes log entries to a file, creating it if necessary.
type FileSink struct {
	mu sync.Mutex
	f  io.WriteCloser
}

// NewFile opens (or creates) the file at path and returns a FileSink.
func NewFile(path string) (Sink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("sink: open file %q: %w", path, err)
	}
	return &FileSink{f: f}, nil
}

func (s *FileSink) Write(entry string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := fmt.Fprintln(s.f, entry)
	return err
}

func (s *FileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.f.Close()
}
