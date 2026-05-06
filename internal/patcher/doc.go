// Package patcher provides utilities for applying targeted key-value patches
// to .env files.
//
// It supports two modes of operation:
//
//   - Apply: reads the target file, updates existing keys in-place, appends
//     missing keys, and writes the result back to disk.
//
//   - DryRun: performs the same analysis but returns results without modifying
//     any files, useful for previewing changes before committing them.
//
// Original file formatting (comments, blank lines) is preserved for lines
// that are not patched.
package patcher
