package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	baseFile := flag.String("base", "", "Base .env file (required)")
	compareFile := flag.String("compare", "", "Comparison .env file (required)")
	flag.Parse()

	if *baseFile == "" || *compareFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --base and --compare flags are required")
		flag.Usage()
		os.Exit(1)
	}

	baseEnv, err := parser.ParseFile(*baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading base file %q: %v\n", *baseFile, err)
		os.Exit(1)
	}

	cmpEnv, err := parser.ParseFile(*compareFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading compare file %q: %v\n", *compareFile, err)
		os.Exit(1)
	}

	results := diff.Compare(baseEnv, cmpEnv)

	if err := reporter.Report(os.Stdout, results, *format); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating report: %v\n", err)
		os.Exit(1)
	}

	for _, r := range results {
		if r.Status != "match" {
			os.Exit(2)
		}
	}
}
