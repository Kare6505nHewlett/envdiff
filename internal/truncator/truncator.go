// Package truncator shortens long env values in diff results for display purposes.
// It is useful when rendering output in terminals or narrow report formats.
package truncator

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

const defaultMaxLen = 64
const defaultEllipsis = "..."

// Options controls truncation behaviour.
type Options struct {
	MaxLen   int
	Ellipsis string
}

// DefaultOptions returns sensible truncation defaults.
func DefaultOptions() Options {
	return Options{
		MaxLen:   defaultMaxLen,
		Ellipsis: defaultEllipsis,
	}
}

// Apply returns a new slice of results with values truncated according to opts.
// Original results are not mutated.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if opts.MaxLen <= 0 {
		opts.MaxLen = defaultMaxLen
	}
	if opts.Ellipsis == "" {
		opts.Ellipsis = defaultEllipsis
	}

	out := make([]diff.Result, len(results))
	for i, r := range results {
		copy := r
		copy.Values = make(map[string]string, len(r.Values))
		for file, val := range r.Values {
			copy.Values[file] = truncate(val, opts.MaxLen, opts.Ellipsis)
		}
		out[i] = copy
	}
	return out
}

func truncate(s string, maxLen int, ellipsis string) string {
	// Normalise newlines for display
	s = strings.ReplaceAll(s, "\n", "\\n")
	if len(s) <= maxLen {
		return s
	}
	cutAt := maxLen - len(ellipsis)
	if cutAt < 0 {
		cutAt = 0
	}
	return s[:cutAt] + ellipsis
}
