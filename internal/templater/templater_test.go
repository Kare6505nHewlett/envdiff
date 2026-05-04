package templater_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/templater"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", File: "production.env", Status: diff.StatusMissing},
		{Key: "API_KEY", File: "staging.env", Status: diff.StatusMismatch},
		{Key: "PORT", File: "production.env", Status: diff.StatusMatch},
	}
}

func TestGenerate_KeysAreSorted(t *testing.T) {
	var buf strings.Builder
	err := templater.Generate(&buf, sampleResults(), templater.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmpty(strings.Split(buf.String(), "\n"))
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "API_KEY=") {
		t.Errorf("expected first key API_KEY, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "DB_HOST=") {
		t.Errorf("expected second key DB_HOST, got %q", lines[1])
	}
}

func TestGenerate_PlaceholderApplied(t *testing.T) {
	var buf strings.Builder
	opts := templater.Options{Placeholder: "CHANGEME"}
	_ = templater.Generate(&buf, sampleResults(), opts)
	if !strings.Contains(buf.String(), "API_KEY=CHANGEME") {
		t.Errorf("placeholder not applied: %s", buf.String())
	}
}

func TestGenerate_IncludeComments(t *testing.T) {
	var buf strings.Builder
	opts := templater.Options{IncludeComments: true}
	_ = templater.Generate(&buf, sampleResults(), opts)
	if !strings.Contains(buf.String(), "# status:") {
		t.Errorf("expected comments in output: %s", buf.String())
	}
}

func TestGenerate_DeduplicatesKeys(t *testing.T) {
	results := []diff.Result{
		{Key: "FOO", File: "a.env", Status: diff.StatusMatch},
		{Key: "FOO", File: "b.env", Status: diff.StatusMismatch},
	}
	var buf strings.Builder
	_ = templater.Generate(&buf, results, templater.Options{})
	count := strings.Count(buf.String(), "FOO=")
	if count != 1 {
		t.Errorf("expected 1 FOO entry, got %d", count)
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env.template")
	err := templater.WriteFile(out, sampleResults(), templater.Options{})
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if len(data) == 0 {
		t.Error("expected non-empty template file")
	}
}

func nonEmpty(lines []string) []string {
	var out []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
