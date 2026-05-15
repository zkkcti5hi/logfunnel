// Package filter provides regex-based log line routing.
//
// A Filter is constructed from a map of sink names to regular expression
// patterns. Each incoming log line is tested against every compiled rule;
// all matching sink names are returned so that callers can forward the
// line to the appropriate output sinks.
//
// Example usage:
//
//	f, err := filter.New(map[string]string{
//		"errors":   `(?i)error`,
//		"warnings": `(?i)warn`,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	sinks := f.Match(line)
//	for _, sink := range sinks {
//		// route line to sink
//	}
package filter
