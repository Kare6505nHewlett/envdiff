// Package snapshotter provides functionality to save, load, and compare
// point-in-time snapshots of envdiff results.
//
// A Snapshot captures the full set of diff.Result entries at a given moment,
// along with a label and timestamp. Snapshots can be persisted to disk as JSON
// and later reloaded for comparison.
//
// Use Compare to compute a Delta between two snapshots, identifying keys that
// were added, removed, or changed status between runs. This is useful for
// tracking environment drift over time or in CI pipelines.
//
// Example usage:
//
//	snapshotter.Save("baseline.snap.json", "v1.0", results)
//	old, _ := snapshotter.Load("baseline.snap.json")
//	new := snapshotter.Snapshot{Results: currentResults}
//	delta := snapshotter.Compare(old, new)
package snapshotter
