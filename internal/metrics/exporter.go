package metrics

import (
	"encoding/json"
	"net/http"
)

// Exporter exposes a Registry's metrics over HTTP as JSON.
type Exporter struct {
	registry *Registry
	mux      *http.ServeMux
}

// NewExporter creates an Exporter backed by the given Registry.
// It registers a handler at /metrics on a new ServeMux.
func NewExporter(r *Registry) *Exporter {
	e := &Exporter{
		registry: r,
		mux:      http.NewServeMux(),
	}
	e.mux.HandleFunc("/metrics", e.handleMetrics)
	return e
}

// Handler returns the HTTP handler for use with http.ListenAndServe.
func (e *Exporter) Handler() http.Handler {
	return e.mux
}

// handleMetrics writes a JSON snapshot of all counters to the response.
func (e *Exporter) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snapshot := e.registry.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
	}
}
