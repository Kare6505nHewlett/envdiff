package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_MissingFlags(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit code when flags missing")
	}
}

func TestMain_MatchingFiles(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	cmp := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	cmd := exec.Command(bin, "--base", base, "--compare", cmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 for matching files, got error: %v\noutput: %s", err, out)
	}
}

func TestMain_MismatchedFiles(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	cmp := writeTempEnv(t, "KEY=different\n")
	cmd := exec.Command(bin, "--base", base, "--compare", cmp)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit code for mismatched files")
	}
}

func TestMain_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "KEY=value\n")
	cmp := writeTempEnv(t, "KEY=value\n")
	cmd := exec.Command(bin, "--base", base, "--compare", cmp, "--format", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty JSON output")
	}
}
