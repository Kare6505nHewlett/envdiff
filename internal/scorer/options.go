package scorer

// Options controls how the health score is computed.
type Options struct {
	// MissingWeight is the penalty applied per missing key (default: 1.0).
	MissingWeight float64

	// MismatchWeight is the penalty applied per mismatched key (default: 0.5).
	MismatchWeight float64

	// MaxScore is the ceiling for the returned score (default: 100.0).
	MaxScore float64
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		MissingWeight:  1.0,
		MismatchWeight: 0.5,
		MaxScore:       100.0,
	}
}
