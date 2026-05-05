// Package scorer assigns a numeric health score to a set of diff results,
// giving callers a single at-a-glance metric for environment parity.
package scorer

import (
	"math"

	"github.com/user/envdiff/internal/diff"
)

// Weights control how much each result status penalises the score.
type Weights struct {
	Missing  float64
	Mismatch float64
}

// DefaultWeights are used when nil is passed to Score.
var DefaultWeights = Weights{
	Missing:  1.0,
	Mismatch: 0.5,
}

// Result holds the computed score and contributing counts.
type Result struct {
	Score    float64 // 0.0 (worst) – 100.0 (perfect)
	Total    int
	Missing  int
	Mismatch int
	Matched  int
}

// Score computes a parity health score from a slice of diff results.
// A nil weights pointer falls back to DefaultWeights.
func Score(results []diff.Result, w *Weights) Result {
	if w == nil {
		w = &DefaultWeights
	}

	var missing, mismatch, matched int
	for _, r := range results {
		switch r.Status {
		case diff.StatusMissing:
			missing++
		case diff.StatusMismatch:
			mismatch++
		case diff.StatusMatch:
			matched++
		}
	}

	total := missing + mismatch + matched
	if total == 0 {
		return Result{Score: 100.0}
	}

	penalty := (float64(missing)*w.Missing + float64(mismatch)*w.Mismatch) /
		float64(total)

	score := math.Max(0, (1.0-penalty)*100.0)
	// Round to two decimal places.
	score = math.Round(score*100) / 100

	return Result{
		Score:    score,
		Total:    total,
		Missing:  missing,
		Mismatch: mismatch,
		Matched:  matched,
	}
}
