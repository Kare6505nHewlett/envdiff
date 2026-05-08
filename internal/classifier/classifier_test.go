package classifier_test

import (
	"testing"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_PASSWORD", Value: "s3cr3t", File: "prod.env", Status: diff.StatusMismatch},
		{Key: "API_URL", Value: "https://api.example.com", File: "prod.env", Status: diff.StatusMatch},
		{Key: "FEATURE_DARK_MODE", Value: "true", File: "prod.env", Status: diff.StatusMatch},
		{Key: "MAX_RETRIES", Value: "5", File: "prod.env", Status: diff.StatusMissing},
		{Key: "APP_NAME", Value: "myapp", File: "prod.env", Status: diff.StatusMatch},
	}
}

func TestApply_ReturnsOneResultPerInput(t *testing.T) {
	results := sampleResults()
	got := classifier.Apply(results)
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApply_SecretKeyword(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_PASSWORD", Value: "abc"},
		{Key: "AUTH_TOKEN", Value: "xyz"},
		{Key: "PRIVATE_KEY", Value: "pem"},
	}
	got := classifier.Apply(results)
	for _, r := range got {
		if r.Category != classifier.CategorySecret {
			t.Errorf("key %q: expected secret, got %s", r.Key, r.Category)
		}
	}
}

func TestApply_URLValue(t *testing.T) {
	results := []diff.Result{
		{Key: "API_URL", Value: "https://example.com"},
		{Key: "CALLBACK", Value: "http://localhost:8080"},
	}
	got := classifier.Apply(results)
	for _, r := range got {
		if r.Category != classifier.CategoryURL {
			t.Errorf("key %q: expected url, got %s", r.Key, r.Category)
		}
	}
}

func TestApply_FlagValue(t *testing.T) {
	for _, v := range []string{"true", "false", "yes", "no", "on", "off", "1", "0"} {
		results := []diff.Result{{Key: "SOME_FLAG", Value: v}}
		got := classifier.Apply(results)
		if got[0].Category != classifier.CategoryFlag {
			t.Errorf("value %q: expected flag, got %s", v, got[0].Category)
		}
	}
}

func TestApply_NumericValue(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", Value: "8080"},
		{Key: "TIMEOUT", Value: "3.14"},
		{Key: "OFFSET", Value: "-42"},
	}
	got := classifier.Apply(results)
	for _, r := range got {
		if r.Category != classifier.CategoryNumeric {
			t.Errorf("key %q value %q: expected numeric, got %s", r.Key, r.Value, r.Category)
		}
	}
}

func TestApply_UnknownValue(t *testing.T) {
	results := []diff.Result{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "REGION", Value: "us-east-1"},
	}
	got := classifier.Apply(results)
	for _, r := range got {
		if r.Category != classifier.CategoryUnknown {
			t.Errorf("key %q: expected unknown, got %s", r.Key, r.Category)
		}
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := classifier.Apply(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d results", len(got))
	}
}
