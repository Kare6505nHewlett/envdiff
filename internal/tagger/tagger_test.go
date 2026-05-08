package tagger_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/tagger"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMismatch},
		{Key: "DB_PORT", Status: diff.StatusMatch},
		{Key: "AWS_ACCESS_KEY", Status: diff.StatusMissing},
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "UNKNOWN_VAR", Status: diff.StatusMismatch},
	}
}

func writeRulesFile(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".envtags")
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeRulesFile: %v", err)
	}
	return path
}

func TestApply_NoRules(t *testing.T) {
	results := sampleResults()
	tagged := tagger.Apply(results, nil)
	if len(tagged) != len(results) {
		t.Fatalf("expected %d tagged results, got %d", len(results), len(tagged))
	}
	for _, tr := range tagged {
		if tr.Tag != "untagged" {
			t.Errorf("key %q: expected 'untagged', got %q", tr.Key, tr.Tag)
		}
	}
}

func TestApply_PrefixRule(t *testing.T) {
	rules := []tagger.Rule{
		{Pattern: "DB_", Tag: "database"},
		{Pattern: "AWS_", Tag: "cloud"},
	}
	tagged := tagger.Apply(sampleResults(), rules)
	expect := map[string]string{
		"DB_HOST":       "database",
		"DB_PORT":       "database",
		"AWS_ACCESS_KEY": "cloud",
		"APP_NAME":      "untagged",
		"UNKNOWN_VAR":   "untagged",
	}
	for _, tr := range tagged {
		if got := tr.Tag; got != expect[tr.Key] {
			t.Errorf("key %q: want tag %q, got %q", tr.Key, expect[tr.Key], got)
		}
	}
}

func TestApply_ExactRule(t *testing.T) {
	rules := []tagger.Rule{
		{Pattern: "APP_NAME", Tag: "identity", Exact: true},
		{Pattern: "APP_", Tag: "app"},
	}
	tagged := tagger.Apply(sampleResults(), rules)
	for _, tr := range tagged {
		if tr.Key == "APP_NAME" && tr.Tag != "identity" {
			t.Errorf("exact rule should win: got %q", tr.Tag)
		}
	}
}

func TestLoadRules_ValidFile(t *testing.T) {
	path := writeRulesFile(t, []string{
		"# comment",
		"",
		"DB_\tdatabase",
		"=AWS_ACCESS_KEY\tcloud",
	})
	rules, err := tagger.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Tag != "database" || rules[0].Exact {
		t.Errorf("rule[0] unexpected: %+v", rules[0])
	}
	if rules[1].Tag != "cloud" || !rules[1].Exact {
		t.Errorf("rule[1] unexpected: %+v", rules[1])
	}
}

func TestLoadRules_MissingFileReturnsNil(t *testing.T) {
	rules, err := tagger.LoadRules(filepath.Join(t.TempDir(), "nonexistent"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if rules != nil {
		t.Errorf("expected nil rules, got %v", rules)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	tagged := tagger.Apply(nil, []tagger.Rule{{Pattern: "DB_", Tag: "database"}})
	if len(tagged) != 0 {
		t.Errorf("expected empty slice, got %d", len(tagged))
	}
}

func TestApply_PreservesResultFields(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", File: "prod.env", Status: diff.StatusMismatch, Values: map[string]string{"prod.env": "localhost"}},
	}
	rules := []tagger.Rule{{Pattern: "DB_", Tag: "database"}}
	tagged := tagger.Apply(results, rules)
	if tagged[0].File != "prod.env" {
		t.Errorf("File field not preserved: %v", tagged[0].File)
	}
	if fmt.Sprint(tagged[0].Values) != fmt.Sprint(results[0].Values) {
		t.Errorf("Values field not preserved")
	}
}
