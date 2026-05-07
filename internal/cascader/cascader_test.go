package cascader_test

import (
	"testing"

	"github.com/user/envdiff/internal/cascader"
)

func sampleLayers() []cascader.Layer {
	return []cascader.Layer{
		{
			Name:   "base",
			Values: map[string]string{"APP_ENV": "development", "DB_HOST": "localhost", "LOG_LEVEL": "debug"},
		},
		{
			Name:   "staging",
			Values: map[string]string{"APP_ENV": "staging", "DB_HOST": "staging-db"},
		},
	}
}

func TestApply_MergesLayers(t *testing.T) {
	results, err := cascader.Apply(sampleLayers(), cascader.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestApply_LaterLayerWins(t *testing.T) {
	results, err := cascader.Apply(sampleLayers(), cascader.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "APP_ENV" && r.Value != "staging" {
			t.Errorf("expected staging to win, got %q", r.Value)
		}
	}
}

func TestApply_OverriddenFlag(t *testing.T) {
	results, _ := cascader.Apply(sampleLayers(), cascader.DefaultOptions())
	for _, r := range results {
		switch r.Key {
		case "APP_ENV", "DB_HOST":
			if !r.Overridden {
				t.Errorf("key %q should be marked overridden", r.Key)
			}
		case "LOG_LEVEL":
			if r.Overridden {
				t.Errorf("key %q should not be marked overridden", r.Key)
			}
		}
	}
}

func TestApply_SortedByKey(t *testing.T) {
	results, _ := cascader.Apply(sampleLayers(), cascader.DefaultOptions())
	for i := 1; i < len(results); i++ {
		if results[i].Key < results[i-1].Key {
			t.Errorf("results not sorted: %q before %q", results[i-1].Key, results[i].Key)
		}
	}
}

func TestApply_SkipEmpty(t *testing.T) {
	layers := []cascader.Layer{
		{Name: "base", Values: map[string]string{"KEY": "value", "EMPTY": ""}},
	}
	opts := cascader.DefaultOptions()
	opts.SkipEmpty = true
	results, _ := cascader.Apply(layers, opts)
	for _, r := range results {
		if r.Key == "EMPTY" {
			t.Error("empty key should have been skipped")
		}
	}
}

func TestApply_EmptyLayersReturnsError(t *testing.T) {
	_, err := cascader.Apply(nil, cascader.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty layers")
	}
}

func TestApply_SourceLayerTracked(t *testing.T) {
	results, _ := cascader.Apply(sampleLayers(), cascader.DefaultOptions())
	for _, r := range results {
		if r.Key == "LOG_LEVEL" && r.SourceLayer != "base" {
			t.Errorf("expected source=base, got %q", r.SourceLayer)
		}
		if r.Key == "APP_ENV" && r.SourceLayer != "staging" {
			t.Errorf("expected source=staging, got %q", r.SourceLayer)
		}
	}
}
