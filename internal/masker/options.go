package masker

// Options controls how values are masked.
type Options struct {
	// PrefixLen is the number of characters revealed at the start of a value.
	PrefixLen int

	// SuffixLen is the number of characters revealed at the end of a value.
	SuffixLen int

	// MaskChar is the character used to fill the hidden portion.
	MaskChar rune

	// MinMaskLen is the minimum number of mask characters inserted, even when
	// the value is shorter than PrefixLen+SuffixLen.
	MinMaskLen int

	// SensitiveKeys is a list of substrings; any key whose upper-cased name
	// contains one of these substrings will have its value masked.
	SensitiveKeys []string
}
