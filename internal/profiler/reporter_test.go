package profiler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/profiler"
)

func sampleProfile() profiler.Profile {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch},
		{Key: "DB_PASS", Status: diff.StatusMissing},
		{Key: "API_KEY", Status: diff.StatusMismatch},
	}
	return profiler.Compute(results)
}

func TestReportText_ContainsGrade(t *testing.T) {
	var buf bytes.Buffer
	p := sampleProfile()
	profiler.ReportText(&buf, p)
	out := buf.String()
	if !strings.Contains(out, "Grade") {
		t.Errorf("expected 'Grade' in output, got:\n%s", out)
	}
	if !strings.Contains(out, string(p.Grade)) {
		t.Errorf("expected grade %s in output", p.Grade)
	}
}

func TestReportText_ContainsCounts(t *testing.T) {
	var buf bytes.Buffer
	p := sampleProfile()
	profiler.ReportText(&buf, p)
	out := buf.String()
	for _, want := range []string{"Missing", "Mismatch", "Match", "Total"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in text output", want)
		}
	}
}

func TestReportJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	p := sampleProfile()
	if err := profiler.ReportJSON(&buf, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, key := range []string{"grade", "score", "total_keys", "missing", "mismatch", "match"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing JSON field %q", key)
		}
	}
}

func TestReportJSON_ScoreRange(t *testing.T) {
	var buf bytes.Buffer
	p := sampleProfile()
	_ = profiler.ReportJSON(&buf, p)
	var m map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &m)
	score, _ := m["score"].(float64)
	if score < 0 || score > 100 {
		t.Errorf("score out of range: %f", score)
	}
}

// ensure fmt and strings are used
var _ = fmt.Sprintf
var _ = strings.Contains
