package differ_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
)

func TestDiff_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	hunks := differ.Diff("a.env", "b.env", a, b)
	if len(hunks) != 0 {
		t.Fatalf("expected 0 hunks, got %d", len(hunks))
	}
}

func TestDiff_Added(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "NEW": "val"}
	hunks := differ.Diff("a.env", "b.env", a, b)
	if len(hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(hunks))
	}
	if !hunks[0].Added {
		t.Error("expected Added hunk")
	}
	if hunks[0].Key != "NEW" {
		t.Errorf("expected key NEW, got %s", hunks[0].Key)
	}
}

func TestDiff_Removed(t *testing.T) {
	a := map[string]string{"FOO": "bar", "OLD": "gone"}
	b := map[string]string{"FOO": "bar"}
	hunks := differ.Diff("a.env", "b.env", a, b)
	if len(hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(hunks))
	}
	if !hunks[0].Removed {
		t.Error("expected Removed hunk")
	}
	if hunks[0].Key != "OLD" {
		t.Errorf("expected key OLD, got %s", hunks[0].Key)
	}
}

func TestDiff_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	hunks := differ.Diff("a.env", "b.env", a, b)
	if len(hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(hunks))
	}
	if !hunks[0].Changed {
		t.Error("expected Changed hunk")
	}
	if hunks[0].ValueA != "old" || hunks[0].ValueB != "new" {
		t.Errorf("unexpected values: %s / %s", hunks[0].ValueA, hunks[0].ValueB)
	}
}

func TestDiff_SortedByKey(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "2", "M": "3"}
	b := map[string]string{}
	hunks := differ.Diff("a.env", "b.env", a, b)
	keys := make([]string, len(hunks))
	for i, h := range hunks {
		keys[i] = h.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("hunks not sorted: %v", keys)
		}
	}
}

func TestHunk_String_Added(t *testing.T) {
	h := differ.Hunk{Key: "X", FileB: "b.env", ValueB: "val", Added: true}
	if !strings.HasPrefix(h.String(), "+") {
		t.Errorf("expected '+' prefix, got: %s", h.String())
	}
}

func TestHunk_String_Removed(t *testing.T) {
	h := differ.Hunk{Key: "X", FileA: "a.env", ValueA: "val", Removed: true}
	if !strings.HasPrefix(h.String(), "-") {
		t.Errorf("expected '-' prefix, got: %s", h.String())
	}
}

func TestSummary_NoChanges(t *testing.T) {
	s := differ.Summary([]differ.Hunk{})
	if s != "no differences" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_Mixed(t *testing.T) {
	hunks := []differ.Hunk{
		{Added: true},
		{Removed: true},
		{Changed: true},
		{Changed: true},
	}
	s := differ.Summary(hunks)
	if !strings.Contains(s, "1 added") {
		t.Errorf("missing added count: %s", s)
	}
	if !strings.Contains(s, "1 removed") {
		t.Errorf("missing removed count: %s", s)
	}
	if !strings.Contains(s, "2 changed") {
		t.Errorf("missing changed count: %s", s)
	}
}
