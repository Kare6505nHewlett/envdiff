// Package profiler computes a letter-grade health profile for a set of
// environment diff results.
//
// # Overview
//
// Given a slice of [diff.Result] values, Compute returns a [Profile] that
// summarises the overall quality of the environment configuration:
//
//	  p := profiler.Compute(results)
//	  fmt.Println(p.Grade, p.Score)
//
// # Grading
//
// Scores are calculated by penalising missing keys (weight 2) and mismatched
// keys (weight 1) relative to the total number of keys.  The resulting
// percentage maps to a letter grade:
//
//	  A  ≥ 95   B  ≥ 80   C  ≥ 65   D  ≥ 50   F  < 50
//
// # Reporting
//
// [ReportText] and [ReportJSON] write the profile to any [io.Writer].
package profiler
