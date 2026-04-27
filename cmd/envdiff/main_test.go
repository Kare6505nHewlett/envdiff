package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %s", out)
	}
	return bin
}

// runEnvdiff is a helper that runs the envdiff binary with the given arguments
// and returns the combined stdout/stderr output and the exit error (if any).
func runEnvdiff(t *testing.T, bin string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestMain_MissingFlags(t *testing.T) {
	bin := buildBinary(t)
	out, _ := runEnvdiff(t, bin)
	if !strings.Contains(out, "--base and --target are required") {
		t.Fatalf("expected usage error, got: %s", out)
	}
}

func TestMain_MatchingFiles(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP=hello\nDB=world\n")
	target := writeTempEnv(t, "APP=hello\nDB=world\n")
	out, err := runEnvdiff(t, bin, "--base", base, "--target", target)
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "match") {
		t.Fatalf("expected match output, got: %s", out)
	}
}

func TestMain_MismatchedFiles(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP=hello\nDB=world\n")
	target := writeTempEnv(t, "APP=hello\nDB=different\n")
	out, _ := runEnvdiff(t, bin, "--base", base, "--target", target)
	if !strings.Contains(out, "mismatch") {
		t.Fatalf("expected mismatch output, got: %s", out)
	}
}

func TestMain_FilterByPrefix(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP_NAME=foo\nDB_HOST=localhost\n")
	target := writeTempEnv(t, "APP_NAME=foo\nDB_HOST=remotehost\n")
	out, _ := runEnvdiff(t, bin, "--base", base, "--target", target, "--prefix", "DB_")
	if strings.Contains(out, "APP_NAME") {
		t.Fatalf("APP_NAME should be filtered out, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Fatalf("DB_HOST should appear, got: %s", out)
	}
}
