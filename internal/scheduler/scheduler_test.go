package scheduler_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scheduler"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestScheduler_FiresHandlerOnChange(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.env")
	b := filepath.Join(dir, "b.env")
	os.WriteFile(a, []byte("KEY=value\n"), 0644)
	os.WriteFile(b, []byte("KEY=other\n"), 0644)

	var calls atomic.Int32
	h := func(_ []diff.Result) { calls.Add(1) }

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	s := scheduler.New([]string{a, b}, 50*time.Millisecond, h, scheduler.DefaultOptions())
	s.Run(ctx)

	if calls.Load() == 0 {
		t.Fatal("expected handler to be called at least once")
	}
}

func TestScheduler_NoExtraCallsWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.env")
	b := filepath.Join(dir, "b.env")
	os.WriteFile(a, []byte("KEY=same\n"), 0644)
	os.WriteFile(b, []byte("KEY=same\n"), 0644)

	var calls atomic.Int32
	h := func(_ []diff.Result) { calls.Add(1) }

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	opts := scheduler.DefaultOptions()
	opts.AlwaysNotify = false
	s := scheduler.New([]string{a, b}, 50*time.Millisecond, h, opts)
	s.Run(ctx)

	// Should fire exactly once (initial tick), then no more because nothing changed.
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected 1 call, got %d", got)
	}
}

func TestScheduler_AlwaysNotify(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.env")
	b := filepath.Join(dir, "b.env")
	os.WriteFile(a, []byte("KEY=same\n"), 0644)
	os.WriteFile(b, []byte("KEY=same\n"), 0644)

	var calls atomic.Int32
	h := func(_ []diff.Result) { calls.Add(1) }

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	opts := scheduler.DefaultOptions()
	opts.AlwaysNotify = true
	s := scheduler.New([]string{a, b}, 50*time.Millisecond, h, opts)
	s.Run(ctx)

	if calls.Load() < 2 {
		t.Fatalf("expected multiple calls with AlwaysNotify, got %d", calls.Load())
	}
}

func TestDefaultOptions_Defaults(t *testing.T) {
	opts := scheduler.DefaultOptions()
	if opts.AlwaysNotify {
		t.Error("expected AlwaysNotify to be false")
	}
	if opts.DefaultInterval != 30*time.Second {
		t.Errorf("unexpected DefaultInterval: %v", opts.DefaultInterval)
	}
	if opts.OnError != nil {
		t.Error("expected OnError to be nil")
	}
}
