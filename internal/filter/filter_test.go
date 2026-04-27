package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "APP_SECRET", Status: diff.StatusMismatch},
		{Key: "DB_HOST", Status: diff.StatusMissing},
		{Key: "DB_PORT", Status: diff.StatusExtra},
		{Key: "LOG_LEVEL", Status: diff.StatusMatch},
	}
}

func TestApply_NoFilter(t *testing.T) {
	results := filter.Apply(sampleResults(), filter.Options{})
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	results := filter.Apply(sampleResults(), filter.Options{OnlyMissing: true})
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Fatalf("expected only DB_HOST, got %+v", results)
	}
}

func TestApply_OnlyMismatch(t *testing.T) {
	results := filter.Apply(sampleResults(), filter.Options{OnlyMismatch: true})
	if len(results) != 1 || results[0].Key != "APP_SECRET" {
		t.Fatalf("expected only APP_SECRET, got %+v", results)
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	results := filter.Apply(sampleResults(), filter.Options{Prefix: "DB_"})
	if len(results) != 2 {
		t.Fatalf("expected 2 DB_ results, got %d", len(results))
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	results := filter.Apply(sampleResults(), filter.Options{ExcludeKeys: []string{"LOG_LEVEL", "APP_NAME"}})
	if len(results) != 3 {
		t.Fatalf("expected 3 results after exclusion, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "LOG_LEVEL" || r.Key == "APP_NAME" {
			t.Fatalf("excluded key %q still present", r.Key)
		}
	}
}

func TestApply_PrefixAndOnlyMissing(t *testing.T) {
	results := filter.Apply(sampleResults(), filter.Options{Prefix: "DB_", OnlyMissing: true})
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Fatalf("expected only DB_HOST, got %+v", results)
	}
}
