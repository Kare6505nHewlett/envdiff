package diff

import "sort"

// Status represents the comparison outcome for a single key.
type Status string

const (
	StatusMatch    Status = "match"
	StatusMismatch Status = "mismatch"
	StatusMissing  Status = "missing" // present in base, absent in target
	StatusExtra    Status = "extra"   // absent in base, present in target
)

// Result holds the comparison result for a single environment key.
type Result struct {
	Key        string `json:"key"`
	Status     Status `json:"status"`
	BaseValue  string `json:"base_value,omitempty"`
	OtherValue string `json:"other_value,omitempty"`
}

// Compare compares two parsed env maps and returns a sorted list of Results.
func Compare(base, target map[string]string) []Result {
	var results []Result

	for key, baseVal := range base {
		if targetVal, ok := target[key]; !ok {
			results = append(results, Result{
				Key:       key,
				Status:    StatusMissing,
				BaseValue: baseVal,
			})
		} else if hasMismatch(baseVal, targetVal) {
			results = append(results, Result{
				Key:        key,
				Status:     StatusMismatch,
				BaseValue:  baseVal,
				OtherValue: targetVal,
			})
		} else {
			results = append(results, Result{
				Key:    key,
				Status: StatusMatch,
			})
		}
	}

	for key, targetVal := range target {
		if _, ok := base[key]; !ok {
			results = append(results, Result{
				Key:        key,
				Status:     StatusExtra,
				OtherValue: targetVal,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

func hasMismatch(a, b string) bool {
	return a != b
}
