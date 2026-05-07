package promoter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/promoter"
)

func buildPlan() promoter.Plan {
	return promoter.Plan{
		SourceFile: "prod.env",
		TargetFile: "staging.env",
		Steps: []promoter.Step{
			{Key: "SECRET_KEY", Action: promoter.ActionAdd, NewValue: "abc123", File: "staging.env"},
			{Key: "APP_ENV", Action: promoter.ActionUpdate, OldValue: "staging", NewValue: "production", File: "staging.env"},
		},
	}
}

func TestReportText_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	promoter.ReportText(&buf, buildPlan())
	out := buf.String()
	if !strings.Contains(out, "prod.env") || !strings.Contains(out, "staging.env") {
		t.Errorf("expected file names in output, got: %s", out)
	}
}

func TestReportText_ListsSteps(t *testing.T) {
	var buf bytes.Buffer
	promoter.ReportText(&buf, buildPlan())
	out := buf.String()
	if !strings.Contains(out, "SECRET_KEY") || !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected step keys in output, got: %s", out)
	}
}

func TestReportText_NoSteps(t *testing.T) {
	plan := promoter.Plan{SourceFile: "a.env", TargetFile: "b.env"}
	var buf bytes.Buffer
	promoter.ReportText(&buf, plan)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", buf.String())
	}
}

func TestReportJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	if err := promoter.ReportJSON(&buf, buildPlan()); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["source_file"] != "prod.env" {
		t.Errorf("unexpected source_file: %v", out["source_file"])
	}
	if out["step_count"].(float64) != 2 {
		t.Errorf("expected step_count 2, got %v", out["step_count"])
	}
}

func TestReportJSON_StepsHaveAction(t *testing.T) {
	var buf bytes.Buffer
	promoter.ReportJSON(&buf, buildPlan()) //nolint:errcheck
	var out map[string]interface{}
	json.Unmarshal(buf.Bytes(), &out) //nolint:errcheck
	steps := out["steps"].([]interface{})
	for _, raw := range steps {
		s := raw.(map[string]interface{})
		if s["action"] == nil || s["action"] == "" {
			t.Errorf("step missing action: %v", s)
		}
	}
}

func TestStepSummary_Mixed(t *testing.T) {
	summary := promoter.StepSummary(buildPlan())
	if !strings.Contains(summary, "add") || !strings.Contains(summary, "update") {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestStepSummary_Empty(t *testing.T) {
	plan := promoter.Plan{SourceFile: "a.env", TargetFile: "b.env"}
	if promoter.StepSummary(plan) != "nothing to promote" {
		t.Errorf("expected 'nothing to promote'")
	}
}
