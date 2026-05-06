// Package differ provides line-level diffing of individual key values
// across two .env files, producing human-readable change hunks.
package differ

import (
	"fmt"
	"strings"
)

// Hunk represents a single value change for a key between two files.
type Hunk struct {
	Key      string
	FileA    string
	FileB    string
	ValueA   string
	ValueB   string
	Removed  bool // key only in FileA
	Added    bool // key only in FileB
	Changed  bool // key in both, values differ
}

// String returns a compact diff-style representation of the hunk.
func (h Hunk) String() string {
	switch {
	case h.Removed:
		return fmt.Sprintf("- [%s] %s=%s", h.FileA, h.Key, h.ValueA)
	case h.Added:
		return fmt.Sprintf("+ [%s] %s=%s", h.FileB, h.Key, h.ValueB)
	default:
		return fmt.Sprintf("~ %s\n  - [%s] %s\n  + [%s] %s",
			h.Key, h.FileA, h.ValueA, h.FileB, h.ValueB)
	}
}

// Diff computes line-level hunks between two parsed env maps.
// fileA and fileB are display names (e.g. file paths).
func Diff(fileA, fileB string, mapA, mapB map[string]string) []Hunk {
	var hunks []Hunk

	for k, va := range mapA {
		vb, ok := mapB[k]
		if !ok {
			hunks = append(hunks, Hunk{Key: k, FileA: fileA, ValueA: va, Removed: true})
			continue
		}
		if va != vb {
			hunks = append(hunks, Hunk{
				Key: k, FileA: fileA, FileB: fileB,
				ValueA: va, ValueB: vb, Changed: true,
			})
		}
	}

	for k, vb := range mapB {
		if _, ok := mapA[k]; !ok {
			hunks = append(hunks, Hunk{Key: k, FileB: fileB, ValueB: vb, Added: true})
		}
	}

	sortHunks(hunks)
	return hunks
}

// Summary returns a brief human-readable summary line.
func Summary(hunks []Hunk) string {
	var added, removed, changed int
	for _, h := range hunks {
		switch {
		case h.Added:
			added++
		case h.Removed:
			removed++
		case h.Changed:
			changed++
		}
	}
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}

func sortHunks(hunks []Hunk) {
	// insertion sort by Key for deterministic output
	for i := 1; i < len(hunks); i++ {
		for j := i; j > 0 && hunks[j].Key < hunks[j-1].Key; j-- {
			hunks[j], hunks[j-1] = hunks[j-1], hunks[j]
		}
	}
}
