package scoper_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scoper"
)

func buildScopedResults() []scoper.Result {
	return []scoper.Result{
		{
			Scope: "production",
			Entries: []diff.Result{
				{Key: "DB_HOST", File: "production.env", Status: diff.StatusMatch},
				{Key: "API_KEY", File: "production.env", Status: diff.StatusMissing},
			},
		},
		{
			Scope: "staging",
			Entries: []diff.Result{
				{Key: "DB_HOST", File: "staging.env", Status: diff.StatusMismatch},
			},
		},
	}
}

func TestReportText_ContainsScopeHeader(t *testing.T) {
	var buf bytes.Buffer
	scoper.ReportText(&buf, buildScopedResults())
	out := buf.String()
	if !strings.Contains(out, "PRODUCTION") {
		t.Error("expected PRODUCTION header in text output")
	}
	if !strings.Contains(out, "STAGING") {
		t.Error("expected STAGING header in text output")
	}
}

func TestReportText_NoResults(t *testing.T) {
	var buf bytes.Buffer
	scoper.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No scoped results") {
		t.Error("expected empty message")
	}
}

func TestReportJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	if err := scoper.ReportJSON(&buf, buildScopedResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(out))
	}
}

func TestReportJSON_CountMatchesEntries(t *testing.T) {
	var buf bytes.Buffer
	_ = scoper.ReportJSON(&buf, buildScopedResults())
	var out []map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)
	for _, s := range out {
		count := int(s["count"].(float64))
		entries := s["entries"].([]interface{})
		if count != len(entries) {
			t.Errorf("scope %q: count %d != entries %d", s["scope"], count, len(entries))
		}
	}
}
