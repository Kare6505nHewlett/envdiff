// Package baseline provides functionality to save and load a reference
// .env file as a baseline for future comparisons.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of parsed key-value pairs.
type Baseline struct {
	CreatedAt time.Time         `json:"created_at"`
	SourceFile string           `json:"source_file"`
	Keys      map[string]string `json:"keys"`
}

// Save writes the provided key-value map and source file path to a JSON
// baseline file at the given destination path.
func Save(destPath, sourceFile string, keys map[string]string) error {
	b := Baseline{
		CreatedAt:  time.Now().UTC(),
		SourceFile: sourceFile,
		Keys:       keys,
	}

	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal failed: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write failed: %w", err)
	}

	return nil
}

// Load reads a baseline JSON file from the given path and returns the
// parsed Baseline struct.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read failed: %w", err)
	}

	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal failed: %w", err)
	}

	if b.Keys == nil {
		b.Keys = make(map[string]string)
	}

	return &b, nil
}
