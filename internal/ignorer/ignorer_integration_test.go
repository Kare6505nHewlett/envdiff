package ignorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignorer"
)

// TestApply_CombinedRules verifies that exact keys and prefix patterns
// can coexist in the same ignore file.
func TestApply_CombinedRules(t *testing.T) {
	path := writeIgnoreFile(t, "# ignore db password and all secret keys\nDB_PASS\nSECRET_*\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	input := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissing},
		{Key: "DB_PASS", Status: diff.StatusMismatch},
		{Key: "SECRET_KEY", Status: diff.StatusMissing},
		{Key: "SECRET_TOKEN", Status: diff.StatusMismatch},
		{Key: "APP_ENV", Status: diff.StatusMatch},
		{Key: "PORT", Status: diff.StatusMatch},
	}

	results := ig.Apply(input)

	expectedKeys := map[string]bool{
		"DB_HOST": true,
		"APP_ENV": true,
		"PORT":    true,
	}

	if len(results) != len(expectedKeys) {
		t.Fatalf("expected %d results, got %d", len(expectedKeys), len(results))
	}

	for _, r := range results {
		if !expectedKeys[r.Key] {
			t.Errorf("unexpected key in results: %q", r.Key)
		}
	}
}

// TestApply_NilSafeOnEmptyInput verifies Apply handles an empty slice gracefully.
func TestApply_NilSafeOnEmptyInput(t *testing.T) {
	path := writeIgnoreFile(t, "DB_HOST\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	results := ig.Apply([]diff.Result{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
