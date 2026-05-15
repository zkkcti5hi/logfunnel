package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Source defines a log source to tail.
type Source struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// Sink defines a destination for matched log entries.
type Sink struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"` // "file" or "stdout"
	Path string `yaml:"path,omitempty"`
}

// Rule maps a regex filter to a sink.
type Rule struct {
	Pattern string `yaml:"pattern"`
	Sink    string `yaml:"sink"`
}

// Config is the top-level configuration structure.
type Config struct {
	Sources []Source `yaml:"sources"`
	Sinks   []Sink   `yaml:"sinks"`
	Rules   []Rule   `yaml:"rules"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks that the configuration is semantically valid.
func (c *Config) Validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source is required")
	}
	if len(c.Sinks) == 0 {
		return fmt.Errorf("at least one sink is required")
	}

	sinkNames := make(map[string]struct{}, len(c.Sinks))
	for _, s := range c.Sinks {
		if s.Name == "" {
			return fmt.Errorf("sink name must not be empty")
		}
		if s.Type != "file" && s.Type != "stdout" {
			return fmt.Errorf("sink %q has unsupported type %q", s.Name, s.Type)
		}
		if s.Type == "file" && s.Path == "" {
			return fmt.Errorf("sink %q of type 'file' requires a path", s.Name)
		}
		sinkNames[s.Name] = struct{}{}
	}

	for _, r := range c.Rules {
		if r.Pattern == "" {
			return fmt.Errorf("rule pattern must not be empty")
		}
		if _, ok := sinkNames[r.Sink]; !ok {
			return fmt.Errorf("rule references unknown sink %q", r.Sink)
		}
	}

	return nil
}
