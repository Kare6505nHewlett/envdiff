// Package deduplicator removes duplicate diff results, keeping the most
// significant entry when the same key appears multiple times across files.
package deduplicator

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// statusPriority assigns a numeric priority to each diff status.
// Higher values indicate more significant findings.
var statusPriority = map[string]int{
	"mismatch": 3,
	"missing":  2,
	"extra":    1,
	"match":    0,
}

// Options controls deduplication behaviour.
type Options struct {
	// KeepAll retains all entries for a key, only removing exact duplicates.
	KeepAll bool
}

// DefaultOptions returns sensible defaults: keep the highest-priority entry
// per key across all files.
func DefaultOptions() Options {
	return Options{KeepAll: false}
}

// Apply deduplicates results according to opts.
// When KeepAll is false, only the highest-priority result per key is kept.
// When KeepAll is true, only exact duplicates (same key + file + status + value) are removed.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if len(results) == 0 {
		return []diff.Result{}
	}

	if opts.KeepAll {
		return deduplicateExact(results)
	}
	return deduplicateByPriority(results)
}

// deduplicateByPriority keeps one result per key — the one with the highest
// status priority. Ties are broken by file name (alphabetical first).
func deduplicateByPriority(results []diff.Result) []diff.Result {
	best := make(map[string]diff.Result)

	for _, r := range results {
		existing, seen := best[r.Key]
		if !seen {
			best[r.Key] = r
			continue
		}
		ep := statusPriority[existing.Status]
		cp := statusPriority[r.Status]
		if cp > ep || (cp == ep && r.File < existing.File) {
			best[r.Key] = r
		}
	}

	out := make([]diff.Result, 0, len(best))
	for _, r := range best {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// deduplicateExact removes entries that are identical in every field.
func deduplicateExact(results []diff.Result) []diff.Result {
	seen := make(map[string]struct{})
	out := make([]diff.Result, 0, len(results))

	for _, r := range results {
		key := r.Key + "\x00" + r.File + "\x00" + r.Status + "\x00" + r.Value
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, r)
	}
	return out
}
