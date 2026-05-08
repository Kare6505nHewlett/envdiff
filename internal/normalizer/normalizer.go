// Package normalizer provides utilities for normalizing .env key-value pairs
// before comparison, such as trimming whitespace, normalizing case, and
// standardizing boolean-like values.
package normalizer

import (
	"strings"
)

// Options controls which normalization steps are applied.
type Options struct {
	TrimSpace       bool
	LowercaseKeys   bool
	NormalizeBools  bool
	CollapseEmpty   bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:      true,
		LowercaseKeys:  false,
		NormalizeBools: true,
		CollapseEmpty:  false,
	}
}

// Entry represents a single normalized key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Apply normalizes a map of env key-value pairs according to the given options.
// It returns a new map and does not mutate the input.
func Apply(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		nk := normalizeKey(k, opts)
		nv := normalizeValue(v, opts)
		out[nk] = nv
	}
	return out
}

func normalizeKey(k string, opts Options) string {
	if opts.TrimSpace {
		k = strings.TrimSpace(k)
	}
	if opts.LowercaseKeys {
		k = strings.ToLower(k)
	}
	return k
}

func normalizeValue(v string, opts Options) string {
	if opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if opts.CollapseEmpty && v == "" {
		return ""
	}
	if opts.NormalizeBools {
		v = normalizeBool(v)
	}
	return v
}

// normalizeBool maps common boolean-like strings to canonical forms.
func normalizeBool(v string) string {
	switch strings.ToLower(v) {
	case "true", "1", "yes", "on":
		return "true"
	case "false", "0", "no", "off":
		return "false"
	}
	return v
}
