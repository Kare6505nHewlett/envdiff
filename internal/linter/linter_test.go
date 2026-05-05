package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestLint_NoFindings(t *testing.T) {
	entries := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"MAX_CONN": "10",
	}
	opts := linter.DefaultOptions()
	findings := linter.Lint("prod.env", entries, opts)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	entries := map[string]string{"app_env": "staging"}
	opts := linter.DefaultOptions()
	findings := linter.Lint("staging.env", entries, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "lowercase-key" {
		t.Errorf("expected rule lowercase-key, got %q", findings[0].Rule)
	}
	if findings[0].Severity != linter.SeverityWarn {
		t.Errorf("expected severity warn, got %q", findings[0].Severity)
	}
}

func TestLint_SpaceInKey(t *testing.T) {
	entries := map[string]string{"APP ENV": "value"}
	opts := linter.DefaultOptions()
	findings := linter.Lint("bad.env", entries, opts)

	var found bool
	for _, f := range findings {
		if f.Rule == "space-in-key" {
			found = true
			if f.Severity != linter.SeverityError {
				t.Errorf("expected severity error for space-in-key, got %q", f.Severity)
			}
		}
	}
	if !found {
		t.Error("expected space-in-key finding")
	}
}

func TestLint_ValueTooLong(t *testing.T) {
	long := make([]byte, 300)
	for i := range long {
		long[i] = 'x'
	}
	entries := map[string]string{"SECRET_KEY": string(long)}
	opts := linter.DefaultOptions()
	findings := linter.Lint("prod.env", entries, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "value-too-long" {
		t.Errorf("expected rule value-too-long, got %q", findings[0].Rule)
	}
}

func TestLint_DisabledRules(t *testing.T) {
	entries := map[string]string{"app_env": "staging"}
	opts := linter.Options{DisallowLowercase: false}
	findings := linter.Lint("staging.env", entries, opts)
	if len(findings) != 0 {
		t.Errorf("expected no findings with rule disabled, got %d", len(findings))
	}
}

func TestLint_FileNamePropagated(t *testing.T) {
	entries := map[string]string{"bad key": "v"}
	opts := linter.DefaultOptions()
	findings := linter.Lint("my.env", entries, opts)
	for _, f := range findings {
		if f.File != "my.env" {
			t.Errorf("expected file my.env, got %q", f.File)
		}
	}
}
