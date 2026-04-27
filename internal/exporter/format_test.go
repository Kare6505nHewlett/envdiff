package exporter_test

import (
	"testing"

	"github.com/user/envdiff/internal/exporter"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected exporter.Format
	}{
		{"csv", exporter.FormatCSV},
		{"CSV", exporter.FormatCSV},
		{"markdown", exporter.FormatMarkdown},
		{"Markdown", exporter.FormatMarkdown},
		{"json", exporter.FormatJSON},
		{"JSON", exporter.FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := exporter.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := exporter.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := exporter.SupportedFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 supported formats, got %d", len(formats))
	}
	set := make(map[string]bool)
	for _, f := range formats {
		set[f] = true
	}
	for _, expected := range []string{"csv", "markdown", "json"} {
		if !set[expected] {
			t.Errorf("expected format %q in supported list", expected)
		}
	}
}
