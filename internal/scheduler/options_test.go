package scheduler_test

import (
	"testing"
	"time"

	"github.com/user/envdiff/internal/scheduler"
)

func TestOptions_CanOverrideAlwaysNotify(t *testing.T) {
	opts := scheduler.DefaultOptions()
	opts.AlwaysNotify = true
	if !opts.AlwaysNotify {
		t.Error("expected AlwaysNotify to be true after override")
	}
}

func TestOptions_CanOverrideDefaultInterval(t *testing.T) {
	opts := scheduler.DefaultOptions()
	opts.DefaultInterval = 5 * time.Second
	if opts.DefaultInterval != 5*time.Second {
		t.Errorf("unexpected DefaultInterval: %v", opts.DefaultInterval)
	}
}

func TestOptions_CanSetOnError(t *testing.T) {
	var called bool
	opts := scheduler.DefaultOptions()
	opts.OnError = func(err error) { called = true }
	if opts.OnError == nil {
		t.Fatal("expected OnError to be set")
	}
	opts.OnError(nil)
	if !called {
		t.Error("expected OnError callback to be invoked")
	}
}

func TestOptions_ZeroIntervalIsAllowed(t *testing.T) {
	opts := scheduler.DefaultOptions()
	opts.DefaultInterval = 0
	if opts.DefaultInterval != 0 {
		t.Errorf("expected zero interval, got %v", opts.DefaultInterval)
	}
}
