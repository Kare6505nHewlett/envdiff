// Package resolver resolves environment variable values by expanding
// references to other variables within the same file or across files.
package resolver

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls resolver behaviour.
type Options struct {
	// MaxDepth limits recursive expansion to prevent infinite loops.
	MaxDepth int
	// FailOnMissing returns an error when a referenced variable is not found.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxDepth:      10,
		FailOnMissing: false,
	}
}

// Resolve expands variable references in values using the provided env map.
// It returns a new map with all resolvable references expanded.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	resolved := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := expand(v, env, opts, 0)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", k, err)
		}
		resolved[k] = expanded
	}
	return resolved, nil
}

// expand recursively resolves variable references in a single value.
func expand(value string, env map[string]string, opts Options, depth int) (string, error) {
	if depth > opts.MaxDepth {
		return value, fmt.Errorf("max expansion depth %d exceeded", opts.MaxDepth)
	}

	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		replacement, ok := env[name]
		if !ok {
			if opts.FailOnMissing {
				expandErr = fmt.Errorf("undefined variable %q", name)
			}
			return match
		}
		inner, err := expand(replacement, env, opts, depth+1)
		if err != nil {
			expandErr = err
			return match
		}
		return inner
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// extractName strips ${ } or $ from a variable reference.
func extractName(ref string) string {
	ref = strings.TrimPrefix(ref, "${") 
	ref = strings.TrimSuffix(ref, "}")
	ref = strings.TrimPrefix(ref, "$")
	return ref
}
