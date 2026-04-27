package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", Status: diff.StatusMatch, File: "b.env"},
		{Key: "ALPHA", Status: diff.StatusMissing, File: "a.env"},
		{Key: "MANGO", Status: diff.StatusMismatch, File: "a.env"},
		{Key: "BETA", Status: diff.StatusMatch, File: "c.env"},
	}
}

func TestSort_ByKeyAscending(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.Options{Field: sorter.SortByKey, Order: sorter.Ascending})

	expected := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_ByKeyDescending(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.Options{Field: sorter.SortByKey, Order: sorter.Descending})

	expected := []string{"ZEBRA", "MANGO", "BETA", "ALPHA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_ByStatus(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.Options{Field: sorter.SortByStatus, Order: sorter.Ascending})

	if sorted[0].Status != diff.StatusMissing {
		t.Errorf("expected first result to be Missing, got %q", sorted[0].Status)
	}
	if sorted[1].Status != diff.StatusMismatch {
		t.Errorf("expected second result to be Mismatch, got %q", sorted[1].Status)
	}
}

func TestSort_ByFile(t *testing.T) {
	results := sampleResults()
	sorted := sorter.Sort(results, sorter.Options{Field: sorter.SortByFile, Order: sorter.Ascending})

	if sorted[0].File != "a.env" {
		t.Errorf("expected first file to be a.env, got %q", sorted[0].File)
	}
	if sorted[len(sorted)-1].File != "c.env" {
		t.Errorf("expected last file to be c.env, got %q", sorted[len(sorted)-1].File)
	}
}

func TestSort_EmptyInput(t *testing.T) {
	results := []diff.Result{}
	sorted := sorter.Sort(results, sorter.Options{Field: sorter.SortByKey, Order: sorter.Ascending})
	if len(sorted) != 0 {
		t.Errorf("expected empty result, got %d items", len(sorted))
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	results := sampleResults()
	originalFirst := results[0].Key
	sorter.Sort(results, sorter.Options{Field: sorter.SortByKey, Order: sorter.Ascending})
	if results[0].Key != originalFirst {
		t.Errorf("original slice was mutated: got %q, want %q", results[0].Key, originalFirst)
	}
}
