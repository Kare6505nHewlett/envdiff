package promoter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReportText writes a human-readable promotion plan to w.
func ReportText(w io.Writer, plan Plan) {
	fmt.Fprintf(w, "Promotion plan: %s → %s\n", plan.SourceFile, plan.TargetFile)
	if len(plan.Steps) == 0 {
		fmt.Fprintln(w, "  No changes required.")
		return
	}
	fmt.Fprintf(w, "  %d step(s):\n", len(plan.Steps))
	for _, s := range plan.Steps {
		fmt.Fprintf(w, "  %s\n", Describe(s))
	}
}

// jsonStep is the JSON-serialisable representation of a Step.
type jsonStep struct {
	Key      string `json:"key"`
	Action   string `json:"action"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value"`
	File     string `json:"file"`
}

// jsonPlan is the JSON-serialisable representation of a Plan.
type jsonPlan struct {
	SourceFile string     `json:"source_file"`
	TargetFile string     `json:"target_file"`
	StepCount  int        `json:"step_count"`
	Steps      []jsonStep `json:"steps"`
}

// ReportJSON writes the promotion plan as a JSON object to w.
func ReportJSON(w io.Writer, plan Plan) error {
	steps := make([]jsonStep, 0, len(plan.Steps))
	for _, s := range plan.Steps {
		steps = append(steps, jsonStep{
			Key:      s.Key,
			Action:   string(s.Action),
			OldValue: s.OldValue,
			NewValue: s.NewValue,
			File:     s.File,
		})
	}
	out := jsonPlan{
		SourceFile: plan.SourceFile,
		TargetFile: plan.TargetFile,
		StepCount:  len(steps),
		Steps:      steps,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// StepSummary returns a compact one-line summary of the plan.
func StepSummary(plan Plan) string {
	if len(plan.Steps) == 0 {
		return "nothing to promote"
	}
	var adds, updates int
	for _, s := range plan.Steps {
		switch s.Action {
		case ActionAdd:
			adds++
		case ActionUpdate:
			updates++
		}
	}
	parts := []string{}
	if adds > 0 {
		parts = append(parts, fmt.Sprintf("%d to add", adds))
	}
	if updates > 0 {
		parts = append(parts, fmt.Sprintf("%d to update", updates))
	}
	return strings.Join(parts, ", ")
}
