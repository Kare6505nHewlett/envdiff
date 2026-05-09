// Package scheduler provides periodic re-diffing of .env files on a configurable interval.
// It runs Compare in the background and invokes a callback whenever results change.
package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Handler is called with the latest diff results whenever a change is detected.
type Handler func(results []diff.Result)

// Scheduler periodically compares env files and fires a Handler on change.
type Scheduler struct {
	files    []string
	interval time.Duration
	handler  Handler
	opts     Options

	mu      sync.Mutex
	lastKey string
}

// New creates a Scheduler for the given files, interval, and change handler.
func New(files []string, interval time.Duration, handler Handler, opts Options) *Scheduler {
	return &Scheduler{
		files:    files,
		interval: interval,
		handler:  handler,
		opts:     opts,
	}
}

// Run starts the scheduler loop and blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run once immediately.
	s.tick()

	for {
		select {
		case <-ticker.C:
			s.tick()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) tick() {
	results, err := diff.Compare(s.files)
	if err != nil {
		if s.opts.OnError != nil {
			s.opts.OnError(err)
		}
		return
	}

	key := resultKey(results)

	s.mu.Lock()
	changed := key != s.lastKey
	if changed {
		s.lastKey = key
	}
	s.mu.Unlock()

	if changed || s.opts.AlwaysNotify {
		s.handler(results)
	}
}

// resultKey produces a cheap fingerprint of results to detect changes.
func resultKey(results []diff.Result) string {
	h := make([]byte, 0, len(results)*32)
	for _, r := range results {
		h = append(h, r.Key...)
		h = append(h, '|')
		h = append(h, r.Status...)
		h = append(h, ';')
	}
	return string(h)
}
