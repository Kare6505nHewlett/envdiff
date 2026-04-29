package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignorer"
)

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envdiffignore")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write ignore file: %v", err)
	}
	return path
}

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissing},
		{Key: "DB_PASS", Status: diff.StatusMismatch},
		{Key: "SECRET_KEY", Status: diff.StatusMissing},
		{Key: "SECRET_TOKEN", Status: diff.StatusMismatch},
		{Key: "APP_ENV", Status: diff.StatusMatch},
	}
}

func TestLoadFile_MissingFileReturnsEmpty(t *testing.T) {
	ig, err := ignorer.LoadFile("/nonexistent/.envdiffignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	results := ig.Apply(sampleResults())
	if len(results) != len(sampleResults()) {
		t.Errorf("expected %d results, got %d", len(sampleResults()), len(results))
	}
}

func TestApply_ExactKey(t *testing.T) {
	path := writeIgnoreFile(t, "DB_HOST\nAPP_ENV\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	results := ig.Apply(sampleResults())
	for _, r := range results {
		if r.Key == "DB_HOST" || r.Key == "APP_ENV" {
			t.Errorf("key %q should have been ignored", r.Key)
		}
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results after filtering, got %d", len(results))
	}
}

func TestApply_PrefixPattern(t *testing.T) {
	path := writeIgnoreFile(t, "SECRET_*\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	results := ig.Apply(sampleResults())
	for _, r := range results {
		if r.Key == "SECRET_KEY" || r.Key == "SECRET_TOKEN" {
			t.Errorf("key %q should have been ignored by prefix", r.Key)
		}
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestApply_CommentsAndBlankLines(t *testing.T) {
	path := writeIgnoreFile(t, "# this is a comment\n\nDB_PASS\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	results := ig.Apply(sampleResults())
	if len(results) != 4 {
		t.Errorf("expected 4 results, got %d", len(results))
	}
}

func TestApply_NoIgnoreRules(t *testing.T) {
	path := writeIgnoreFile(t, "# only comments\n")
	ig, err := ignorer.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	results := ig.Apply(sampleResults())
	if len(results) != len(sampleResults()) {
		t.Errorf("expected all %d results, got %d", len(sampleResults()), len(results))
	}
}
