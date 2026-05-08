// Package tagger assigns user-defined tags to diff results based on key
// patterns, enabling downstream grouping, filtering, and reporting by tag.
package tagger

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Rule maps a key prefix or exact key name to a tag label.
type Rule struct {
	Pattern string
	Tag     string
	Exact   bool
}

// TaggedResult wraps a diff.Result with an associated tag.
type TaggedResult struct {
	diff.Result
	Tag string
}

// Apply tags each result according to the first matching rule. Results that
// match no rule receive the tag "untagged".
func Apply(results []diff.Result, rules []Rule) []TaggedResult {
	tagged := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		tag := matchTag(r.Key, rules)
		tagged = append(tagged, TaggedResult{Result: r, Tag: tag})
	}
	return tagged
}

// LoadRules reads a tag rules file. Each non-blank, non-comment line must
// follow the format:
//
//	<pattern>\t<tag>
//
// Prefix a pattern with "=" to require an exact key match.
func LoadRules(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var rules []Rule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		pattern := strings.TrimSpace(parts[0])
		tag := strings.TrimSpace(parts[1])
		exact := strings.HasPrefix(pattern, "=")
		if exact {
			pattern = pattern[1:]
		}
		rules = append(rules, Rule{Pattern: pattern, Tag: tag, Exact: exact})
	}
	return rules, scanner.Err()
}

func matchTag(key string, rules []Rule) string {
	for _, r := range rules {
		if r.Exact {
			if key == r.Pattern {
				return r.Tag
			}
		} else {
			if strings.HasPrefix(key, r.Pattern) {
				return r.Tag
			}
		}
	}
	return "untagged"
}
