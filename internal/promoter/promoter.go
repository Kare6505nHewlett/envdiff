// Package promoter applies diff results from one environment pair to generate
// a promotion plan — a set of key/value patches needed to bring a target
// environment up to parity with a source environment.
package promoter

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Action describes what should happen to a key during promotion.
type Action string

const (
	ActionAdd    Action = "add"
	ActionUpdate Action = "update"
	ActionNone   Action = "none"
)

// Step represents a single promotion step for one key.
type Step struct {
	Key      string
	Action   Action
	OldValue string // empty when Action is ActionAdd
	NewValue string
	File     string // target file to patch
}

// Plan is the full set of steps required to promote changes.
type Plan struct {
	Steps      []Step
	SourceFile string
	TargetFile string
}

// Build computes a promotion Plan from a slice of diff.Result entries.
// sourceFile is treated as the authoritative environment; targetFile is the
// environment being brought up to parity.
func Build(results []diff.Result, sourceFile, targetFile string) Plan {
	plan := Plan{
		SourceFile: sourceFile,
		TargetFile: targetFile,
	}

	for _, r := range results {
		switch r.Status {
		case diff.StatusMissing:
			// Key exists in source but not in target — add it.
			if r.File == sourceFile {
				plan.Steps = append(plan.Steps, Step{
					Key:      r.Key,
					Action:   ActionAdd,
					NewValue: r.Value,
					File:     targetFile,
				})
			}
		case diff.StatusMismatch:
			// Value differs; use the source value for the target.
			if r.File == targetFile {
				plan.Steps = append(plan.Steps, Step{
					Key:      r.Key,
					Action:   ActionUpdate,
					OldValue: r.Value,
					NewValue: sourceValueFor(results, r.Key, sourceFile),
					File:     targetFile,
				})
			}
		}
	}

	sort.Slice(plan.Steps, func(i, j int) bool {
		return plan.Steps[i].Key < plan.Steps[j].Key
	})

	return plan
}

// Describe returns a human-readable summary line for a Step.
func Describe(s Step) string {
	switch s.Action {
	case ActionAdd:
		return fmt.Sprintf("[add]    %s = %q → %s", s.Key, s.NewValue, s.File)
	case ActionUpdate:
		return fmt.Sprintf("[update] %s: %q → %q in %s", s.Key, s.OldValue, s.NewValue, s.File)
	default:
		return fmt.Sprintf("[none]   %s", s.Key)
	}
}

func sourceValueFor(results []diff.Result, key, sourceFile string) string {
	for _, r := range results {
		if r.Key == key && r.File == sourceFile {
			return r.Value
		}
	}
	return ""
}
