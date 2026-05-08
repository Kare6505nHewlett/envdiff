// Package scoper partitions diff results into named environment scopes.
//
// A Scope associates a human-readable name (e.g. "production", "staging")
// with one or more .env file paths. Apply distributes a flat []diff.Result
// slice into per-scope buckets; any result whose file does not match a
// declared scope is placed in the reserved "default" scope.
//
// Typical usage:
//
//	scopes := []scoper.Scope{
//	    {Name: "production", Files: []string{"prod.env", "prod.secrets.env"}},
//	    {Name: "staging",    Files: []string{"staging.env"}},
//	}
//	partitioned := scoper.Apply(results, scopes)
//	scoper.ReportText(os.Stdout, partitioned)
package scoper
