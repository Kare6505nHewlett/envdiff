// Package ignorer loads and applies .envdiffignore files to suppress
// specific keys or key prefixes from envdiff comparison results.
//
// An ignore file contains one rule per line:
//
//	# comment lines and blank lines are skipped
//	EXACT_KEY          // suppress a specific key
//	PREFIX_*           // suppress all keys with the given prefix
//
// Usage:
//
//	ig, err := ignorer.LoadFile(".envdiffignore")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	filtered := ig.Apply(results)
package ignorer
