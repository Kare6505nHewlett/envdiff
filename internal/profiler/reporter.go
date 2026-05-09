package profiler

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable profile report to w.
func ReportText(w io.Writer, p Profile) {
	fmt.Fprintf(w, "Environment Health Profile\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 30))
	fmt.Fprintf(w, "Grade   : %s\n", p.Grade)
	fmt.Fprintf(w, "Score   : %.1f / 100\n", p.Score)
	fmt.Fprintf(w, "Total   : %d keys\n", p.TotalKeys)
	fmt.Fprintf(w, "Match   : %d\n", p.MatchCount)
	fmt.Fprintf(w, "Missing : %d\n", p.MissingCount)
	fmt.Fprintf(w, "Mismatch: %d\n", p.MismatchCount)
	if p.Summary != "" {
		fmt.Fprintf(w, "Summary : %s\n", p.Summary)
	}
}

// jsonProfile is the JSON-serialisable shape of a Profile.
type jsonProfile struct {
	Grade         string  `json:"grade"`
	Score         float64 `json:"score"`
	TotalKeys     int     `json:"total_keys"`
	MatchCount    int     `json:"match"`
	MissingCount  int     `json:"missing"`
	MismatchCount int     `json:"mismatch"`
	Summary       string  `json:"summary"`
}

// ReportJSON writes a JSON-encoded profile to w.
func ReportJSON(w io.Writer, p Profile) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonProfile{
		Grade:         string(p.Grade),
		Score:         p.Score,
		TotalKeys:     p.TotalKeys,
		MatchCount:    p.MatchCount,
		MissingCount:  p.MissingCount,
		MismatchCount: p.MismatchCount,
		Summary:       p.Summary,
	})
}
