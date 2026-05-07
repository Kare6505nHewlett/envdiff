package promoter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/promoter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Value: "production", File: "prod.env", Status: diff.StatusMatch},
		{Key: "APP_ENV", Value: "staging", File: "staging.env", Status: diff.StatusMismatch},
		{Key: "SECRET_KEY", Value: "abc123", File: "prod.env", Status: diff.StatusMissing},
		{Key: "DEBUG", Value: "false", File: "prod.env", Status: diff.StatusMatch},
		{Key: "DEBUG", Value: "false", File: "staging.env", Status: diff.StatusMatch},
	}
}

func TestBuild_AddsStepForMissingKey(t *testing.T) {
	results := sampleResults()
	plan := promoter.Build(results, "prod.env", "staging.env")

	var addSteps []promoter.Step
	for _, s := range plan.Steps {
		if s.Action == promoter.ActionAdd {
			addSteps = append(addSteps, s)
		}
	}

	if len(addSteps) != 1 {
		t.Fatalf("expected 1 add step, got %d", len(addSteps))
	}
	if addSteps[0].Key != "SECRET_KEY" {
		t.Errorf("expected key SECRET_KEY, got %s", addSteps[0].Key)
	}
	if addSteps[0].NewValue != "abc123" {
		t.Errorf("expected value abc123, got %s", addSteps[0].NewValue)
	}
	if addSteps[0].File != "staging.env" {
		t.Errorf("expected file staging.env, got %s", addSteps[0].File)
	}
}

func TestBuild_AddsStepForMismatchedKey(t *testing.T) {
	results := sampleResults()
	plan := promoter.Build(results, "prod.env", "staging.env")

	var updateSteps []promoter.Step
	for _, s := range plan.Steps {
		if s.Action == promoter.ActionUpdate {
			updateSteps = append(updateSteps, s)
		}
	}

	if len(updateSteps) != 1 {
		t.Fatalf("expected 1 update step, got %d", len(updateSteps))
	}
	if updateSteps[0].Key != "APP_ENV" {
		t.Errorf("expected key APP_ENV, got %s", updateSteps[0].Key)
	}
	if updateSteps[0].OldValue != "staging" {
		t.Errorf("expected old value staging, got %s", updateSteps[0].OldValue)
	}
	if updateSteps[0].NewValue != "production" {
		t.Errorf("expected new value production, got %s", updateSteps[0].NewValue)
	}
}

func TestBuild_StepsAreSortedByKey(t *testing.T) {
	results := sampleResults()
	plan := promoter.Build(results, "prod.env", "staging.env")

	for i := 1; i < len(plan.Steps); i++ {
		if plan.Steps[i].Key < plan.Steps[i-1].Key {
			t.Errorf("steps not sorted: %s before %s", plan.Steps[i-1].Key, plan.Steps[i].Key)
		}
	}
}

func TestBuild_NoStepsForMatchingKeys(t *testing.T) {
	results := []diff.Result{
		{Key: "DEBUG", Value: "false", File: "prod.env", Status: diff.StatusMatch},
		{Key: "DEBUG", Value: "false", File: "staging.env", Status: diff.StatusMatch},
	}
	plan := promoter.Build(results, "prod.env", "staging.env")
	if len(plan.Steps) != 0 {
		t.Errorf("expected 0 steps, got %d", len(plan.Steps))
	}
}

func TestDescribe_AddAction(t *testing.T) {
	s := promoter.Step{Key: "FOO", Action: promoter.ActionAdd, NewValue: "bar", File: "staging.env"}
	out := promoter.Describe(s)
	if !strings.Contains(out, "[add]") || !strings.Contains(out, "FOO") {
		t.Errorf("unexpected describe output: %s", out)
	}
}

func TestDescribe_UpdateAction(t *testing.T) {
	s := promoter.Step{Key: "APP_ENV", Action: promoter.ActionUpdate, OldValue: "staging", NewValue: "production", File: "staging.env"}
	out := promoter.Describe(s)
	if !strings.Contains(out, "[update]") || !strings.Contains(out, "APP_ENV") {
		t.Errorf("unexpected describe output: %s", out)
	}
}
