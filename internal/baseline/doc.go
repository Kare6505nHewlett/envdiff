// Package baseline provides save and load functionality for persisting a
// parsed .env file as a JSON snapshot (baseline). This baseline can later
// be loaded and used as the reference side when comparing against other
// environment files, enabling drift detection over time.
//
// Usage:
//
//	// Save a baseline from parsed keys
//	err := baseline.Save("baseline.json", ".env.production", keys)
//
//	// Load it back for comparison
//	b, err := baseline.Load("baseline.json")
//	fmt.Println(b.Keys)
package baseline
