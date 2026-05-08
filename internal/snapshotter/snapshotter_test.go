package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshotter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", File: "prod.env", Status: "match"},
		{Key: "API_KEY", File: "prod.env", Status: "missing"},
		{Key: "TIMEOUT", File: "prod.env", Status: "mismatch"},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	err := snapshotter.Save(path, "test-label", sampleResults())
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to exist")
	}
}

func TestLoad_ReturnsCorrectData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	_ = snapshotter.Save(path, "my-label", sampleResults())
	snap, err := snapshotter.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if snap.Label != "my-label" {
		t.Errorf("expected label 'my-label', got %q", snap.Label)
	}
	if len(snap.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(snap.Results))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshotter.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	results := sampleResults()

	_ = snapshotter.Save(path, "round-trip", results)
	snap, err := snapshotter.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	for i, r := range snap.Results {
		if r.Key != results[i].Key || r.Status != results[i].Status || r.File != results[i].File {
			t.Errorf("result[%d] mismatch: got %+v, want %+v", i, r, results[i])
		}
	}
}

func TestCompare_Added(t *testing.T) {
	old := snapshotter.Snapshot{Results: sampleResults()[:1]}
	new := snapshotter.Snapshot{Results: sampleResults()}

	delta := snapshotter.Compare(old, new)
	if len(delta.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(delta.Added))
	}
}

func TestCompare_Removed(t *testing.T) {
	old := snapshotter.Snapshot{Results: sampleResults()}
	new := snapshotter.Snapshot{Results: sampleResults()[:1]}

	delta := snapshotter.Compare(old, new)
	if len(delta.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(delta.Removed))
	}
}

func TestCompare_Changed(t *testing.T) {
	old := snapshotter.Snapshot{Results: []diff.Result{
		{Key: "DB_HOST", File: "prod.env", Status: "missing"},
	}}
	new := snapshotter.Snapshot{Results: []diff.Result{
		{Key: "DB_HOST", File: "prod.env", Status: "match"},
	}}

	delta := snapshotter.Compare(old, new)
	if len(delta.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(delta.Changed))
	}
}

func TestCompare_NoChanges(t *testing.T) {
	old := snapshotter.Snapshot{Results: sampleResults()}
	new := snapshotter.Snapshot{Results: sampleResults()}

	delta := snapshotter.Compare(old, new)
	if len(delta.Added)+len(delta.Removed)+len(delta.Changed) != 0 {
		t.Errorf("expected no changes, got delta: %+v", delta)
	}
}
