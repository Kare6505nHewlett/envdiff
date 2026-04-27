package diff

import (
	"testing"
)

func TestCompare_AllMatch(t *testing.T) {
	ref := map[string]string{"FOO": "bar", "BAZ": "qux"}
	tgt := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Compare("ref.env", ref, "target.env", tgt)
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	for _, e := range result.Entries {
		if e.Status != StatusMatch {
			t.Errorf("key %q: expected match, got %s", e.Key, e.Status)
		}
	}
}

func TestCompare_MissingKey(t *testing.T) {
	ref := map[string]string{"FOO": "bar", "MISSING": "val"}
	tgt := map[string]string{"FOO": "bar"}

	result := Compare("ref.env", ref, "target.env", tgt)
	found := false
	for _, e := range result.Entries {
		if e.Key == "MISSING" {
			found = true
			if e.Status != StatusMissing {
				t.Errorf("expected missing, got %s", e.Status)
			}
		}
	}
	if !found {
		t.Error("expected MISSING key in entries")
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	ref := map[string]string{"FOO": "bar"}
	tgt := map[string]string{"FOO": "bar", "EXTRA": "bonus"}

	result := Compare("ref.env", ref, "target.env", tgt)
	found := false
	for _, e := range result.Entries {
		if e.Key == "EXTRA" {
			found = true
			if e.Status != StatusExtra {
				t.Errorf("expected extra, got %s", e.Status)
			}
		}
	}
	if !found {
		t.Error("expected EXTRA key in entries")
	}
}

func TestCompare_Mismatch(t *testing.T) {
	ref := map[string]string{"FOO": "original"}
	tgt := map[string]string{"FOO": "changed"}

	result := Compare("ref.env", ref, "target.env", tgt)
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	e := result.Entries[0]
	if e.Status != StatusMismatch {
		t.Errorf("expected mismatch, got %s", e.Status)
	}
	if e.RefValue != "original" || e.TargetValue != "changed" {
		t.Errorf("unexpected values: ref=%q target=%q", e.RefValue, e.TargetValue)
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	ref := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	tgt := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}

	result := Compare("ref.env", ref, "target.env", tgt)
	keys := make([]string, len(result.Entries))
	for i, e := range result.Entries {
		keys[i] = e.Key
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestResult_Summary(t *testing.T) {
	result := Result{
		Entries: []Entry{
			{Key: "A", Status: StatusMatch},
			{Key: "B", Status: StatusMissing},
			{Key: "C", Status: StatusMismatch},
			{Key: "D", Status: StatusExtra},
			{Key: "E", Status: StatusMissing},
		},
	}
	s := result.Summary()
	if s[StatusMatch] != 1 || s[StatusMissing] != 2 || s[StatusMismatch] != 1 || s[StatusExtra] != 1 {
		t.Errorf("unexpected summary: %v", s)
	}
}
