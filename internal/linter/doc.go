// Package linter implements style and convention checks for .env file key/value
// pairs. It is designed to catch common authoring mistakes such as lowercase
// environment variable names, keys containing spaces, or values that are
// suspiciously long.
//
// Usage:
//
//	opts := linter.DefaultOptions()
//	findings := linter.Lint("production.env", entries, opts)
//	linter.ReportText(os.Stdout, findings)
//
// Lint rules can be selectively disabled via the Options struct. The reporter
// helpers support both plain-text and JSON output formats.
package linter
