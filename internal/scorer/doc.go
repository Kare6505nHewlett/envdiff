// Package scorer computes a numeric health score (0–100) that summarises
// the overall parity between compared .env files.
//
// A score of 100 means every key matches across all environments.
// Missing keys incur a full-weight penalty; mismatched values incur a
// half-weight penalty by default.  Both weights are configurable via
// the Weights struct so callers can tune sensitivity to their needs.
//
// Typical usage:
//
//	results := diff.Compare(base, target)
//	sr := scorer.Score(results, nil)
//	fmt.Printf("Parity score: %.2f/100\n", sr.Score)
package scorer
