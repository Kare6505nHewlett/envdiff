// Package merger provides functionality to merge multiple .env file
// comparison results into a unified canonical view.
package merger

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// MergedKey represents a key seen across all compared environments,
// with its value per file and an overall status.
type MergedKey struct {
	Key    string
	Values map[string]string // file path -> value (empty string if missing)
	Status string            // "match", "missing", "mismatch"
}

// Merge takes a slice of diff.Result and collapses them into a slice of
// MergedKey entries, one per unique key across all files.
func Merge(results []diff.Result) []MergedKey {
	type keyFile struct {
		value  string
		present bool
	}

	// key -> file -> {value, present}
	index := map[string]map[string]keyFile{}

	for _, r := range results {
		if _, ok := index[r.Key]; !ok {
			index[r.Key] = map[string]keyFile{}
		}
		index[r.Key][r.File] = keyFile{value: r.Value, present: r.Status != "missing"}
	}

	merged := make([]MergedKey, 0, len(index))
	for key, files := range index {
		mk := MergedKey{
			Key:    key,
			Values: make(map[string]string, len(files)),
		}
		for file, kf := range files {
			if kf.present {
				mk.Values[file] = kf.value
			} else {
				mk.Values[file] = ""
			}
		}
		mk.Status = computeStatus(files)
		merged = append(merged, mk)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Key < merged[j].Key
	})
	return merged
}

// computeStatus derives an overall status from per-file presence and values.
func computeStatus(files map[string]keyFile) string {
	var presentValues []string
	for _, kf := range files {
		if !kf.present {
			return "missing"
		}
		presentValues = append(presentValues, kf.value)
	}
	if len(presentValues) == 0 {
		return "missing"
	}
	first := presentValues[0]
	for _, v := range presentValues[1:] {
		if v != first {
			return "mismatch"
		}
	}
	return "match"
}
