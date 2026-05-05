// Package annotator enriches diff results with human-readable descriptions.
//
// Descriptions are loaded from a plain-text file in KEY=description format
// (typically named .env.descriptions and co-located with the .env files).
// The Apply function merges loaded descriptions into diff.Result slices,
// producing AnnotatedResult values that downstream reporters and exporters
// can use to surface context alongside each finding.
//
// Example usage:
//
//	desc, err := annotator.LoadDescriptions(".env.descriptions")
//	if err != nil { ... }
//	annotated := annotator.Apply(results, desc)
package annotator
