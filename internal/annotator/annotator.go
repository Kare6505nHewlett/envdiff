// Package annotator attaches human-readable descriptions to diff results
// by reading a companion .env.descriptions file or an inline comment map.
package annotator

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Annotation holds a description for a single key.
type Annotation struct {
	Key         string
	Description string
}

// AnnotatedResult wraps a diff.Result with an optional description.
type AnnotatedResult struct {
	diff.Result
	Description string
}

// LoadDescriptions parses a file where each line is "KEY=description text".
// Lines starting with '#' and blank lines are ignored.
func LoadDescriptions(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	desc := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if key != "" {
			desc[key] = val
		}
	}
	return desc, scanner.Err()
}

// Apply merges a descriptions map into a slice of diff.Result values,
// returning AnnotatedResult entries. Results without a description get
// an empty Description field.
func Apply(results []diff.Result, descriptions map[string]string) []AnnotatedResult {
	out := make([]AnnotatedResult, 0, len(results))
	for _, r := range results {
		ar := AnnotatedResult{Result: r}
		if descriptions != nil {
			ar.Description = descriptions[r.Key]
		}
		out = append(out, ar)
	}
	return out
}
