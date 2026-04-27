package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortField defines the field to sort results by.
type SortField string

const (
	SortByKey    SortField = "key"
	SortByStatus SortField = "status"
	SortByFile   SortField = "file"
)

// Order defines the sort direction.
type Order string

const (
	Ascending  Order = "asc"
	Descending Order = "desc"
)

// Options holds sorting configuration.
type Options struct {
	Field SortField
	Order Order
}

// Sort returns a sorted copy of the diff results based on the given options.
func Sort(results []diff.Result, opts Options) []diff.Result {
	if len(results) == 0 {
		return results
	}

	copy_ := make([]diff.Result, len(results))
	copy(copy_, results)

	sort.SliceStable(copy_, func(i, j int) bool {
		var less bool
		switch opts.Field {
		case SortByStatus:
			less = statusRank(copy_[i].Status) < statusRank(copy_[j].Status)
		case SortByFile:
			less = copy_[i].File < copy_[j].File
		default: // SortByKey
			less = copy_[i].Key < copy_[j].Key
		}
		if opts.Order == Descending {
			return !less
		}
		return less
	})

	return copy_
}

// statusRank assigns a numeric rank to each status for ordering.
func statusRank(status diff.Status) int {
	switch status {
	case diff.StatusMissing:
		return 0
	case diff.StatusMismatch:
		return 1
	case diff.StatusMatch:
		return 2
	default:
		return 3
	}
}
