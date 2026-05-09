package patcher_test

import (
	"testing"

	"github.com/user/envdiff/internal/patcher"
)

func TestDefaultOptions_Defaults(t *testing.T) {
	opts := patcher.DefaultOptions()
	if opts.DryRun {
		t.Error("expected DryRun to be false by default")
	}
	if opts.SkipExisting {
		t.Error("expected SkipExisting to be false by default")
	}
	if opts.Backup {
		t.Error("expected Backup to be false by default")
	}
}

func TestOptions_CanOverride(t *testing.T) {
	opts := patcher.DefaultOptions()
	opts.DryRun = true
	opts.SkipExisting = true
	opts.Backup = true

	if !opts.DryRun {
		t.Error("expected DryRun to be true after override")
	}
	if !opts.SkipExisting {
		t.Error("expected SkipExisting to be true after override")
	}
	if !opts.Backup {
		t.Error("expected Backup to be true after override")
	}
}

func TestDefaultOptions_IndependentInstances(t *testing.T) {
	opts1 := patcher.DefaultOptions()
	opts2 := patcher.DefaultOptions()

	opts1.DryRun = true
	opts1.Backup = true

	if opts2.DryRun {
		t.Error("expected modifying opts1 to not affect opts2 DryRun")
	}
	if opts2.Backup {
		t.Error("expected modifying opts1 to not affect opts2 Backup")
	}
}
