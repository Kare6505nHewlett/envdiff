package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "HOST", File: "staging.env", Status: diff.StatusMatch},
		{Key: "PORT", File: "staging.env", Status: diff.StatusMismatch},
		{Key: "SECRET", File: "staging.env", Status: diff.StatusMissing},
		{Key: "DEBUG", File: "prod.env", Status: diff.StatusMatch},
		{Key: "DB_URL", File: "prod.env", Status: diff.StatusMissing},
	}
}

func TestCompute_Counts(t *testing.T) {
	stats := summary.Compute(sampleResults())

	if stats.Total != 5 {
		t.Errorf("expected Total=5, got %d", stats.Total)
	}
	if stats.Matched != 2 {
		t.Errorf("expected Matched=2, got %d", stats.Matched)
	}
	if stats.Missing != 2 {
		t.Errorf("expected Missing=2, got %d", stats.Missing)
	}
	if stats.Mismatch != 1 {
		t.Errorf("expected Mismatch=1, got %d", stats.Mismatch)
	}
}

func TestCompute_FileSet(t *testing.T) {
	stats := summary.Compute(sampleResults())
	if len(stats.Files) != 2 {
		t.Errorf("expected 2 unique files, got %d", len(stats.Files))
	}
}

func TestCompute_EmptyResults(t *testing.T) {
	stats := summary.Compute([]diff.Result{})
	if stats.Total != 0 || stats.Matched != 0 || stats.Missing != 0 || stats.Mismatch != 0 {
		t.Error("expected all zero stats for empty input")
	}
}

func TestPrint_DifferencesFound(t *testing.T) {
	stats := summary.Compute(sampleResults())
	var buf bytes.Buffer
	summary.Print(&buf, stats)
	out := buf.String()

	if !strings.Contains(out, "DIFFERENCES FOUND") {
		t.Errorf("expected 'DIFFERENCES FOUND' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Missing        : 2") {
		t.Errorf("expected missing count in output, got:\n%s", out)
	}
}

func TestPrint_AllMatched(t *testing.T) {
	results := []diff.Result{
		{Key: "HOST", File: "a.env", Status: diff.StatusMatch},
		{Key: "PORT", File: "a.env", Status: diff.StatusMatch},
	}
	stats := summary.Compute(results)
	var buf bytes.Buffer
	summary.Print(&buf, stats)
	out := buf.String()

	if !strings.Contains(out, "ALL MATCHED") {
		t.Errorf("expected 'ALL MATCHED' in output, got:\n%s", out)
	}
}
