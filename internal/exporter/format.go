package exporter

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format, returning an error if unrecognized.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(strings.TrimSpace(s))) {
	case FormatCSV:
		return FormatCSV, nil
	case FormatMarkdown:
		return FormatMarkdown, nil
	case FormatJSON:
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown format %q: supported formats are csv, markdown, json", s)
	}
}

// SupportedFormats returns all supported export format names.
func SupportedFormats() []string {
	return []string{
		string(FormatCSV),
		string(FormatMarkdown),
		string(FormatJSON),
	}
}
