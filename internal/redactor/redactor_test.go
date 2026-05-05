package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/redactor"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Value: "myapp", File: "prod.env", Status: diff.StatusMatch},
		{Key: "DB_PASSWORD", Value: "s3cr3t", File: "prod.env", Status: diff.StatusMatch},
		{Key: "API_TOKEN", Value: "tok_abc123", File: "prod.env", Status: diff.StatusMismatch,
			OtherValues: map[string]string{"staging.env": "tok_xyz999"}},
		{Key: "SECRET_KEY", Value: "hunter2", File: "prod.env", Status: diff.StatusMissing},
		{Key: "PORT", Value: "8080", File: "prod.env", Status: diff.StatusMatch},
	}
}

func TestApply_NonSensitiveKeysUnchanged(t *testing.T) {
	results := sampleResults()
	out := redactor.Apply(results, redactor.Options{})

	for _, r := range out {
		if r.Key == "APP_NAME" && r.Value != "myapp" {
			t.Errorf("APP_NAME should not be redacted, got %q", r.Value)
		}
		if r.Key == "PORT" && r.Value != "8080" {
			t.Errorf("PORT should not be redacted, got %q", r.Value)
		}
	}
}

func TestApply_SensitiveKeysRedacted(t *testing.T) {
	out := redactor.Apply(sampleResults(), redactor.Options{})

	sensitiveKeys := []string{"DB_PASSWORD", "API_TOKEN", "SECRET_KEY"}
	for _, r := range out {
		for _, sk := range sensitiveKeys {
			if r.Key == sk && r.Value != redactor.DefaultPlaceholder {
				t.Errorf("key %q value should be redacted, got %q", r.Key, r.Value)
			}
		}
	}
}

func TestApply_OtherValuesRedacted(t *testing.T) {
	out := redactor.Apply(sampleResults(), redactor.Options{})

	for _, r := range out {
		if r.Key == "API_TOKEN" {
			for file, val := range r.OtherValues {
				if val != redactor.DefaultPlaceholder {
					t.Errorf("OtherValues[%q] should be redacted, got %q", file, val)
				}
			}
		}
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	opts := redactor.Options{Placeholder: "<hidden>"}
	out := redactor.Apply(sampleResults(), opts)

	for _, r := range out {
		if r.Key == "DB_PASSWORD" && r.Value != "<hidden>" {
			t.Errorf("expected <hidden>, got %q", r.Value)
		}
	}
}

func TestApply_ExtraPatterns(t *testing.T) {
	opts := redactor.Options{ExtraPatterns: []string{"INTERNAL"}}
	input := []diff.Result{
		{Key: "INTERNAL_URL", Value: "http://internal", File: "a.env"},
		{Key: "PUBLIC_URL", Value: "http://public", File: "a.env"},
	}
	out := redactor.Apply(input, opts)

	if out[0].Value != redactor.DefaultPlaceholder {
		t.Errorf("INTERNAL_URL should be redacted, got %q", out[0].Value)
	}
	if out[1].Value != "http://public" {
		t.Errorf("PUBLIC_URL should not be redacted, got %q", out[1].Value)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	original := sampleResults()
	redactor.Apply(original, redactor.Options{})

	for _, r := range original {
		if r.Key == "DB_PASSWORD" && r.Value != "s3cr3t" {
			t.Error("original results should not be mutated")
		}
	}
}
