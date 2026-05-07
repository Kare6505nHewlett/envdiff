// Package profiler analyses a set of diff results and produces a health
// profile summarising the overall quality of the environment configuration.
package profiler

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Grade represents a letter-grade health assessment.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// Profile holds the computed health profile for a set of results.
type Profile struct {
	Grade        Grade
	Score        float64 // 0–100
	TotalKeys    int
	MissingCount int
	MismatchCount int
	MatchCount   int
	Summary      string
}

// Compute derives a Profile from a slice of diff results.
func Compute(results []diff.Result) Profile {
	p := Profile{TotalKeys: len(results)}
	for _, r := range results {
		switch r.Status {
		case diff.StatusMatch:
			p.MatchCount++
		case diff.StatusMissing:
			p.MissingCount++
		case diff.StatusMismatch:
			p.MismatchCount++
		}
	}

	if p.TotalKeys == 0 {
		p.Score = 100
	} else {
		penalty := float64(p.MissingCount*2+p.MismatchCount) / float64(p.TotalKeys*2) * 100
		p.Score = max(0, 100-penalty)
	}

	p.Grade = gradeFor(p.Score)
	p.Summary = buildSummary(p)
	return p
}

// IsHealthy reports whether the profile meets a minimum acceptable threshold,
// defined as a score of 80 or above (grade B or better).
func (p Profile) IsHealthy() bool {
	return p.Score >= 80
}

func gradeFor(score float64) Grade {
	switch {
	case score >= 95:
		return GradeA
	case score >= 80:
		return GradeB
	case score >= 65:
		return GradeC
	case score >= 50:
		return GradeD
	default:
		return GradeF
	}
}

func buildSummary(p Profile) string {
	parts := []string{
		fmt.Sprintf("score=%.1f", p.Score),
		fmt.Sprintf("grade=%s", p.Grade),
		fmt.Sprintf("match=%d", p.MatchCount),
		fmt.Sprintf("missing=%d", p.MissingCount),
		fmt.Sprintf("mismatch=%d", p.MismatchCount),
	}
	return strings.Join(parts, " ")
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
