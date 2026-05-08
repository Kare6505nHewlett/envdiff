package deduplicator_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduplicator"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", File: "staging.env", Status: "match", Value: "localhost"},
		{Key: "DB_HOST", File: "prod.env", Status: "mismatch", Value: "db.prod.example.com"},
		{Key: "API_KEY", File: "staging.env", Status: "missing", Value: ""},
		{Key: "API_KEY", File: "prod.env", Status: "missing", Value: ""},
		{Key: "TIMEOUT", File: "staging.env", Status: "extra", Value: "30"},
	}
}

func TestApply_DefaultKeepsHighestPriority(t *testing.T) {
	results := sampleResults()
	out := deduplicator.Apply(results, deduplicator.DefaultOptions())

	index := make(map[string]diff.Result)
	for _, r := range out {
		index[r.Key] = r
	}

	if len(out) != 3 {
		t.Fatalf("expected 3 unique keys, got %d", len(out))
	}
	if index["DB_HOST"].Status != "mismatch" {
		t.Errorf("DB_HOST: expected mismatch, got %s", index["DB_HOST"].Status)
	}
	if index["API_KEY"].Status != "missing" {
		t.Errorf("API_KEY: expected missing, got %s", index["API_KEY"].Status)
	}
}

func TestApply_SortedByKey(t *testing.T) {
	out := deduplicator.Apply(sampleResults(), deduplicator.DefaultOptions())
	for i := 1; i < len(out); i++ {
		if out[i].Key < out[i-1].Key {
			t.Errorf("results not sorted: %s before %s", out[i-1].Key, out[i].Key)
		}
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := deduplicator.Apply([]diff.Result{}, deduplicator.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d results", len(out))
	}
}

func TestApply_KeepAll_RemovesExactDuplicates(t *testing.T) {
	input := []diff.Result{
		{Key: "FOO", File: "a.env", Status: "match", Value: "bar"},
		{Key: "FOO", File: "a.env", Status: "match", Value: "bar"},
		{Key: "FOO", File: "b.env", Status: "mismatch", Value: "baz"},
	}
	out := deduplicator.Apply(input, deduplicator.Options{KeepAll: true})
	if len(out) != 2 {
		t.Errorf("expected 2 results (exact dup removed), got %d", len(out))
	}
}

func TestApply_KeepAll_PreservesDistinctEntries(t *testing.T) {
	input := []diff.Result{
		{Key: "X", File: "a.env", Status: "missing", Value: ""},
		{Key: "X", File: "b.env", Status: "missing", Value: ""},
	}
	out := deduplicator.Apply(input, deduplicator.Options{KeepAll: true})
	if len(out) != 2 {
		t.Errorf("expected 2 distinct entries, got %d", len(out))
	}
}

func TestDefaultOptions_KeepAllIsFalse(t *testing.T) {
	opts := deduplicator.DefaultOptions()
	if opts.KeepAll {
		t.Error("expected KeepAll to be false by default")
	}
}
