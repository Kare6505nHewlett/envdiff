// Package differ provides fine-grained, line-level diffing between two
// parsed .env maps. Unlike the top-level diff package which works on
// multi-file comparison results, differ operates directly on two
// map[string]string values and produces Hunk records that describe
// exactly what was added, removed, or changed between them.
//
// Typical usage:
//
//	mapA, _ := parser.ParseFile("staging.env")
//	mapB, _ := parser.ParseFile("production.env")
//	hunks := differ.Diff("staging.env", "production.env", mapA, mapB)
//	for _, h := range hunks {
//		fmt.Println(h)
//	}
//	fmt.Println(differ.Summary(hunks))
package differ
