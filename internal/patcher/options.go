package patcher

// Options controls the behaviour of the patcher.
type Options struct {
	// DryRun, when true, prevents any files from being written.
	DryRun bool

	// SkipExisting, when true, only adds missing keys and never updates
	// keys that are already present in the file.
	SkipExisting bool

	// Backup, when true, writes a copy of the original file to <path>.bak
	// before applying patches.
	Backup bool
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{
		DryRun:       false,
		SkipExisting: false,
		Backup:       false,
	}
}
