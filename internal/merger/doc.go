// Package merger collapses per-file diff.Result slices into a unified
// MergedKey view that maps each key to its values and overall status
// across all compared environment files.
//
// Usage:
//
//	results := diff.Compare(baseEnv, otherEnvs...)
//	merged := merger.Merge(results)
//	for _, mk := range merged {
//		fmt.Printf("%s [%s]\n", mk.Key, mk.Status)
//	}
package merger
