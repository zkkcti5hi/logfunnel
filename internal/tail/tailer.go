package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Entry represents a single log line read from a source.
type Entry struct {
	Source string
	Line   string
	TS     time.Time
}

// Tailer tails a log file and emits entries on a channel.
type Tailer struct {
	path   string
	out    chan<- Entry
	pollInterval time.Duration
}

// New creates a new Tailer for the given file path.
func New(path string, out chan<- Entry, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = 500 * time.Millisecond
	}
	return &Tailer{
		path:         path,
		out:          out,
		pollInterval: pollInterval,
	}
}

// Run opens the file, seeks to the end, and streams new lines until ctx is cancelled.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			// No new data yet; back off and retry.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(t.pollInterval):
			}
			continue
		}

		if len(line) > 0 {
			select {
			case t.out <- Entry{Source: t.path, Line: line, TS: time.Now()}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
