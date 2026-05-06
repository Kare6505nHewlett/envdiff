package resolver

import "testing"

func TestDefaultOptions_MaxDepth(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxDepth != 10 {
		t.Errorf("expected MaxDepth 10, got %d", opts.MaxDepth)
	}
}

func TestDefaultOptions_FailOnMissingIsFalse(t *testing.T) {
	opts := DefaultOptions()
	if opts.FailOnMissing {
		t.Error("expected FailOnMissing to be false by default")
	}
}

func TestOptions_CanOverrideMaxDepth(t *testing.T) {
	opts := DefaultOptions()
	opts.MaxDepth = 5
	if opts.MaxDepth != 5 {
		t.Errorf("expected MaxDepth 5, got %d", opts.MaxDepth)
	}
}

func TestOptions_CanEnableFailOnMissing(t *testing.T) {
	opts := DefaultOptions()
	opts.FailOnMissing = true
	if !opts.FailOnMissing {
		t.Error("expected FailOnMissing to be true after override")
	}
}

func TestResolve_EmptyEnv(t *testing.T) {
	got, err := Resolve(map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestResolve_PreservesAllKeys(t *testing.T) {
	env := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}
	got, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(got))
	}
}
