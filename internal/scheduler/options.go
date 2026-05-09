package scheduler

import "time"

// Options configures Scheduler behaviour.
type Options struct {
	// AlwaysNotify fires the handler on every tick even when results are unchanged.
	AlwaysNotify bool

	// OnError is called when Compare returns an error. May be nil.
	OnError func(err error)

	// DefaultInterval is the fallback when no interval is specified by callers.
	DefaultInterval time.Duration
}

// DefaultOptions returns sensible defaults for the Scheduler.
func DefaultOptions() Options {
	return Options{
		AlwaysNotify:    false,
		DefaultInterval: 30 * time.Second,
		OnError:         nil,
	}
}
