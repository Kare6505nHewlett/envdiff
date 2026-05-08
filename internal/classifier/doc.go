// Package classifier assigns semantic categories to diff results based on
// key naming conventions and value patterns.
//
// Supported categories:
//
//   - secret   — keys containing words like "password", "token", "secret"
//   - url      — values that begin with http:// or https://
//   - flag     — boolean-like values such as true/false, yes/no, on/off
//   - numeric  — integer or floating-point values
//   - unknown  — anything that does not match the above
//
// Usage:
//
//	classified := classifier.Apply(diffResults)
//	for _, r := range classified {
//		fmt.Printf("%s [%s]\n", r.Key, r.Category)
//	}
package classifier
