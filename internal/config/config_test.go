package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	content := `
sources:
  - name: app
    path: /var/log/app.log
sinks:
  - name: errors_out
    type: file
    path: /tmp/errors.log
rules:
  - pattern: "ERROR"
    sink: errors_out
`
	path := writeTemp(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(cfg.Sources) != 1 || cfg.Sources[0].Name != "app" {
		t.Errorf("unexpected sources: %+v", cfg.Sources)
	}
	if len(cfg.Sinks) != 1 || cfg.Sinks[0].Type != "file" {
		t.Errorf("unexpected sinks: %+v", cfg.Sinks)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Pattern != "ERROR" {
		t.Errorf("unexpected rules: %+v", cfg.Rules)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_NoSources(t *testing.T) {
	cfg := &Config{
		Sinks: []Sink{{Name: "out", Type: "stdout"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing sources")
	}
}

func TestValidate_UnknownSinkInRule(t *testing.T) {
	cfg := &Config{
		Sources: []Source{{Name: "s", Path: "/tmp/x"}},
		Sinks:   []Sink{{Name: "out", Type: "stdout"}},
		Rules:   []Rule{{Pattern: "ERR", Sink: "ghost"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unknown sink in rule")
	}
}

func TestValidate_FileSinkMissingPath(t *testing.T) {
	cfg := &Config{
		Sources: []Source{{Name: "s", Path: "/tmp/x"}},
		Sinks:   []Sink{{Name: "out", Type: "file"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for file sink without path")
	}
}
