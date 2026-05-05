// Package linter implements style and convention checks for .env file key/value
// pairs. It is designed to catch common authoring mistakes such as lowercase
// environment variable names, keys containing spaces, or values that are
// suspiciously long.
//
// # Lint Rules
//
// The following rules are applied by default:
//
//   - Keys must be uppercase (e.g. DATABASE_URL, not database_url)
//   - Keys must not contain spaces
//   - Values must not exceed a configurable maximum length
//   - Keys must not be empty
//
// # Usage
//
//	opts := linter.DefaultOptions()
//	findings := linter.Lint("production.env", entries, opts)
//	linter.ReportText(os.Stdout, findings)
//
// Lint rules can be selectively disabled via the Options struct. The reporter
// helpers support both plain-text and JSON output formats.
package linter
