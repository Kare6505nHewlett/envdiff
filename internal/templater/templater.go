// Package templater generates a .env.template file from diff results,
// producing a file with keys but redacted or empty values.
package templater

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options controls template generation behaviour.
type Options struct {
	// Placeholder is written as the value for each key.
	// Defaults to "" (empty) if not set.
	Placeholder string
	// IncludeComments adds a comment above each key indicating its diff status.
	IncludeComments bool
}

// Generate writes a .env template derived from results to w.
// Keys are sorted alphabetically. Each key appears exactly once.
func Generate(w io.Writer, results []diff.Result, opts Options) error {
	if opts.Placeholder == "" {
		opts.Placeholder = ""
	}

	seen := make(map[string]diff.Result)
	for _, r := range results {
		if _, exists := seen[r.Key]; !exists {
			seen[r.Key] = r
		}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		r := seen[k]
		if opts.IncludeComments {
			status := strings.ToUpper(string(r.Status))
			if _, err := fmt.Fprintf(w, "# status: %s\n", status); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, opts.Placeholder); err != nil {
			return err
		}
	}
	return nil
}

// WriteFile generates a template and writes it to the given file path.
func WriteFile(path string, results []diff.Result, opts Options) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("templater: create %q: %w", path, err)
	}
	defer f.Close()
	return Generate(f, results, opts)
}
