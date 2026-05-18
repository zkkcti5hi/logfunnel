package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExporter_EmptyRegistry(t *testing.T) {
	reg := NewRegistry()
	exporter := NewExporter(reg)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	exporter.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]int64
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty snapshot, got %v", result)
	}
}

func TestExporter_WithCounters(t *testing.T) {
	reg := NewRegistry()
	reg.GetOrCreate("lines_read").Add(42)
	reg.GetOrCreate("lines_routed").Add(7)

	exporter := NewExporter(reg)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	exporter.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result map[string]int64
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["lines_read"] != 42 {
		t.Errorf("expected lines_read=42, got %d", result["lines_read"])
	}
	if result["lines_routed"] != 7 {
		t.Errorf("expected lines_routed=7, got %d", result["lines_routed"])
	}
}

func TestExporter_ContentTypeHeader(t *testing.T) {
	reg := NewRegistry()
	exporter := NewExporter(reg)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	exporter.Handler().ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
