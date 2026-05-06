// Package resolver expands variable references within .env file values.
//
// It supports both ${VAR} and $VAR reference styles and handles
// nested expansions up to a configurable depth limit.
//
// Example usage:
//
//	env := map[string]string{
//		"HOST": "localhost",
//		"DSN":  "postgres://${HOST}/mydb",
//	}
//
//	resolved, err := resolver.Resolve(env, resolver.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	// resolved["DSN"] == "postgres://localhost/mydb"
package resolver
