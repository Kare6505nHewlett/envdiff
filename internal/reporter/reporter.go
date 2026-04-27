package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes the diff results to w in the specified format.
func Report(w io.Writer, results []diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return reportJSON(w, results)
	case FormatText:
		return reportText(w, results)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func reportText(w io.Writer, results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "✓ All keys match across environments.")
		return err
	}

	for _, r := range results {
		var line string
		switch r.Status {
		case diff.StatusMissing:
			line = fmt.Sprintf("[MISSING]  %-30s missing in: %s", r.Key, strings.Join(r.Files, ", "))
		case diff.StatusExtra:
			line = fmt.Sprintf("[EXTRA]    %-30s only in:   %s", r.Key, strings.Join(r.Files, ", "))
		case diff.StatusMismatch:
			line = fmt.Sprintf("[MISMATCH] %-30s values differ across files", r.Key)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func reportJSON(w io.Writer, results []diff.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
