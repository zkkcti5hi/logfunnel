package pipeline_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/logfunnel/internal/config"
	"github.com/user/logfunnel/internal/pipeline"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "src-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestBuild_ValidConfig(t *testing.T) {
	srcPath := writeTempLog(t, "")
	cfg := &config.Config{
		Sources: []config.Source{{Name: "app", Path: srcPath}},
		Sinks:   []config.SinkConfig{{Name: "out", Type: "stdout"}},
		Rules:   []config.Rule{{Source: "app", Pattern: ".*", Sink: "out"}},
	}
	_, err := pipeline.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_UnknownSinkType(t *testing.T) {
	srcPath := writeTempLog(t, "")
	cfg := &config.Config{
		Sources: []config.Source{{Name: "app", Path: srcPath}},
		Sinks:   []config.SinkConfig{{Name: "bad", Type: "kafka"}},
		Rules:   []config.Rule{{Source: "app", Pattern: ".*", Sink: "bad"}},
	}
	_, err := pipeline.Build(cfg)
	if err == nil {
		t.Fatal("expected error for unknown sink type")
	}
}

func TestRun_RoutesLinesToFileSink(t *testing.T) {
	srcPath := writeTempLog(t, "")
	dstPath := filepath.Join(t.TempDir(), "out.log")

	cfg := &config.Config{
		Sources: []config.Source{{Name: "app", Path: srcPath}},
		Sinks:   []config.SinkConfig{{Name: "out", Type: "file", Path: dstPath}},
		Rules:   []config.Rule{{Source: "app", Pattern: "ERROR", Sink: "out"}},
	}

	p, err := pipeline.Build(cfg)
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- p.Run(ctx) }()

	time.Sleep(50 * time.Millisecond)
	f, _ := os.OpenFile(srcPath, os.O_APPEND|os.O_WRONLY, 0644)
	f.WriteString("ERROR something broke\n")
	f.Close()

	<-done

	data, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("reading sink file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected sink file to contain routed line")
	}
}
