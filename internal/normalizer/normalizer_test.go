package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/normalizer"
)

func TestApply_TrimSpace(t *testing.T) {
	env := map[string]string{
		"  KEY  ": "  value  ",
	}
	opts := normalizer.DefaultOptions()
	opts.LowercaseKeys = false
	out := normalizer.Apply(env, opts)
	if v, ok := out["KEY"]; !ok || v != "value" {
		t.Errorf("expected KEY=value, got %q=%q", "KEY", v)
	}
}

func TestApply_LowercaseKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := normalizer.DefaultOptions()
	opts.LowercaseKeys = true
	out := normalizer.Apply(env, opts)
	if _, ok := out["db_host"]; !ok {
		t.Error("expected key to be lowercased to db_host")
	}
}

func TestApply_NormalizeBools_True(t *testing.T) {
	cases := []string{"true", "1", "yes", "on", "TRUE", "YES"}
	opts := normalizer.DefaultOptions()
	for _, raw := range cases {
		out := normalizer.Apply(map[string]string{"FLAG": raw}, opts)
		if out["FLAG"] != "true" {
			t.Errorf("expected true for input %q, got %q", raw, out["FLAG"])
		}
	}
}

func TestApply_NormalizeBools_False(t *testing.T) {
	cases := []string{"false", "0", "no", "off", "FALSE", "No"}
	opts := normalizer.DefaultOptions()
	for _, raw := range cases {
		out := normalizer.Apply(map[string]string{"FLAG": raw}, opts)
		if out["FLAG"] != "false" {
			t.Errorf("expected false for input %q, got %q", raw, out["FLAG"])
		}
	}
}

func TestApply_NonBoolValueUnchanged(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	opts := normalizer.DefaultOptions()
	out := normalizer.Apply(env, opts)
	if out["NAME"] != "alice" {
		t.Errorf("expected alice, got %q", out["NAME"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"FLAG": "yes"}
	opts := normalizer.DefaultOptions()
	normalizer.Apply(env, opts)
	if env["FLAG"] != "yes" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_CollapseEmpty(t *testing.T) {
	env := map[string]string{"EMPTY": "", "SET": "val"}
	opts := normalizer.DefaultOptions()
	opts.CollapseEmpty = true
	out := normalizer.Apply(env, opts)
	if out["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", out["EMPTY"])
	}
	if out["SET"] != "val" {
		t.Errorf("expected val, got %q", out["SET"])
	}
}

func TestDefaultOptions_Defaults(t *testing.T) {
	opts := normalizer.DefaultOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true by default")
	}
	if opts.LowercaseKeys {
		t.Error("expected LowercaseKeys to be false by default")
	}
	if !opts.NormalizeBools {
		t.Error("expected NormalizeBools to be true by default")
	}
}
