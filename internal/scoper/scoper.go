// Package scoper filters and organizes diff results by environment scope.
// A scope is a named boundary (e.g. "production", "staging") that groups
// results from specific files and applies targeted filtering.
package scoper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Scope defines a named grouping of env files.
type Scope struct {
	Name  string
	Files []string
}

// Result holds the results associated with a named scope.
type Result struct {
	Scope   string
	Entries []diff.Result
}

// Apply partitions the given results into scopes based on file membership.
// Results whose file does not match any scope are placed in a "default" scope.
func Apply(results []diff.Result, scopes []Scope) []Result {
	index := buildIndex(scopes)
	buckets := make(map[string][]diff.Result)

	for _, r := range results {
		scope := resolveScope(r.File, index)
		buckets[scope] = append(buckets[scope], r)
	}

	var out []Result
	for _, s := range scopes {
		if entries, ok := buckets[s.Name]; ok {
			out = append(out, Result{Scope: s.Name, Entries: entries})
		}
	}
	if entries, ok := buckets["default"]; ok {
		out = append(out, Result{Scope: "default", Entries: entries})
	}
	return out
}

// Names returns a sorted list of unique scope names from the given results.
func Names(results []Result) []string {
	seen := make(map[string]struct{})
	for _, r := range results {
		seen[r.Scope] = struct{}{}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// buildIndex maps each file path to its scope name.
func buildIndex(scopes []Scope) map[string]string {
	idx := make(map[string]string)
	for _, s := range scopes {
		for _, f := range s.Files {
			idx[strings.ToLower(f)] = s.Name
		}
	}
	return idx
}

func resolveScope(file string, index map[string]string) string {
	if name, ok := index[strings.ToLower(file)]; ok {
		return name
	}
	return "default"
}
