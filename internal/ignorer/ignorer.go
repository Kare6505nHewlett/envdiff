// Package ignorer provides functionality to load and apply .envdiffignore
// files, allowing users to suppress specific keys from diff results.
package ignorer

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Ignorer holds a set of keys and prefix patterns to ignore.
type Ignorer struct {
	exact    map[string]struct{}
	prefixes []string
}

// LoadFile reads an ignore file from the given path and returns an Ignorer.
// Lines starting with '#' or empty lines are skipped.
// Lines ending with '*' are treated as prefix patterns.
func LoadFile(path string) (*Ignorer, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Ignorer{exact: make(map[string]struct{})}, nil
		}
		return nil, err
	}
	defer f.Close()

	ig := &Ignorer{exact: make(map[string]struct{})}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "*") {
			ig.prefixes = append(ig.prefixes, strings.TrimSuffix(line, "*"))
		} else {
			ig.exact[line] = struct{}{}
		}
	}
	return ig, scanner.Err()
}

// Apply filters out any diff results whose key matches an ignored key or prefix.
func (ig *Ignorer) Apply(results []diff.Result) []diff.Result {
	if ig == nil {
		return results
	}
	filtered := make([]diff.Result, 0, len(results))
	for _, r := range results {
		if !ig.matches(r.Key) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// matches reports whether a key should be ignored.
func (ig *Ignorer) matches(key string) bool {
	if _, ok := ig.exact[key]; ok {
		return true
	}
	for _, p := range ig.prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
