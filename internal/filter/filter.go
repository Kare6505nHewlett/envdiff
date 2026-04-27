package filter

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options holds filtering criteria for diff results.
type Options struct {
	OnlyMissing  bool
	OnlyMismatch bool
	Prefix       string
	ExcludeKeys  []string
}

// Apply filters a slice of diff.Result according to the given Options.
func Apply(results []diff.Result, opts Options) []diff.Result {
	var filtered []diff.Result
	for _, r := range results {
		if opts.OnlyMissing && r.Status != diff.StatusMissing {
			continue
		}
		if opts.OnlyMismatch && r.Status != diff.StatusMismatch {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(r.Key, opts.Prefix) {
			continue
		}
		if isExcluded(r.Key, opts.ExcludeKeys) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// isExcluded returns true if key matches any entry in the exclusion list.
func isExcluded(key string, excludeKeys []string) bool {
	for _, ex := range excludeKeys {
		if strings.EqualFold(key, ex) {
			return true
		}
	}
	return false
}
