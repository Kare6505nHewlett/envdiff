package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/profiler"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	out := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Result{Key: fmt.Sprintf("KEY_%d", i), Status: s}
	}
	return out
}

func TestCompute_EmptyInput(t *testing.T) {
	p := profiler.Compute(nil)
	if p.Score != 100 {
		t.Errorf("expected score 100, got %.1f", p.Score)
	}
	if p.Grade != profiler.GradeA {
		t.Errorf("expected grade A, got %s", p.Grade)
	}
}

func TestCompute_AllMatch(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMatch},
	}
	p := profiler.Compute(results)
	if p.Score != 100 {
		t.Errorf("expected 100, got %.1f", p.Score)
	}
	if p.Grade != profiler.GradeA {
		t.Errorf("expected A, got %s", p.Grade)
	}
	if p.MatchCount != 2 {
		t.Errorf("expected 2 matches, got %d", p.MatchCount)
	}
}

func TestCompute_AllMissing(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMissing},
		{Key: "B", Status: diff.StatusMissing},
	}
	p := profiler.Compute(results)
	if p.Score != 0 {
		t.Errorf("expected 0, got %.1f", p.Score)
	}
	if p.Grade != profiler.GradeF {
		t.Errorf("expected F, got %s", p.Grade)
	}
}

func TestCompute_MixedResults(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMatch},
		{Key: "C", Status: diff.StatusMatch},
		{Key: "D", Status: diff.StatusMismatch},
	}
	p := profiler.Compute(results)
	if p.MismatchCount != 1 {
		t.Errorf("expected 1 mismatch, got %d", p.MismatchCount)
	}
	if p.Score <= 80 {
		t.Errorf("expected score > 80, got %.1f", p.Score)
	}
}

func TestCompute_SummaryContainsGrade(t *testing.T) {
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
	}
	p := profiler.Compute(results)
	if !strings.Contains(p.Summary, "grade=A") {
		t.Errorf("summary missing grade: %s", p.Summary)
	}
}
