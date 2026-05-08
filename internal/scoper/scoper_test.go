package scoper_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scoper"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", File: "production.env", Status: diff.StatusMatch},
		{Key: "API_KEY", File: "production.env", Status: diff.StatusMissing},
		{Key: "DB_HOST", File: "staging.env", Status: diff.StatusMismatch},
		{Key: "CACHE_URL", File: "local.env", Status: diff.StatusMatch},
	}
}

func TestApply_PartitionsIntoScopes(t *testing.T) {
	scopes := []scoper.Scope{
		{Name: "production", Files: []string{"production.env"}},
		{Name: "staging", Files: []string{"staging.env"}},
	}
	results := Apply(sampleResults(), scopes)
	if len(results) != 3 {
		t.Fatalf("expected 3 scope results (prod, staging, default), got %d", len(results))
	}
}

func TestApply_DefaultScopeForUnmatchedFiles(t *testing.T) {
	scopes := []scoper.Scope{
		{Name: "production", Files: []string{"production.env"}},
	}
	results := scoper.Apply(sampleResults(), scopes)
	var def *scoper.Result
	for i := range results {
		if results[i].Scope == "default" {
			def = &results[i]
			break
		}
	}
	if def == nil {
		t.Fatal("expected a default scope")
	}
	if len(def.Entries) != 2 {
		t.Errorf("expected 2 default entries, got %d", len(def.Entries))
	}
}

func TestApply_EmptyScopes(t *testing.T) {
	results := scoper.Apply(sampleResults(), nil)
	if len(results) != 1 || results[0].Scope != "default" {
		t.Errorf("expected single default scope, got %v", results)
	}
}

func TestApply_EmptyResults(t *testing.T) {
	scopes := []scoper.Scope{{Name: "production", Files: []string{"production.env"}}}
	results := scoper.Apply(nil, scopes)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestNames_ReturnsSorted(t *testing.T) {
	input := []scoper.Result{
		{Scope: "staging"},
		{Scope: "production"},
		{Scope: "default"},
	}
	names := scoper.Names(input)
	expected := []string{"default", "production", "staging"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("index %d: expected %q got %q", i, expected[i], n)
		}
	}
}

func TestApply_CaseInsensitiveFileMatch(t *testing.T) {
	scopes := []scoper.Scope{
		{Name: "prod", Files: []string{"Production.Env"}},
	}
	input := []diff.Result{
		{Key: "X", File: "production.env", Status: diff.StatusMatch},
	}
	results := scoper.Apply(input, scopes)
	if len(results) != 1 || results[0].Scope != "prod" {
		t.Errorf("expected prod scope, got %v", results)
	}
}
