package diff

import "sort"

// Status represents the comparison outcome for a single key.
type Status string

const (
	StatusMatch    Status = "match"
	StatusMissing  Status = "missing"
	StatusMismatch Status = "mismatch"
)

// Result holds the comparison result for one key in one file.
type Result struct {
	Key      string
	File     string
	Status   Status
	BaseVal  string
	OtherVal string
}

// Compare compares a base env map against one or more target env maps.
// Each target is identified by its filename key in the targets map.
func Compare(base map[string]string, targets map[string]map[string]string) []Result {
	var results []Result

	keys := make([]string, 0, len(base))
	for k := range base {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, file := range sortedKeys(targets) {
		target := targets[file]
		for _, key := range keys {
			baseVal := base[key]
			otherVal, exists := target[key]
			switch {
			case !exists:
				results = append(results, Result{
					Key: key, File: file,
					Status: StatusMissing,
					BaseVal: baseVal,
				})
			case hasMismatch(baseVal, otherVal):
				results = append(results, Result{
					Key: key, File: file,
					Status: StatusMismatch,
					BaseVal: baseVal, OtherVal: otherVal,
				})
			default:
				results = append(results, Result{
					Key: key, File: file,
					Status: StatusMatch,
					BaseVal: baseVal, OtherVal: otherVal,
				})
			}
		}
	}
	return results
}

func hasMismatch(a, b string) bool {
	return a != b
}

func sortedKeys(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
