// Package masker partially masks sensitive values in diff results so that
// operators can verify the presence and rough shape of a secret without
// exposing the full plaintext.
//
// A key is considered sensitive when its name (case-insensitive) contains
// any of the configured substrings such as "SECRET", "PASSWORD", "TOKEN",
// "KEY", "PRIVATE", or "CREDENTIAL".
//
// Example:
//
//	results := diff.Compare(envA, envB)
//	masked := masker.Apply(results, masker.DefaultOptions())
//	// DB_PASSWORD value "supersecret" becomes "sup*****ret"
//
// The original result slice is never modified; Apply always returns a
// new slice with copied entries.
package masker
