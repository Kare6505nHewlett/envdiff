// Package digester computes and compares cryptographic digests of .env file
// key sets, enabling quick detection of structural drift between environments.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Digest holds the computed hash and the ordered keys that produced it.
type Digest struct {
	// Hash is a SHA-256 hex string derived from the sorted key names.
	Hash string `json:"hash"`
	// Keys is the sorted list of key names that were hashed.
	Keys []string `json:"keys"`
	// File is the source label (filename or environment name).
	File string `json:"file"`
}

// Compute returns a Digest for each unique file referenced in results.
// Only key names are hashed — values are intentionally excluded so that
// the digest reflects structural shape, not secret content.
func Compute(results []diff.Result) []Digest {
	if len(results) == 0 {
		return nil
	}

	// Group keys by file.
	fileKeys := make(map[string]map[string]struct{})
	for _, r := range results {
		if _, ok := fileKeys[r.File]; !ok {
			fileKeys[r.File] = make(map[string]struct{})
		}
		fileKeys[r.File][r.Key] = struct{}{}
	}

	digests := make([]Digest, 0, len(fileKeys))
	for file, keySet := range fileKeys {
		keys := make([]string, 0, len(keySet))
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		h := sha256.New()
		for _, k := range keys {
			fmt.Fprintf(h, "%s\n", k)
		}
		hash := hex.EncodeToString(h.Sum(nil))

		digests = append(digests, Digest{
			Hash: hash,
			Keys: keys,
			File: file,
		})
	}

	sort.Slice(digests, func(i, j int) bool {
		return digests[i].File < digests[j].File
	})
	return digests
}

// Match returns true when all digests share the same hash, meaning every
// referenced file contains an identical set of key names.
func Match(digests []Digest) bool {
	if len(digests) == 0 {
		return true
	}
	base := digests[0].Hash
	for _, d := range digests[1:] {
		if d.Hash != base {
			return false
		}
	}
	return true
}

// Diff returns a human-readable summary of which files diverge from the first
// digest in the slice.
func Diff(digests []Digest) string {
	if len(digests) < 2 {
		return ""
	}
	var sb strings.Builder
	base := digests[0]
	for _, d := range digests[1:] {
		if d.Hash != base.Hash {
			fmt.Fprintf(&sb, "%s (%s) differs from %s (%s)\n",
				d.File, d.Hash[:8], base.File, base.Hash[:8])
		}
	}
	return sb.String()
}
