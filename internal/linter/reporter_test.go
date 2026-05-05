package linter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func sampleFindings() []linter.Finding {
	return []linter.Finding{
		{Key: "app_env", File: "dev.env", Rule: "lowercase-key", Message: "key has lowercase", Severity: linter.SeverityWarn},
		{Key: "APP ENV", File: "dev.env", Rule: "space-in-key", Message: "key has space", Severity: linter.SeverityError},
	}
}

func TestReportText_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	linter.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "no issues") {
		t.Errorf("expected 'no issues' message, got: %q", buf.String())
	}
}

func TestReportText_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	linter.ReportText(&buf, sampleFindings())
	out := buf.String()
	if !strings.Contains(out, "[warn]") {
		t.Error("expected [warn] in output")
	}
	if !strings.Contains(out, "[error]") {
		t.Error("expected [error] in output")
	}
	if !strings.Contains(out, "lowercase-key") {
		t.Error("expected rule name in output")
	}
}

func TestReportJSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	err := linter.ReportJSON(&buf, sampleFindings())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []linter.Finding
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 findings, got %d", len(out))
	}
}

func TestReportJSON_EmptyFindings(t *testing.T) {
	var buf bytes.Buffer
	err := linter.ReportJSON(&buf, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Error("expected JSON array even for empty findings")
	}
}

func TestCountBySeverity(t *testing.T) {
	counts := linter.CountBySeverity(sampleFindings())
	if counts[linter.SeverityWarn] != 1 {
		t.Errorf("expected 1 warn, got %d", counts[linter.SeverityWarn])
	}
	if counts[linter.SeverityError] != 1 {
		t.Errorf("expected 1 error, got %d", counts[linter.SeverityError])
	}
}
