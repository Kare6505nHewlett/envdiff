package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/merger"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", File: "prod.env", Status: "match", Value: "production"},
		{Key: "APP_ENV", File: "staging.env", Status: "match", Value: "production"},
		{Key: "DB_HOST", File: "prod.env", Status: "mismatch", Value: "db.prod.internal"},
		{Key: "DB_HOST", File: "staging.env", Status: "mismatch", Value: "db.staging.internal"},
		{Key: "SECRET_KEY", File: "prod.env", Status: "missing", Value: ""},
		{Key: "SECRET_KEY", File: "staging.env", Status: "match", Value: "s3cr3t"},
	}
}

func TestMerge_ReturnsOneEntryPerKey(t *testing.T) {
	results := sampleResults()
	merged := merger.Merge(results)
	if len(merged) != 3 {
		t.Fatalf("expected 3 merged keys, got %d", len(merged))
	}
}

func TestMerge_SortedByKey(t *testing.T) {
	merged := merger.Merge(sampleResults())
	keys := []string{merged[0].Key, merged[1].Key, merged[2].Key}
	expected := []string{"APP_ENV", "DB_HOST", "SECRET_KEY"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestMerge_StatusMatch(t *testing.T) {
	merged := merger.Merge(sampleResults())
	for _, mk := range merged {
		if mk.Key == "APP_ENV" && mk.Status != "match" {
			t.Errorf("APP_ENV: expected status 'match', got %q", mk.Status)
		}
	}
}

func TestMerge_StatusMismatch(t *testing.T) {
	merged := merger.Merge(sampleResults())
	for _, mk := range merged {
		if mk.Key == "DB_HOST" && mk.Status != "mismatch" {
			t.Errorf("DB_HOST: expected status 'mismatch', got %q", mk.Status)
		}
	}
}

func TestMerge_StatusMissing(t *testing.T) {
	merged := merger.Merge(sampleResults())
	for _, mk := range merged {
		if mk.Key == "SECRET_KEY" && mk.Status != "missing" {
			t.Errorf("SECRET_KEY: expected status 'missing', got %q", mk.Status)
		}
	}
}

func TestMerge_ValuesPopulated(t *testing.T) {
	merged := merger.Merge(sampleResults())
	for _, mk := range merged {
		if mk.Key == "DB_HOST" {
			if mk.Values["prod.env"] != "db.prod.internal" {
				t.Errorf("unexpected prod value: %q", mk.Values["prod.env"])
			}
			if mk.Values["staging.env"] != "db.staging.internal" {
				t.Errorf("unexpected staging value: %q", mk.Values["staging.env"])
			}
		}
	}
}

func TestMerge_EmptyInput(t *testing.T) {
	merged := merger.Merge([]diff.Result{})
	if len(merged) != 0 {
		t.Errorf("expected empty result, got %d entries", len(merged))
	}
}
