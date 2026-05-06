// Package patcher applies a set of key-value patches to an existing .env file,
// adding missing keys and updating mismatched ones while preserving formatting.
package patcher

import (
	"fmt"
	"os"
	"strings"
)

// Patch represents a single key-value change to apply.
type Patch struct {
	Key   string
	Value string
}

// Result describes the outcome of applying a patch to a file.
type Result struct {
	Key    string
	Action string // "added", "updated", "skipped"
}

// Apply reads the file at path, applies the given patches, and writes the result back.
// Keys already present are updated in-place; missing keys are appended.
func Apply(path string, patches []Patch) ([]Result, error) {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("patcher: read %s: %w", path, err)
	}

	lines := []string{}
	if len(data) > 0 {
		lines = strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	}

	results := make([]Result, 0, len(patches))
	patchMap := make(map[string]string, len(patches))
	for _, p := range patches {
		patchMap[p.Key] = p.Value
	}

	applied := make(map[string]bool)
	for i, line := range lines {
		key, _, ok := parsePatchLine(line)
		if !ok {
			continue
		}
		if newVal, found := patchMap[key]; found {
			lines[i] = fmt.Sprintf("%s=%s", key, quoteIfNeeded(newVal))
			applied[key] = true
			results = append(results, Result{Key: key, Action: "updated"})
		}
	}

	for _, p := range patches {
		if !applied[p.Key] {
			lines = append(lines, fmt.Sprintf("%s=%s", p.Key, quoteIfNeeded(p.Value)))
			results = append(results, Result{Key: p.Key, Action: "added"})
		}
	}

	output := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(output), 0644); err != nil {
		return nil, fmt.Errorf("patcher: write %s: %w", path, err)
	}
	return results, nil
}

// DryRun returns what Apply would do without writing any files.
func DryRun(path string, patches []Patch) ([]Result, error) {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("patcher: read %s: %w", path, err)
	}

	existing := make(map[string]bool)
	for _, line := range strings.Split(string(data), "\n") {
		if key, _, ok := parsePatchLine(line); ok {
			existing[key] = true
		}
	}

	results := make([]Result, 0, len(patches))
	for _, p := range patches {
		action := "added"
		if existing[p.Key] {
			action = "updated"
		}
		results = append(results, Result{Key: p.Key, Action: action})
	}
	return results, nil
}

func parsePatchLine(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	line = strings.TrimPrefix(line, "export ")
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", false
	}
	return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
