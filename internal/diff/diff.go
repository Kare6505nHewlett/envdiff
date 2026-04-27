package diff

import "sort"

// KeyStatus represents the comparison status of a key across environments.
type KeyStatus string

const (
	StatusMissing   KeyStatus = "missing"   // key exists in reference but not in target
	StatusExtra     KeyStatus = "extra"     // key exists in target but not in reference
	StatusMismatch  KeyStatus = "mismatch"  // key exists in both but values differ
	StatusMatch     KeyStatus = "match"     // key exists in both with identical values
)

// Entry describes the diff result for a single key.
type Entry struct {
	Key        string
	Status     KeyStatus
	RefValue   string // value from the reference env file
	TargetValue string // value from the target env file
}

// Result holds the full diff between two env files.
type Result struct {
	Reference string  // label / path of the reference file
	Target    string  // label / path of the target file
	Entries   []Entry
}

// Compare compares two parsed env maps and returns a Result.
// refLabel and targetLabel are used for display purposes (e.g. file paths).
func Compare(refLabel string, ref map[string]string, targetLabel string, target map[string]string) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, rv := range ref {
		seen[k] = true
		if tv, ok := target[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: StatusMissing, RefValue: rv})
		} else if rv != tv {
			entries = append(entries, Entry{Key: k, Status: StatusMismatch, RefValue: rv, TargetValue: tv})
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusMatch, RefValue: rv, TargetValue: tv})
		}
	}

	for k, tv := range target {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: StatusExtra, TargetValue: tv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Result{
		Reference: refLabel,
		Target:    targetLabel,
		Entries:   entries,
	}
}

// Summary returns counts of each status in the result.
func (r *Result) Summary() map[KeyStatus]int {
	counts := map[KeyStatus]int{}
	for _, e := range r.Entries {
		counts[e.Status]++
	}
	return counts
}
