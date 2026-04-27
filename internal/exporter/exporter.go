package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the export format type.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
	FormatJSON     Format = "json"
)

// Export writes diff results to w in the given format.
func Export(w io.Writer, results []diff.Result, format Format) error {
	switch format {
	case FormatCSV:
		return exportCSV(w, results)
	case FormatMarkdown:
		return exportMarkdown(w, results)
	case FormatJSON:
		return exportJSON(w, results)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportCSV(w io.Writer, results []diff.Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "status", "file", "base_value", "compare_value"}); err != nil {
		return err
	}
	for _, r := range results {
		row := []string{r.Key, string(r.Status), r.File, r.BaseValue, r.CompareValue}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportMarkdown(w io.Writer, results []diff.Result) error {
	fmt.Fprintln(w, "| Key | Status | File | Base Value | Compare Value |")
	fmt.Fprintln(w, "|-----|--------|------|------------|---------------|")
	for _, r := range results {
		line := fmt.Sprintf("| %s | %s | %s | %s | %s |",
			r.Key, r.Status, r.File,
			escapeMD(r.BaseValue), escapeMD(r.CompareValue))
		fmt.Fprintln(w, line)
	}
	return nil
}

func exportJSON(w io.Writer, results []diff.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func escapeMD(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
