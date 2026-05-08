// Package classifier categorises diff results into named tiers based on
// configurable value patterns (e.g. secrets, URLs, feature flags).
package classifier

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Category represents a named classification tier.
type Category string

const (
	CategorySecret  Category = "secret"
	CategoryURL     Category = "url"
	CategoryFlag    Category = "flag"
	CategoryNumeric Category = "numeric"
	CategoryUnknown Category = "unknown"
)

// Result wraps a diff.Result with its assigned category.
type Result struct {
	diff.Result
	Category Category
}

var (
	secretKeyRe  = regexp.MustCompile(`(?i)(secret|password|passwd|token|apikey|api_key|private)`)
	urlValueRe   = regexp.MustCompile(`(?i)^https?://`)
	numericRe    = regexp.MustCompile(`^-?\d+(\.\d+)?$`)
	flagValueRe  = regexp.MustCompile(`(?i)^(true|false|yes|no|on|off|1|0)$`)
)

// Apply classifies each diff result and returns a slice of classified results.
func Apply(results []diff.Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Result:   r,
			Category: classify(r.Key, r.Value),
		})
	}
	return out
}

// classify determines the category for a given key/value pair.
func classify(key, value string) Category {
	if secretKeyRe.MatchString(key) {
		return CategorySecret
	}
	v := strings.TrimSpace(value)
	if urlValueRe.MatchString(v) {
		return CategoryURL
	}
	if flagValueRe.MatchString(v) {
		return CategoryFlag
	}
	if numericRe.MatchString(v) {
		return CategoryNumeric
	}
	return CategoryUnknown
}
