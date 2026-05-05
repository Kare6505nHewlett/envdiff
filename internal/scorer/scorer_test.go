package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

func makeResults(matched, missing, mismatch int) []diff.Result {
	var out []diff.Result
	for i := 0; i < matched; i++ {
		out = append(out, diff.Result{Key: "K", Status: diff.StatusMatch})
	}
	for i := 0; i < missing; i++ {
		out = append(out, diff.Result{Key: "K", Status: diff.StatusMissing})
	}
	for i := 0; i < mismatch; i++ {
		out = append(out, diff.Result{Key: "K", Status: diff.StatusMismatch})
	}
	return out
}

func TestScore_PerfectScore(t *testing.T) {
	res := scorer.Score(makeResults(10, 0, 0), nil)
	if res.Score != 100.0 {
		t.Errorf("expected 100.0, got %v", res.Score)
	}
}

func TestScore_EmptyInput(t *testing.T) {
	res := scorer.Score(nil, nil)
	if res.Score != 100.0 {
		t.Errorf("expected 100.0 for empty input, got %v", res.Score)
	}
	if res.Total != 0 {
		t.Errorf("expected total 0, got %d", res.Total)
	}
}

func TestScore_AllMissing(t *testing.T) {
	res := scorer.Score(makeResults(0, 10, 0), nil)
	if res.Score != 0.0 {
		t.Errorf("expected 0.0 for all missing, got %v", res.Score)
	}
}

func TestScore_MixedResults(t *testing.T) {
	// 8 matched, 1 missing, 1 mismatch with default weights
	// penalty = (1*1.0 + 1*0.5) / 10 = 0.15 → score = 85.0
	res := scorer.Score(makeResults(8, 1, 1), nil)
	if res.Score != 85.0 {
		t.Errorf("expected 85.0, got %v", res.Score)
	}
	if res.Total != 10 {
		t.Errorf("expected total 10, got %d", res.Total)
	}
}

func TestScore_CustomWeights(t *testing.T) {
	w := &scorer.Weights{Missing: 0.5, Mismatch: 0.5}
	// 8 matched, 2 missing → penalty = 1.0/10 = 0.1 → score = 90.0
	res := scorer.Score(makeResults(8, 2, 0), w)
	if res.Score != 90.0 {
		t.Errorf("expected 90.0, got %v", res.Score)
	}
}

func TestScore_CountsAreCorrect(t *testing.T) {
	res := scorer.Score(makeResults(5, 3, 2), nil)
	if res.Matched != 5 || res.Missing != 3 || res.Mismatch != 2 {
		t.Errorf("unexpected counts: %+v", res)
	}
}
