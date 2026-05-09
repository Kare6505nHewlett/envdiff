// Package scheduler implements periodic background diffing of .env files.
//
// It wraps the diff.Compare function and calls a user-supplied Handler
// whenever the comparison results change between ticks.  Callers control
// the polling interval and can opt in to always-notify mode so the handler
// fires on every tick regardless of whether results have changed.
//
// Typical usage:
//
//	s := scheduler.New(
//		[]string{".env.staging", ".env.production"},
//		30*time.Second,
//		func(results []diff.Result) { reporter.Report(os.Stdout, results, "text") },
//		scheduler.DefaultOptions(),
//	)
//	s.Run(ctx)
package scheduler
