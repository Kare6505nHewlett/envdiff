package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", File: "production.env", Status: "missing"},
		{Key: "DB_PASS", File: "production.env", Status: "mismatch"},
		{Key: "APP_NAME", File: "staging.env", Status: "match"},
		{Key: "APP_PORT", File: "staging.env", Status: "missing"},
		{Key: "SECRET", File: "production.env", Status: "mismatch"},
	}
}

func TestApply_ByFile(t *testing.T) {
	groups := grouper.Apply(sampleResults(), grouper.ByFile)

	if len(groups) != 2 {
		t.Fatalf("expected 2 file groups, got %d", len(groups))
	}
	if groups[0].Label != "production.env" {
		t.Errorf("expected first label 'production.env', got %q", groups[0].Label)
	}
	if len(groups[0].Results) != 3 {
		t.Errorf("expected 3 results for production.env, got %d", len(groups[0].Results))
	}
}

func TestApply_ByPrefix(t *testing.T) {
	groups := grouper.Apply(sampleResults(), grouper.ByPrefix)

	labels := map[string]bool{}
	for _, g := range groups {
		labels[g.Label] = true
	}

	for _, want := range []string{"APP", "DB", "(no prefix)"} {
		if !labels[want] {
			t.Errorf("expected group label %q to exist", want)
		}
	}
}

func TestApply_ByStatus(t *testing.T) {
	groups := grouper.Apply(sampleResults(), grouper.ByStatus)

	statusCount := map[string]int{}
	for _, g := range groups {
		statusCount[g.Label] = len(g.Results)
	}

	if statusCount["missing"] != 2 {
		t.Errorf("expected 2 missing, got %d", statusCount["missing"])
	}
	if statusCount["mismatch"] != 2 {
		t.Errorf("expected 2 mismatch, got %d", statusCount["mismatch"])
	}
	if statusCount["match"] != 1 {
		t.Errorf("expected 1 match, got %d", statusCount["match"])
	}
}

func TestApply_EmptyInput(t *testing.T) {
	groups := grouper.Apply([]diff.Result{}, grouper.ByFile)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestApply_SortedLabels(t *testing.T) {
	groups := grouper.Apply(sampleResults(), grouper.ByFile)
	for i := 1; i < len(groups); i++ {
		if groups[i-1].Label > groups[i].Label {
			t.Errorf("groups not sorted: %q > %q", groups[i-1].Label, groups[i].Label)
		}
	}
}
