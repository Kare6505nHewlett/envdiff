package scorer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Result holds the computed health score and a human-readable grade.
type Result struct {
	Score   float64 `json:"score"`
	Grade   string  `json:"grade"`
	Total   int     `json:"total"`
	Missing int     `json:"missing"`
	Mismatch int    `json:"mismatch"`
}

// gradeFor maps a score to a letter grade.
func gradeFor(score float64) string {
	switch {
	case score >= 95:
		return "A"
	case score >= 80:
		return "B"
	case score >= 65:
		return "C"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}

// ReportText writes a human-readable score summary to w.
func ReportText(w io.Writer, r Result) {
	fmt.Fprintf(w, "Health Score: %.1f/100 (%s)\n", r.Score, r.Grade)
	fmt.Fprintf(w, "  Total keys : %d\n", r.Total)
	fmt.Fprintf(w, "  Missing    : %d\n", r.Missing)
	fmt.Fprintf(w, "  Mismatched : %d\n", r.Mismatch)
}

// ReportJSON writes a JSON-encoded score result to w.
func ReportJSON(w io.Writer, r Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
