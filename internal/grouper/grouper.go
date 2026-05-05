// Package grouper organises diff results by a chosen grouping dimension,
// such as file name or key prefix, returning a map of labelled slices.
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// GroupBy controls how results are grouped.
type GroupBy string

const (
	// ByFile groups each result under the file it was found in.
	ByFile GroupBy = "file"
	// ByPrefix groups results by the segment of the key before the first
	// underscore (e.g. "DB" for "DB_HOST").
	ByPrefix GroupBy = "prefix"
	// ByStatus groups results by their diff status string.
	ByStatus GroupBy = "status"
)

// Group holds a labelled collection of diff results.
type Group struct {
	Label   string
	Results []diff.Result
}

// Apply partitions results according to the requested dimension and returns
// groups sorted alphabetically by label.
func Apply(results []diff.Result, by GroupBy) []Group {
	buckets := make(map[string][]diff.Result)

	for _, r := range results {
		key := labelFor(r, by)
		buckets[key] = append(buckets[key], r)
	}

	labels := make([]string, 0, len(buckets))
	for l := range buckets {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	groups := make([]Group, 0, len(labels))
	for _, l := range labels {
		groups = append(groups, Group{Label: l, Results: buckets[l]})
	}
	return groups
}

func labelFor(r diff.Result, by GroupBy) string {
	switch by {
	case ByFile:
		if r.File != "" {
			return r.File
		}
		return "(unknown)"
	case ByPrefix:
		parts := strings.SplitN(r.Key, "_", 2)
		if len(parts) > 1 && parts[0] != "" {
			return parts[0]
		}
		return "(no prefix)"
	case ByStatus:
		if r.Status != "" {
			return r.Status
		}
		return "(unknown)"
	default:
		return "(unknown)"
	}
}
