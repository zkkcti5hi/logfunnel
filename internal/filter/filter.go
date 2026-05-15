package filter

import (
	"fmt"
	"regexp"
)

// Rule represents a single routing rule that matches log lines
// against a regex pattern and routes them to a named sink.
type Rule struct {
	Pattern *regexp.Regexp
	SinkName string
}

// Filter holds compiled routing rules and matches log lines to sinks.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a map of sink name -> regex pattern strings.
// Returns an error if any pattern fails to compile.
func New(rules map[string]string) (*Filter, error) {
	compiled := make([]Rule, 0, len(rules))
	for sink, pattern := range rules {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid pattern for sink %q: %w", sink, err)
		}
		compiled = append(compiled, Rule{
			Pattern:  re,
			SinkName: sink,
		})
	}
	return &Filter{rules: compiled}, nil
}

// Match returns the names of all sinks whose pattern matches line.
// A line may match multiple rules.
func (f *Filter) Match(line string) []string {
	var matched []string
	for _, r := range f.rules {
		if r.Pattern.MatchString(line) {
			matched = append(matched, r.SinkName)
		}
	}
	return matched
}

// Rules returns the compiled rules held by the filter.
func (f *Filter) Rules() []Rule {
	return f.rules
}
