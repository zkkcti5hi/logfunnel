// Package metrics provides a minimal, dependency-free metrics layer for
// logfunnel.
//
// # Overview
//
// A [Registry] owns a set of named [Counter] values. Counters are created on
// first access via [Registry.Counter] and are safe for concurrent use.
//
// # Usage
//
//	reg := metrics.NewRegistry()
//
//	// Increment counters wherever log lines are processed.
//	reg.Counter("lines_in").Inc()
//	reg.Counter("lines_routed").Inc()
//	reg.Counter("lines_dropped").Inc()
//
//	// Dump a human-readable summary to stdout at shutdown.
//	reg.WriteTo(os.Stdout)
//
// Counters are intentionally simple — there is no histogram or gauge type.
// If richer instrumentation is required, replace this package with a
// Prometheus client.
package metrics
