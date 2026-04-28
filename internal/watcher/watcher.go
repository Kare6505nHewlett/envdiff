// Package watcher monitors .env files for changes and triggers a callback
// when modifications are detected, enabling live diff updates.
package watcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// ChangeFunc is called when one or more watched files change.
type ChangeFunc func(changedFiles []string)

// Watcher polls a set of files at a given interval and fires a callback on change.
type Watcher struct {
	files    []string
	interval time.Duration
	hashes   map[string]string
	stop     chan struct{}
}

// New creates a Watcher for the given files, polling at the specified interval.
func New(files []string, interval time.Duration) *Watcher {
	return &Watcher{
		files:    files,
		interval: interval,
		hashes:   make(map[string]string),
		stop:     make(chan struct{}),
	}
}

// Start begins polling and calls onChange whenever files change.
// It blocks until Stop is called.
func (w *Watcher) Start(onChange ChangeFunc) {
	// Capture initial hashes so first tick only fires on real changes.
	for _, f := range w.files {
		if h, err := hashFile(f); err == nil {
			w.hashes[f] = h
		}
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			if changed := w.detectChanges(); len(changed) > 0 {
				onChange(changed)
			}
		}
	}
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stop)
}

// detectChanges returns files whose content hash has changed since last check.
func (w *Watcher) detectChanges() []string {
	var changed []string
	for _, f := range w.files {
		h, err := hashFile(f)
		if err != nil {
			continue
		}
		if prev, ok := w.hashes[f]; !ok || prev != h {
			w.hashes[f] = h
			changed = append(changed, f)
		}
	}
	return changed
}

// hashFile returns the MD5 hex digest of the file at path.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
