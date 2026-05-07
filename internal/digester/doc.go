// Package digester provides structural fingerprinting for .env file key sets.
//
// A Digest is a SHA-256 hash derived exclusively from sorted key names,
// deliberately omitting values so that sensitive data is never hashed.
// This makes digests safe to log, store, or transmit for drift detection.
//
// Typical usage:
//
//	results := diff.Compare(envA, envB)
//	digests := digester.Compute(results)
//	if !digester.Match(digests) {
//		fmt.Print(digester.Diff(digests))
//	}
package digester
