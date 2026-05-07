package truncator_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/truncator"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{
			Key:    "SHORT",
			Status: diff.StatusMatch,
			Values: map[string]string{"a.env": "hello", "b.env": "hello"},
		},
		{
			Key:    "LONG_VAL",
			Status: diff.StatusMismatch,
			Values: map[string]string{
				"a.env": strings.Repeat("x", 80),
				"b.env": strings.Repeat("y", 80),
			},
		},
		{
			Key:    "MULTILINE",
			Status: diff.StatusMatch,
			Values: map[string]string{"a.env": "line1\nline2"},
		},
	}
}

func TestApply_ShortValuesUnchanged(t *testing.T) {
	results := sampleResults()
	out := truncator.Apply(results, truncator.DefaultOptions())
	if out[0].Values["a.env"] != "hello" {
		t.Errorf("expected 'hello', got %q", out[0].Values["a.env"])
	}
}

func TestApply_LongValuesTruncated(t *testing.T) {
	opts := truncator.DefaultOptions()
	out := truncator.Apply(sampleResults(), opts)
	val := out[1].Values["a.env"]
	if len(val) > opts.MaxLen {
		t.Errorf("expected value <= %d chars, got %d", opts.MaxLen, len(val))
	}
	if !strings.HasSuffix(val, "...") {
		t.Errorf("expected ellipsis suffix, got %q", val)
	}
}

func TestApply_NewlinesEscaped(t *testing.T) {
	out := truncator.Apply(sampleResults(), truncator.DefaultOptions())
	val := out[2].Values["a.env"]
	if strings.Contains(val, "\n") {
		t.Errorf("expected newlines to be escaped, got %q", val)
	}
	if !strings.Contains(val, `\n`) {
		t.Errorf("expected escaped newline marker in %q", val)
	}
}

func TestApply_OriginalNotMutated(t *testing.T) {
	original := sampleResults()
	long := strings.Repeat("x", 80)
	truncator.Apply(original, truncator.DefaultOptions())
	if original[1].Values["a.env"] != long {
		t.Error("original result was mutated")
	}
}

func TestApply_CustomEllipsisAndMaxLen(t *testing.T) {
	opts := truncator.Options{MaxLen: 10, Ellipsis: "~~"}
	out := truncator.Apply(sampleResults(), opts)
	val := out[1].Values["a.env"]
	if len(val) > 10 {
		t.Errorf("expected max 10 chars, got %d", len(val))
	}
	if !strings.HasSuffix(val, "~~") {
		t.Errorf("expected custom ellipsis, got %q", val)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := truncator.Apply(nil, truncator.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d results", len(out))
	}
}
