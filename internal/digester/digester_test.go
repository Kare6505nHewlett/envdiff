package digester_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/digester"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", File: "staging.env", Status: diff.StatusMatch},
		{Key: "DB_HOST", File: "staging.env", Status: diff.StatusMatch},
		{Key: "APP_ENV", File: "prod.env", Status: diff.StatusMatch},
		{Key: "DB_HOST", File: "prod.env", Status: diff.StatusMismatch},
	}
}

func TestCompute_ReturnsOneDigestPerFile(t *testing.T) {
	digests := digester.Compute(sampleResults())
	if len(digests) != 2 {
		t.Fatalf("expected 2 digests, got %d", len(digests))
	}
}

func TestCompute_DigestsSortedByFile(t *testing.T) {
	digests := digester.Compute(sampleResults())
	if digests[0].File != "prod.env" {
		t.Errorf("expected prod.env first, got %s", digests[0].File)
	}
	if digests[1].File != "staging.env" {
		t.Errorf("expected staging.env second, got %s", digests[1].File)
	}
}

func TestCompute_SameKeysSameHash(t *testing.T) {
	digests := digester.Compute(sampleResults())
	if digests[0].Hash != digests[1].Hash {
		t.Errorf("expected identical hashes for same key sets, got %s vs %s",
			digests[0].Hash, digests[1].Hash)
	}
}

func TestCompute_DifferentKeysProduceDifferentHash(t *testing.T) {
	results := []diff.Result{
		{Key: "APP_ENV", File: "a.env"},
		{Key: "APP_ENV", File: "b.env"},
		{Key: "EXTRA_KEY", File: "b.env"},
	}
	digests := digester.Compute(results)
	if digests[0].Hash == digests[1].Hash {
		t.Error("expected different hashes for different key sets")
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	digests := digester.Compute(nil)
	if digests != nil {
		t.Errorf("expected nil for empty input, got %v", digests)
	}
}

func TestMatch_AllSameHash(t *testing.T) {
	digests := digester.Compute(sampleResults())
	if !digester.Match(digests) {
		t.Error("expected Match to return true for identical key sets")
	}
}

func TestMatch_DifferentHashes(t *testing.T) {
	results := []diff.Result{
		{Key: "A", File: "x.env"},
		{Key: "B", File: "y.env"},
	}
	digests := digester.Compute(results)
	if digester.Match(digests) {
		t.Error("expected Match to return false for different key sets")
	}
}

func TestMatch_EmptySlice(t *testing.T) {
	if !digester.Match(nil) {
		t.Error("expected Match to return true for empty input")
	}
}

func TestDiff_ReturnsEmptyWhenAllMatch(t *testing.T) {
	digests := digester.Compute(sampleResults())
	out := digester.Diff(digests)
	if out != "" {
		t.Errorf("expected empty diff output, got %q", out)
	}
}

func TestDiff_DescribesDivergentFile(t *testing.T) {
	results := []diff.Result{
		{Key: "A", File: "base.env"},
		{Key: "A", File: "other.env"},
		{Key: "B", File: "other.env"},
	}
	digests := digester.Compute(results)
	out := digester.Diff(digests)
	if !strings.Contains(out, "other.env") {
		t.Errorf("expected diff output to mention other.env, got %q", out)
	}
}
