package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/renamer"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", File: "prod.env", Status: diff.StatusMissing},
		{Key: "DB_PASS", File: "prod.env", Status: diff.StatusMismatch},
		{Key: "API_KEY", File: "staging.env", Status: diff.StatusMissing},
	}
}

func TestApply_NoRules(t *testing.T) {
	results := renamer.Apply(sampleResults(), nil, false)
	if results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestApply_NoResults(t *testing.T) {
	rules := []renamer.Rule{{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"}}
	out := renamer.Apply(nil, rules, false)
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestApply_MatchesKey(t *testing.T) {
	rules := []renamer.Rule{{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"}}
	out := renamer.Apply(sampleResults(), rules, false)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Rule.NewKey != "DATABASE_HOST" {
		t.Errorf("unexpected new key: %s", out[0].Rule.NewKey)
	}
	if len(out[0].Matched) != 1 {
		t.Errorf("expected 1 matched result, got %d", len(out[0].Matched))
	}
}

func TestApply_IgnoreCase(t *testing.T) {
	rules := []renamer.Rule{{OldKey: "db_host", NewKey: "DATABASE_HOST"}}
	out := renamer.Apply(sampleResults(), rules, true)
	if len(out) != 1 {
		t.Fatalf("expected 1 match with ignoreCase=true, got %d", len(out))
	}
}

func TestApply_NoMatch(t *testing.T) {
	rules := []renamer.Rule{{OldKey: "UNKNOWN_KEY", NewKey: "OTHER_KEY"}}
	out := renamer.Apply(sampleResults(), rules, false)
	if len(out) != 0 {
		t.Errorf("expected 0 matches, got %d", len(out))
	}
}

func TestLoadRules_Valid(t *testing.T) {
	lines := []string{
		"# comment",
		"",
		"DB_HOST=DATABASE_HOST",
		"API_KEY=SERVICE_API_KEY",
	}
	rules := renamer.LoadRules(lines)
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].OldKey != "DB_HOST" || rules[0].NewKey != "DATABASE_HOST" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].OldKey != "API_KEY" || rules[1].NewKey != "SERVICE_API_KEY" {
		t.Errorf("unexpected rule[1]: %+v", rules[1])
	}
}

func TestLoadRules_SkipsMalformed(t *testing.T) {
	lines := []string{"NOEQUALSIGN", "VALID=KEY"}
	rules := renamer.LoadRules(lines)
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}
