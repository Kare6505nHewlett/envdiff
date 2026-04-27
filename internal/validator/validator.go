package validator

import (
	"fmt"
	"strings"
)

// Rule defines a validation rule for environment variable keys or values.
type Rule struct {
	RequiredPrefix string
	ForbiddenKeys  []string
	NoEmptyValues  bool
}

// Violation represents a single validation issue found in an env file.
type Violation struct {
	File    string
	Key     string
	Message string
}

// Validate checks a map of env key-value pairs against the given Rule and
// returns a slice of Violations describing any issues found.
func Validate(file string, env map[string]string, rule Rule) []Violation {
	var violations []Violation

	forbidden := make(map[string]bool, len(rule.ForbiddenKeys))
	for _, k := range rule.ForbiddenKeys {
		forbidden[strings.ToUpper(k)] = true
	}

	for key, value := range env {
		upper := strings.ToUpper(key)

		if forbidden[upper] {
			violations = append(violations, Violation{
				File:    file,
				Key:     key,
				Message: fmt.Sprintf("key %q is forbidden", key),
			})
		}

		if rule.RequiredPrefix != "" && !strings.HasPrefix(upper, strings.ToUpper(rule.RequiredPrefix)) {
			violations = append(violations, Violation{
				File:    file,
				Key:     key,
				Message: fmt.Sprintf("key %q does not have required prefix %q", key, rule.RequiredPrefix),
			})
		}

		if rule.NoEmptyValues && strings.TrimSpace(value) == "" {
			violations = append(violations, Violation{
				File:    file,
				Key:     key,
				Message: fmt.Sprintf("key %q has an empty value", key),
			})
		}
	}

	return violations
}
