package filter_test

import (
	"testing"

	"github.com/yourorg/logfunnel/internal/filter"
)

func TestNew_ValidPatterns(t *testing.T) {
	rules := map[string]string{
		"errors": `(?i)error`,
		"warnings": `(?i)warn`,
	}
	f, err := filter.New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules()) != 2 {
		t.Errorf("expected 2 rules, got %d", len(f.Rules()))
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	rules := map[string]string{
		"bad": `[invalid`,
	}
	_, err := filter.New(rules)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatch_SingleRule(t *testing.T) {
	f, err := filter.New(map[string]string{
		"errors": `(?i)error`,
	})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	sinks := f.Match("2024/01/02 ERROR something went wrong")
	if len(sinks) != 1 || sinks[0] != "errors" {
		t.Errorf("expected [errors], got %v", sinks)
	}
}

func TestMatch_NoMatch(t *testing.T) {
	f, err := filter.New(map[string]string{
		"errors": `(?i)error`,
	})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	sinks := f.Match("INFO everything is fine")
	if len(sinks) != 0 {
		t.Errorf("expected no matches, got %v", sinks)
	}
}

func TestMatch_MultipleRules(t *testing.T) {
	f, err := filter.New(map[string]string{
		"errors":   `(?i)error`,
		"critical": `(?i)critical`,
	})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	sinks := f.Match("CRITICAL ERROR: disk full")
	if len(sinks) != 2 {
		t.Errorf("expected 2 matches, got %v", sinks)
	}
}

func TestMatch_EmptyFilter(t *testing.T) {
	f, err := filter.New(map[string]string{})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	sinks := f.Match("any line at all")
	if len(sinks) != 0 {
		t.Errorf("expected no matches on empty filter, got %v", sinks)
	}
}
