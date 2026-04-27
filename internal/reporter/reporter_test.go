package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reporter"
)

func TestReportText_NoResults(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Report(&buf, []diff.Result{}, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "All keys match") {
		t.Errorf("expected match message, got: %s", buf.String())
	}
}

func TestReportText_WithResults(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissing, Files: []string{".env.production"}},
		{Key: "DEBUG", Status: diff.StatusExtra, Files: []string{".env.local"}},
		{Key: "API_KEY", Status: diff.StatusMismatch, Files: []string{".env.local", ".env.production"}},
	}

	var buf bytes.Buffer
	if err := reporter.Report(&buf, results, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[MISSING]") {
		t.Error("expected MISSING label")
	}
	if !strings.Contains(out, "[EXTRA]") {
		t.Error("expected EXTRA label")
	}
	if !strings.Contains(out, "[MISMATCH]") {
		t.Error("expected MISMATCH label")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key DB_HOST in output")
	}
}

func TestReportJSON_ValidOutput(t *testing.T) {
	results := []diff.Result{
		{Key: "SECRET", Status: diff.StatusMissing, Files: []string{".env.staging"}},
	}

	var buf bytes.Buffer
	if err := reporter.Report(&buf, results, reporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded []diff.Result
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded) != 1 || decoded[0].Key != "SECRET" {
		t.Errorf("unexpected decoded results: %+v", decoded)
	}
}

func TestReportUnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.Report(&buf, []diff.Result{}, reporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
