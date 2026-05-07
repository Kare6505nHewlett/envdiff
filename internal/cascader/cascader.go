// Package cascader merges multiple .env files in priority order,
// with later files overriding keys from earlier ones.
package cascader

import (
	"fmt"
	"sort"
)

// Layer represents a named .env file layer in the cascade.
type Layer struct {
	Name   string
	Values map[string]string
}

// Result holds the resolved value for a key across all layers.
type Result struct {
	Key        string
	Value      string
	SourceLayer string
	Overridden bool // true if at least one earlier layer also defined this key
}

// Options controls cascade behaviour.
type Options struct {
	// SkipEmpty causes keys with empty values to be ignored during merge.
	SkipEmpty bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{SkipEmpty: false}
}

// Apply merges layers in order (index 0 = lowest priority) and returns
// one Result per unique key, reflecting the winning value and its source.
func Apply(layers []Layer, opts Options) ([]Result, error) {
	if len(layers) == 0 {
		return nil, fmt.Errorf("cascader: at least one layer is required")
	}

	type entry struct {
		value  string
		source string
		count  int // how many layers defined this key
	}

	resolved := make(map[string]*entry)

	for _, layer := range layers {
		for k, v := range layer.Values {
			if opts.SkipEmpty && v == "" {
				continue
			}
			if e, exists := resolved[k]; exists {
				e.value = v
				e.source = layer.Name
				e.count++
			} else {
				resolved[k] = &entry{value: v, source: layer.Name, count: 1}
			}
		}
	}

	keys := make([]string, 0, len(resolved))
	for k := range resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]Result, 0, len(keys))
	for _, k := range keys {
		e := resolved[k]
		results = append(results, Result{
			Key:         k,
			Value:       e.value,
			SourceLayer: e.source,
			Overridden:  e.count > 1,
		})
	}
	return results, nil
}
