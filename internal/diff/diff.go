package diff

import (
	"sort"
)

// Status represents the type of difference found for a key.
type Status string

const (
	StatusMissing  Status = "missing"
	StatusExtra    Status = "extra"
	StatusMismatch Status = "mismatch"
)

// Result holds the comparison outcome for a single key.
type Result struct {
	Key    string   `json:"key"`
	Status Status   `json:"status"`
	Files  []string `json:"files"`
}

// Compare takes a map of filename -> key/value pairs and returns a sorted
// slice of Results describing any missing, extra, or mismatched keys.
func Compare(envs map[string]map[string]string) []Result {
	// Collect all unique keys across all files.
	keySet := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	files := make([]string, 0, len(envs))
	for f := range envs {
		files = append(files, f)
	}
	sort.Strings(files)

	var results []Result

	for key := range keySet {
		var presentIn []string
		var missingIn []string
		values := map[string]string{}

		for _, f := range files {
			if v, ok := envs[f][key]; ok {
				presentIn = append(presentIn, f)
				values[f] = v
			} else {
				missingIn = append(missingIn, f)
			}
		}

		switch {
		case len(missingIn) == len(files):
			// Key exists in no file — shouldn't happen, skip.
			continue
		case len(missingIn) > 0 && len(presentIn) == 1:
			results = append(results, Result{Key: key, Status: StatusExtra, Files: presentIn})
		case len(missingIn) > 0:
			results = append(results, Result{Key: key, Status: StatusMissing, Files: missingIn})
		default:
			// All files have the key — check for value mismatch.
			if hasMismatch(values) {
				results = append(results, Result{Key: key, Status: StatusMismatch, Files: files})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

func hasMismatch(values map[string]string) bool {
	var first string
	set := false
	for _, v := range values {
		if !set {
			first = v
			set = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}
