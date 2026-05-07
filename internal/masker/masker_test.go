package masker_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/masker"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch, Values: map[string]string{"a.env": "myapp", "b.env": "myapp"}},
		{Key: "DB_PASSWORD", Status: diff.StatusMismatch, Values: map[string]string{"a.env": "supersecret", "b.env": "anothersecret"}},
		{Key: "API_TOKEN", Status: diff.StatusMissing, Values: map[string]string{"a.env": "tok_abc123xyz"}},
		{Key: "PORT", Status: diff.StatusMatch, Values: map[string]string{"a.env": "8080", "b.env": "8080"}},
	}
}

func TestApply_NonSensitiveKeysUnchanged(t *testing.T) {
	opts := masker.DefaultOptions()
	results := sampleResults()
	out := masker.Apply(results, opts)

	for _, r := range out {
		if r.Key == "APP_NAME" || r.Key == "PORT" {
			for file, val := range r.Values {
				original := results[indexOf(results, r.Key)].Values[file]
				if val != original {
					t.Errorf("key %s file %s: expected %q, got %q", r.Key, file, original, val)
				}
			}
		}
	}
}

func TestApply_SensitiveKeysMasked(t *testing.T) {
	opts := masker.DefaultOptions()
	out := masker.Apply(sampleResults(), opts)

	for _, r := range out {
		if r.Key == "DB_PASSWORD" || r.Key == "API_TOKEN" {
			for _, val := range r.Values {
				if !strings.Contains(val, "*") {
					t.Errorf("key %s: expected masked value, got %q", r.Key, val)
				}
			}
		}
	}
}

func TestApply_MaskPreservesPrefix(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.PrefixLen = 3
	opts.SuffixLen = 3
	out := masker.Apply(sampleResults(), opts)

	for _, r := range out {
		if r.Key == "API_TOKEN" {
			val := r.Values["a.env"]
			if !strings.HasPrefix(val, "tok") {
				t.Errorf("expected prefix 'tok', got %q", val)
			}
			if !strings.HasSuffix(val, "xyz") {
				t.Errorf("expected suffix 'xyz', got %q", val)
			}
		}
	}
}

func TestApply_ShortValueFullyMasked(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.PrefixLen = 4
	opts.SuffixLen = 4
	opts.MinMaskLen = 3
	results := []diff.Result{
		{Key: "DB_PASSWORD", Status: diff.StatusMatch, Values: map[string]string{"a.env": "ab"}},
	}
	out := masker.Apply(results, opts)
	val := out[0].Values["a.env"]
	if !strings.ContainsRune(val, '*') {
		t.Errorf("expected fully masked short value, got %q", val)
	}
	if strings.Contains(val, "ab") {
		t.Errorf("short value should not reveal original chars, got %q", val)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	opts := masker.DefaultOptions()
	original := sampleResults()
	masker.Apply(original, opts)
	if original[1].Values["a.env"] != "supersecret" {
		t.Error("original results were mutated")
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := masker.Apply(nil, masker.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d results", len(out))
	}
}

func indexOf(results []diff.Result, key string) int {
	for i, r := range results {
		if r.Key == key {
			return i
		}
	}
	return -1
}
