// Package watcher provides file-change detection for .env files used by envdiff.
//
// It polls a list of file paths at a configurable interval, computing an MD5
// hash of each file's contents on every tick. When a hash differs from the
// previously recorded value the file is considered changed, and a user-supplied
// ChangeFunc callback is invoked with the list of changed paths.
//
// Example usage:
//
//	w := watcher.New([]string{".env.dev", ".env.prod"}, time.Second)
//	go w.Start(func(changed []string) {
//		fmt.Println("changed:", changed)
//	})
//	// ... later ...
//	w.Stop()
package watcher
