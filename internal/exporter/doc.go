// Package exporter provides functionality to export diff results to various
// structured file formats including CSV, Markdown, and JSON.
//
// Usage:
//
//	results := diff.Compare(base, targets)
//	f, _ := os.Create("report.csv")
//	defer f.Close()
//	exporter.Export(f, results, exporter.FormatCSV)
//
// Supported formats:
//   - csv      — comma-separated values, suitable for spreadsheets
//   - markdown — GitHub-flavoured markdown table
//   - json     — pretty-printed JSON array
package exporter
