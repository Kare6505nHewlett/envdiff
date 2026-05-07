// Package masker provides value masking for .env entries based on
// configurable patterns. Unlike the redactor (which replaces sensitive
// values with a placeholder), the masker partially reveals values so
// that engineers can confirm the correct value is present without
// exposing the full secret (e.g. "abc***xyz").
package masker

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// DefaultOptions returns a MaskOptions with sensible defaults.
func DefaultOptions() Options {
	return Options{
		PrefixLen:    3,
		SuffixLen:    3,
		MaskChar:     '*',
		MinMaskLen:   3,
		SensitiveKeys: []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"},
	}
}

// Apply returns a new slice of diff.Result where values of sensitive keys
// are partially masked. The original slice is never modified.
func Apply(results []diff.Result, opts Options) []diff.Result {
	out := make([]diff.Result, len(results))
	for i, r := range results {
		copy := r
		if isSensitive(r.Key, opts.SensitiveKeys) {
			copy.Values = maskMap(r.Values, opts)
		}
		out[i] = copy
	}
	return out
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func maskMap(values map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(values))
	for file, val := range values {
		out[file] = maskValue(val, opts)
	}
	return out
}

// maskValue partially reveals a string value.
func maskValue(val string, opts Options) string {
	if val == "" {
		return val
	}
	runes := []rune(val)
	n := len(runes)
	prefix := opts.PrefixLen
	suffix := opts.SuffixLen
	// If the value is too short to mask meaningfully, replace entirely.
	if n <= prefix+suffix {
		return strings.Repeat(string(opts.MaskChar), max(n, opts.MinMaskLen))
	}
	midLen := n - prefix - suffix
	if midLen < opts.MinMaskLen {
		midLen = opts.MinMaskLen
	}
	mask := strings.Repeat(string(opts.MaskChar), midLen)
	return string(runes[:prefix]) + mask + string(runes[n-suffix:])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
