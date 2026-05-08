package scoper

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable summary of scoped results to w.
func ReportText(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No scoped results.")
		return
	}
	for _, r := range results {
		fmt.Fprintf(w, "[%s] %d entries\n", strings.ToUpper(r.Scope), len(r.Entries))
		for _, e := range r.Entries {
			fmt.Fprintf(w, "  %-30s %-12s %s\n", e.Key, e.Status, e.File)
		}
	}
}

// jsonScope is the serialisable form of a Result.
type jsonScope struct {
	Scope   string        `json:"scope"`
	Count   int           `json:"count"`
	Entries []jsonEntry   `json:"entries"`
}

type jsonEntry struct {
	Key    string `json:"key"`
	File   string `json:"file"`
	Status string `json:"status"`
}

// ReportJSON writes a JSON array of scoped results to w.
func ReportJSON(w io.Writer, results []Result) error {
	var out []jsonScope
	for _, r := range results {
		js := jsonScope{Scope: r.Scope, Count: len(r.Entries)}
		for _, e := range r.Entries {
			js.Entries = append(js.Entries, jsonEntry{
				Key:    e.Key,
				File:   e.File,
				Status: string(e.Status),
			})
		}
		out = append(out, js)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
