// Package renamer provides utilities for suggesting or applying key renames
// across diff results, mapping old key names to new ones.
package renamer

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Rule defines a single rename mapping from an old key name to a new one.
type Rule struct {
	OldKey string
	NewKey string
}

// Result holds the outcome of applying a rename rule to a set of diff results.
type Result struct {
	Rule    Rule
	Matched []diff.Result
}

// Apply scans the provided diff results and returns rename suggestions based on
// the given rules. A result is matched if its key equals Rule.OldKey (case-insensitive
// when ignoreCase is true).
func Apply(results []diff.Result, rules []Rule, ignoreCase bool) []Result {
	if len(results) == 0 || len(rules) == 0 {
		return nil
	}

	ruleMap := make(map[string]Rule, len(rules))
	for _, r := range rules {
		lookup := r.OldKey
		if ignoreCase {
			lookup = strings.ToUpper(lookup)
		}
		ruleMap[lookup] = r
	}

	matches := make(map[string]*Result)

	for _, res := range results {
		lookup := res.Key
		if ignoreCase {
			lookup = strings.ToUpper(lookup)
		}
		if rule, ok := ruleMap[lookup]; ok {
			if _, exists := matches[rule.OldKey]; !exists {
				matches[rule.OldKey] = &Result{Rule: rule}
			}
			matches[rule.OldKey].Matched = append(matches[rule.OldKey].Matched, res)
		}
	}

	out := make([]Result, 0, len(matches))
	for _, v := range matches {
		out = append(out, *v)
	}
	return out
}

// LoadRules parses a slice of "OLD_KEY=NEW_KEY" strings into Rule values.
// Lines that are blank or start with '#' are skipped.
func LoadRules(lines []string) []Rule {
	var rules []Rule
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		rules = append(rules, Rule{
			OldKey: strings.TrimSpace(parts[0]),
			NewKey: strings.TrimSpace(parts[1]),
		})
	}
	return rules
}
