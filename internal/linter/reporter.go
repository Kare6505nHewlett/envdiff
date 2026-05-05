package linter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ReportText writes lint findings in a human-readable format to w.
func ReportText(w io.Writer, findings []Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "lint: no issues found")
		return
	}

	sorted := make([]Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].File != sorted[j].File {
			return sorted[i].File < sorted[j].File
		}
		return sorted[i].Key < sorted[j].Key
	})

	for _, f := range sorted {
		fmt.Fprintf(w, "[%s] %s (%s): %s\n", f.Severity, f.File, f.Rule, f.Message)
	}
}

// ReportJSON writes lint findings as a JSON array to w.
func ReportJSON(w io.Writer, findings []Finding) error {
	if findings == nil {
		findings = []Finding{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(findings)
}

// CountBySeverity returns a map of severity -> count.
func CountBySeverity(findings []Finding) map[Severity]int {
	counts := make(map[Severity]int)
	for _, f := range findings {
		counts[f.Severity]++
	}
	return counts
}
