package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	base := flag.String("base", "", "path to base .env file (required)")
	target := flag.String("target", "", "path to target .env file (required)")
	format := flag.String("format", "text", "output format: text or json")
	onlyMissing := flag.Bool("only-missing", false, "show only missing keys")
	onlyMismatch := flag.Bool("only-mismatch", false, "show only mismatched keys")
	prefix := flag.String("prefix", "", "filter keys by prefix")
	exclude := flag.String("exclude", "", "comma-separated list of keys to exclude")
	flag.Parse()

	if *base == "" || *target == "" {
		fmt.Fprintln(os.Stderr, "error: --base and --target are required")
		flag.Usage()
		os.Exit(1)
	}

	baseMap, err := parser.ParseFile(*base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		os.Exit(1)
	}

	targetMap, err := parser.ParseFile(*target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target file: %v\n", err)
		os.Exit(1)
	}

	results := diff.Compare(baseMap, targetMap)

	var excludeKeys []string
	if *exclude != "" {
		for _, k := range strings.Split(*exclude, ",") {
			excludeKeys = append(excludeKeys, strings.TrimSpace(k))
		}
	}

	results = filter.Apply(results, filter.Options{
		OnlyMissing:  *onlyMissing,
		OnlyMismatch: *onlyMismatch,
		Prefix:       *prefix,
		ExcludeKeys:  excludeKeys,
	})

	if err := reporter.Report(os.Stdout, results, *format); err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		os.Exit(1)
	}
}
