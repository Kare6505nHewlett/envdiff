package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed .env file as a map of key-value pairs.
type EnvMap map[string]string

// ParseFile reads and parses a .env file at the given path.
// It skips blank lines and comments (lines starting with '#').
// Returns an EnvMap and any error encountered.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single env line into key and value.
// Supports optional export prefix, inline comments, and quoted values.
func parseLine(line string) (string, string, error) {
	// Strip optional "export " prefix
	line = strings.TrimPrefix(line, "export ")

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid line (missing '='): %q", line)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line: %q", line)
	}

	// Keys must not contain spaces or special characters that would be invalid
	// in a shell environment variable name.
	if strings.ContainsAny(key, " \t") {
		return "", "", fmt.Errorf("key contains whitespace in line: %q", line)
	}

	value := strings.TrimSpace(parts[1])
	value = stripInlineComment(value)
	value = unquote(value)

	return key, value, nil
}

// stripInlineComment removes trailing inline comments (unquoted # and beyond).
func stripInlineComment(value string) string {
	if !strings.HasPrefix(value, `"`) && !strings.HasPrefix(value, `'`) {
		if idx := strings.Index(value, " #"); idx != -1 {
			value = strings.TrimSpace(value[:idx])
		}
	}
	return value
}

// unquote removes surrounding single or double quotes from a value.
func unquote(value string) string {
	if len(value) >= 2 {
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
			return value[1 : len(value)-1]
		}
	}
	return value
}
