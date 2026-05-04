package templater

import "fmt"

// DefaultPlaceholder is used when Options.Placeholder is empty and
// explicit empty-string behaviour is not desired.
const DefaultPlaceholder = ""

// Validate checks that the Options struct contains sensible values.
// It returns a descriptive error for any invalid combination.
func (o Options) Validate() error {
	for _, ch := range o.Placeholder {
		if ch == '\n' || ch == '\r' {
			return fmt.Errorf("templater: placeholder must not contain newline characters")
		}
	}
	return nil
}

// WithDefaults returns a copy of Options with sensible defaults applied.
func (o Options) WithDefaults() Options {
	if o.Placeholder == "" {
		o.Placeholder = DefaultPlaceholder
	}
	return o
}
