// Package snapshotter captures and compares point-in-time snapshots of diff
// results, enabling detection of changes between runs.
package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot holds a saved set of diff results with metadata.
type Snapshot struct {
	CreatedAt time.Time    `json:"created_at"`
	Label     string       `json:"label"`
	Results   []diff.Result `json:"results"`
}

// Delta describes what changed between two snapshots.
type Delta struct {
	Added   []diff.Result `json:"added"`
	Removed []diff.Result `json:"removed"`
	Changed []diff.Result `json:"changed"`
}

// Save writes a snapshot to the given file path.
func Save(path, label string, results []diff.Result) error {
	s := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshotter: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshotter: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshotter: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshotter: unmarshal: %w", err)
	}
	if s.Results == nil {
		s.Results = []diff.Result{}
	}
	return s, nil
}

// Compare returns the delta between an old and new snapshot.
func Compare(old, new Snapshot) Delta {
	oldMap := indexResults(old.Results)
	newMap := indexResults(new.Results)

	var delta Delta
	for key, newR := range newMap {
		oldR, exists := oldMap[key]
		if !exists {
			delta.Added = append(delta.Added, newR)
		} else if oldR.Status != newR.Status {
			delta.Changed = append(delta.Changed, newR)
		}
	}
	for key, oldR := range oldMap {
		if _, exists := newMap[key]; !exists {
			delta.Removed = append(delta.Removed, oldR)
		}
	}
	return delta
}

func indexResults(results []diff.Result) map[string]diff.Result {
	m := make(map[string]diff.Result, len(results))
	for _, r := range results {
		m[r.Key+"|"+r.File] = r
	}
	return m
}
