package annotator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/diff"
)

func writeDescFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.descriptions")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeDescFile: %v", err)
	}
	return p
}

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: "match"},
		{Key: "API_KEY", Status: "missing"},
		{Key: "LOG_LEVEL", Status: "mismatch"},
	}
}

func TestLoadDescriptions_ValidFile(t *testing.T) {
	p := writeDescFile(t, "DB_HOST=Database hostname\n# comment\nAPI_KEY=Secret API key\n")
	desc, err := annotator.LoadDescriptions(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if desc["DB_HOST"] != "Database hostname" {
		t.Errorf("expected 'Database hostname', got %q", desc["DB_HOST"])
	}
	if desc["API_KEY"] != "Secret API key" {
		t.Errorf("expected 'Secret API key', got %q", desc["API_KEY"])
	}
}

func TestLoadDescriptions_MissingFileReturnsEmpty(t *testing.T) {
	desc, err := annotator.LoadDescriptions("/nonexistent/.env.descriptions")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(desc) != 0 {
		t.Errorf("expected empty map, got %v", desc)
	}
}

func TestLoadDescriptions_SkipsBlankAndCommentLines(t *testing.T) {
	p := writeDescFile(t, "\n# skip me\n  \nFOO=bar\n")
	desc, err := annotator.LoadDescriptions(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(desc) != 1 || desc["FOO"] != "bar" {
		t.Errorf("unexpected desc map: %v", desc)
	}
}

func TestApply_AttachesDescriptions(t *testing.T) {
	desc := map[string]string{
		"DB_HOST":  "Database hostname",
		"LOG_LEVEL": "Logging verbosity",
	}
	annotated := annotator.Apply(sampleResults(), desc)
	if len(annotated) != 3 {
		t.Fatalf("expected 3 results, got %d", len(annotated))
	}
	if annotated[0].Description != "Database hostname" {
		t.Errorf("DB_HOST description mismatch: %q", annotated[0].Description)
	}
	if annotated[1].Description != "" {
		t.Errorf("API_KEY should have no description, got %q", annotated[1].Description)
	}
	if annotated[2].Description != "Logging verbosity" {
		t.Errorf("LOG_LEVEL description mismatch: %q", annotated[2].Description)
	}
}

func TestApply_NilDescriptionsIsSafe(t *testing.T) {
	annotated := annotator.Apply(sampleResults(), nil)
	for _, a := range annotated {
		if a.Description != "" {
			t.Errorf("expected empty description, got %q for key %s", a.Description, a.Key)
		}
	}
}
