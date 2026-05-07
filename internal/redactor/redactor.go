// Package redactor masks sensitive values in diff results before output.
// Keys matching known secret patterns (e.g. PASSWORD, SECRET, TOKEN, KEY)
// have their values replaced with a redaction placeholder.
package redactor

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// DefaultPlaceholder is the string used to replace sensitive values.
const DefaultPlaceholder = "***REDACTED***"

// Options controls redactor behaviour.
type Options struct {
	// Placeholder overrides the default redaction string.
	Placeholder string
	// ExtraPatterns are additional case-insensitive substrings that
	// trigger redaction when found in a key name.
	ExtraPatterns []string
}

var builtinPatterns = []string{
	"PASSWORD", "PASSWD", "SECRET", "TOKEN", "API_KEY",
	"PRIVATE_KEY", "ACCESS_KEY", "AUTH", "CREDENTIAL",
}

// Apply returns a copy of results with sensitive values masked.
// Original results are not modified.
func Apply(results []diff.Result, opts Options) []diff.Result {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = DefaultPlaceholder
	}

	patterns := make([]string, len(builtinPatterns))
	copy(patterns, builtinPatterns)
	for _, p := range opts.ExtraPatterns {
		patterns = append(patterns, strings.ToUpper(p))
	}

	out := make([]diff.Result, len(results))
	for i, r := range results {
		if isSensitive(r.Key, patterns) {
			r.Value = placeholder
			r.OtherValues = redactMap(r.OtherValues, placeholder)
		}
		out[i] = r
	}
	return out
}

// IsSensitiveKey reports whether the given key name would be redacted
// using the provided options. This is useful for callers that need to
// check sensitivity without running a full Apply pass.
func IsSensitiveKey(key string, opts Options) bool {
	patterns := make([]string, len(builtinPatterns))
	copy(patterns, builtinPatterns)
	for _, p := range opts.ExtraPatterns {
		patterns = append(patterns, strings.ToUpper(p))
	}
	return isSensitive(key, patterns)
}

// isSensitive returns true when key contains any of the given patterns.
func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// redactMap returns a new map with all values replaced by placeholder.
func redactMap(m map[string]string, placeholder string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k := range m {
		out[k] = placeholder
	}
	return out
}
