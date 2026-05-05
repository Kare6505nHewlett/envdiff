// Package linter provides style and convention checks for .env file entries.
// It flags issues such as lowercase keys, keys with spaces, or suspiciously
// long values that may indicate accidental data inclusion.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Severity represents the importance of a lint finding.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Finding describes a single lint violation.
type Finding struct {
	Key      string
	File     string
	Rule     string
	Message  string
	Severity Severity
}

// Options controls which lint rules are active.
type Options struct {
	DisallowLowercase bool
	DisallowSpaces    bool
	MaxValueLength    int // 0 means unlimited
}

// DefaultOptions returns a sensible default lint configuration.
func DefaultOptions() Options {
	return Options{
		DisallowLowercase: true,
		DisallowSpaces:    true,
		MaxValueLength:    256,
	}
}

// Lint checks a map of key/value pairs parsed from file and returns findings.
func Lint(file string, entries map[string]string, opts Options) []Finding {
	var findings []Finding

	for key, value := range entries {
		if opts.DisallowLowercase && hasLowercaseLetter(key) {
			findings = append(findings, Finding{
				Key:      key,
				File:     file,
				Rule:     "lowercase-key",
				Message:  fmt.Sprintf("key %q contains lowercase letters", key),
				Severity: SeverityWarn,
			})
		}

		if opts.DisallowSpaces && strings.ContainsRune(key, ' ') {
			findings = append(findings, Finding{
				Key:      key,
				File:     file,
				Rule:     "space-in-key",
				Message:  fmt.Sprintf("key %q contains a space character", key),
				Severity: SeverityError,
			})
		}

		if opts.MaxValueLength > 0 && len(value) > opts.MaxValueLength {
			findings = append(findings, Finding{
				Key:      key,
				File:     file,
				Rule:     "value-too-long",
				Message:  fmt.Sprintf("value for key %q exceeds max length %d (got %d)", key, opts.MaxValueLength, len(value)),
				Severity: SeverityWarn,
			})
		}
	}

	return findings
}

func hasLowercaseLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}
