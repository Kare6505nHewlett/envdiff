// Package normalizer provides key-value normalization for .env entries.
//
// It supports trimming whitespace from keys and values, optional
// case-folding of keys, and canonicalization of boolean-like values
// (e.g. "yes", "1", "on" → "true"; "no", "0", "off" → "false").
//
// Normalization is applied before diff or comparison operations to reduce
// false positives caused by superficial formatting differences.
//
// Usage:
//
//	opts := normalizer.DefaultOptions()
//	normalized := normalizer.Apply(env, opts)
package normalizer
