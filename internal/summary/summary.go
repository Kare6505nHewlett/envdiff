package summary

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Stats holds aggregate counts from a diff result set.
type Stats struct {
	Total    int
	Matched  int
	Missing  int
	Mismatch int
	Files    []string
}

// Compute calculates statistics from a slice of diff results.
func Compute(results []diff.Result) Stats {
	fileSet := make(map[string]struct{})
	stats := Stats{}

	for _, r := range results {
		stats.Total++
		switch r.Status {
		case diff.StatusMatch:
			stats.Matched++
		case diff.StatusMissing:
			stats.Missing++
		case diff.StatusMismatch:
			stats.Mismatch++
		}
		fileSet[r.File] = struct{}{}
	}

	for f := range fileSet {
		stats.Files = append(stats.Files, f)
	}
	return stats
}

// Print writes a human-readable summary to the provided writer.
func Print(w io.Writer, stats Stats) {
	fmt.Fprintf(w, "Summary\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 30))
	fmt.Fprintf(w, "Files compared : %d\n", len(stats.Files))
	fmt.Fprintf(w, "Total keys     : %d\n", stats.Total)
	fmt.Fprintf(w, "Matched        : %d\n", stats.Matched)
	fmt.Fprintf(w, "Missing        : %d\n", stats.Missing)
	fmt.Fprintf(w, "Mismatched     : %d\n", stats.Mismatch)

	if stats.Missing > 0 || stats.Mismatch > 0 {
		fmt.Fprintf(w, "Status         : DIFFERENCES FOUND\n")
	} else {
		fmt.Fprintf(w, "Status         : ALL MATCHED\n")
	}
}
