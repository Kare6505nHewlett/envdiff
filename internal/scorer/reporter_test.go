package scorer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleResult() Result {
	return Result{
		Score:    82.5,
		Grade:    "B",
		Total:    20,
		Missing:  2,
		Mismatch: 1,
	}
}

func TestGradeFor(t *testing.T) {
	cases := []struct {
		score float64
		want  string
	}{
		{100, "A"}, {95, "A"}, {94, "B"}, {80, "B"},
		{79, "C"}, {65, "C"}, {64, "D"}, {50, "D"}, {49, "F"}, {0, "F"},
	}
	for _, c := range cases {
		got := gradeFor(c.score)
		if got != c.want {
			t.Errorf("gradeFor(%.0f) = %q, want %q", c.score, got, c.want)
		}
	}
}

func TestReportText_ContainsScore(t *testing.T) {
	var buf bytes.Buffer
	ReportText(&buf, sampleResult())
	out := buf.String()
	if !strings.Contains(out, "82.5") {
		t.Errorf("expected score in output, got: %s", out)
	}
	if !strings.Contains(out, "(B)") {
		t.Errorf("expected grade in output, got: %s", out)
	}
}

func TestReportText_ContainsCounts(t *testing.T) {
	var buf bytes.Buffer
	ReportText(&buf, sampleResult())
	out := buf.String()
	for _, want := range []string{"20", "2", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestReportJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	if err := ReportJSON(&buf, sampleResult()); err != nil {
		t.Fatalf("ReportJSON error: %v", err)
	}
	var got Result
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got.Score != 82.5 {
		t.Errorf("score = %.1f, want 82.5", got.Score)
	}
	if got.Grade != "B" {
		t.Errorf("grade = %q, want B", got.Grade)
	}
}

func TestDefaultOptions_Weights(t *testing.T) {
	opts := DefaultOptions()
	if opts.MissingWeight != 1.0 {
		t.Errorf("MissingWeight = %.1f, want 1.0", opts.MissingWeight)
	}
	if opts.MismatchWeight != 0.5 {
		t.Errorf("MismatchWeight = %.1f, want 0.5", opts.MismatchWeight)
	}
	if opts.MaxScore != 100.0 {
		t.Errorf("MaxScore = %.1f, want 100.0", opts.MaxScore)
	}
}
