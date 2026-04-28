package watcher_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "KEY=original\n")

	var mu sync.Mutex
	var got []string

	w := watcher.New([]string{path}, 20*time.Millisecond)
	go w.Start(func(changed []string) {
		mu.Lock()
		got = append(got, changed...)
		mu.Unlock()
	})
	defer w.Stop()

	// Allow watcher to capture initial hash.
	time.Sleep(30 * time.Millisecond)

	// Modify the file.
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Wait for at least one poll cycle.
	time.Sleep(60 * time.Millisecond)
	w.Stop()

	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected change callback to fire, got none")
	}
	if got[0] != path {
		t.Errorf("expected %q, got %q", path, got[0])
	}
}

func TestWatcher_NoCallbackOnUnchangedFile(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, ".env", "KEY=stable\n")

	var mu sync.Mutex
	var called int

	w := watcher.New([]string{path}, 20*time.Millisecond)
	go w.Start(func(_ []string) {
		mu.Lock()
		called++
		mu.Unlock()
	})
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)
	w.Stop()

	mu.Lock()
	defer mu.Unlock()
	if called != 0 {
		t.Errorf("expected 0 callbacks, got %d", called)
	}
}

func TestWatcher_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempEnv(t, dir, ".env.dev", "A=1\n")
	p2 := writeTempEnv(t, dir, ".env.prod", "B=2\n")

	var mu sync.Mutex
	changedSet := map[string]bool{}

	w := watcher.New([]string{p1, p2}, 20*time.Millisecond)
	go w.Start(func(changed []string) {
		mu.Lock()
		for _, c := range changed {
			changedSet[c] = true
		}
		mu.Unlock()
	})
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	os.WriteFile(p1, []byte("A=updated\n"), 0644)
	os.WriteFile(p2, []byte("B=updated\n"), 0644)
	time.Sleep(60 * time.Millisecond)
	w.Stop()

	mu.Lock()
	defer mu.Unlock()
	if !changedSet[p1] || !changedSet[p2] {
		t.Errorf("expected both files detected; got %v", changedSet)
	}
}
