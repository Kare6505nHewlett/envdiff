// Package templater generates a .env.template file from a set of diff results.
//
// It collects all unique keys found across compared environment files and
// writes them out with redacted or placeholder values, making it easy to
// bootstrap a new environment configuration without leaking real secrets.
//
// Usage:
//
//	results := diff.Compare(base, other, "staging.env")
//	err := templater.WriteFile(".env.template", results, templater.Options{
//		Placeholder:     "CHANGEME",
//		IncludeComments: true,
//	})
package templater
