package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestValidate_NoViolations(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	rule := validator.Rule{}
	got := validator.Validate("prod.env", env, rule)
	if len(got) != 0 {
		t.Errorf("expected no violations, got %d", len(got))
	}
}

func TestValidate_ForbiddenKey(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "abc123",
		"APP_HOST":   "localhost",
	}
	rule := validator.Rule{ForbiddenKeys: []string{"SECRET_KEY"}}
	got := validator.Validate("dev.env", env, rule)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "SECRET_KEY" {
		t.Errorf("expected violation for SECRET_KEY, got %q", got[0].Key)
	}
}

func TestValidate_RequiredPrefix(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"DATABASE": "postgres",
	}
	rule := validator.Rule{RequiredPrefix: "APP_"}
	got := validator.Validate("staging.env", env, rule)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "DATABASE" {
		t.Errorf("expected violation for DATABASE, got %q", got[0].Key)
	}
}

func TestValidate_NoEmptyValues(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "",
		"APP_PORT": "8080",
	}
	rule := validator.Rule{NoEmptyValues: true}
	got := validator.Validate("prod.env", env, rule)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "APP_HOST" {
		t.Errorf("expected violation for APP_HOST, got %q", got[0].Key)
	}
}

func TestValidate_MultipleRules(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "",
		"SECRET":   "value",
	}
	rule := validator.Rule{
		NoEmptyValues:  true,
		ForbiddenKeys:  []string{"SECRET"},
		RequiredPrefix: "APP_",
	}
	got := validator.Validate("dev.env", env, rule)
	// APP_HOST: empty value + no prefix violation? prefix matches, empty value = 1
	// SECRET: forbidden + no prefix = 2
	// Total expected: 3
	if len(got) != 3 {
		t.Errorf("expected 3 violations, got %d", len(got))
	}
}

func TestValidate_FileNamePropagated(t *testing.T) {
	env := map[string]string{"SECRET": "x"}
	rule := validator.Rule{ForbiddenKeys: []string{"SECRET"}}
	got := validator.Validate("prod.env", env, rule)
	if len(got) == 0 {
		t.Fatal("expected at least one violation")
	}
	if got[0].File != "prod.env" {
		t.Errorf("expected file prod.env, got %q", got[0].File)
	}
}
