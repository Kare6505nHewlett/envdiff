package exporter_test

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMismatch, File: "prod.env", BaseValue: "localhost", CompareValue: "db.prod"},
		{Key: "API_KEY", Status: diff.StatusMissing, File: "prod.env", BaseValue: "abc123", CompareValue: ""},
	}
}

func TestExportCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, sampleResults(), exporter.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 rows (header+2), got %d", len(records))
	}
	if records[0][0] != "key" {
		t.Errorf("expected header 'key', got %q", records[0][0])
	}
	if records[1][0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in row 1, got %q", records[1][0])
	}
}

func TestExportMarkdown(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, sampleResults(), exporter.FormatMarkdown); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key |") {
		t.Error("expected markdown header row")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
}

func TestExportJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, sampleResults(), exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Error("expected DB_HOST in JSON output")
	}
}

func TestExportUnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleResults(), exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExportCSV_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, []diff.Result{}, exporter.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header row, got %d lines", len(lines))
	}
}
