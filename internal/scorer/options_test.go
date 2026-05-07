package scorer

import "testing"

func TestDefaultOptions_MaxScore(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxScore != 100.0 {
		t.Errorf("MaxScore = %.1f, want 100.0", opts.MaxScore)
	}
}

func TestOptions_CanOverrideMissingWeight(t *testing.T) {
	opts := DefaultOptions()
	opts.MissingWeight = 2.0
	if opts.MissingWeight != 2.0 {
		t.Errorf("MissingWeight override failed, got %.1f", opts.MissingWeight)
	}
}

func TestOptions_CanOverrideMismatchWeight(t *testing.T) {
	opts := DefaultOptions()
	opts.MismatchWeight = 0.0
	if opts.MismatchWeight != 0.0 {
		t.Errorf("MismatchWeight override failed, got %.1f", opts.MismatchWeight)
	}
}

func TestOptions_CanOverrideMaxScore(t *testing.T) {
	opts := DefaultOptions()
	opts.MaxScore = 50.0
	if opts.MaxScore != 50.0 {
		t.Errorf("MaxScore override failed, got %.1f", opts.MaxScore)
	}
}

func TestOptions_ZeroWeightsAllowed(t *testing.T) {
	opts := Options{
		MissingWeight:  0,
		MismatchWeight: 0,
		MaxScore:       100,
	}
	if opts.MissingWeight != 0 || opts.MismatchWeight != 0 {
		t.Error("expected zero weights to be accepted")
	}
}
