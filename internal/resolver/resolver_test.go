package resolver

import (
	"testing"
)

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	got, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["HOST"] != "localhost" || got["PORT"] != "5432" {
		t.Errorf("expected unchanged values, got %v", got)
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE": "postgres",
		"URL":  "${BASE}://localhost",
	}
	got, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["URL"] != "postgres://localhost" {
		t.Errorf("expected expanded URL, got %q", got["URL"])
	}
}

func TestResolve_DollarStyle(t *testing.T) {
	env := map[string]string{
		"SCHEME": "https",
		"DOMAIN": "example.com",
		"FULL":   "$SCHEME://$DOMAIN",
	}
	got, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FULL"] != "https://example.com" {
		t.Errorf("unexpected value: %q", got["FULL"])
	}
}

func TestResolve_NestedReference(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	got, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["C"] != "hello_world!" {
		t.Errorf("nested expansion failed: %q", got["C"])
	}
}

func TestResolve_MissingVarSilent(t *testing.T) {
	env := map[string]string{
		"URL": "${UNDEFINED}/path",
	}
	got, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// unresolved reference left as-is
	if got["URL"] != "${UNDEFINED}/path" {
		t.Errorf("expected original value, got %q", got["URL"])
	}
}

func TestResolve_FailOnMissing(t *testing.T) {
	env := map[string]string{
		"URL": "${MISSING}/path",
	}
	opts := DefaultOptions()
	opts.FailOnMissing = true
	_, err := Resolve(env, opts)
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestResolve_MaxDepthExceeded(t *testing.T) {
	// A -> B -> A creates a cycle; MaxDepth should stop it.
	env := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	opts := DefaultOptions()
	opts.MaxDepth = 3
	_, err := Resolve(env, opts)
	if err == nil {
		t.Fatal("expected error for max depth exceeded")
	}
}
