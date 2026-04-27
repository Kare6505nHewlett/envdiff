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

func TestMain_MissingFlags(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "--base and --target are required") {
		t.Fatalf("expected usage error, got: %s", out)
	}
}

func TestMain_MatchingFiles(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP=hello\nDB=world\n")
	target := writeTempEnv(t, "APP=hello\nDB=world\n")
	cmd := exec.Command(bin, "--base", base, "--target", target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "match") {
		t.Fatalf("expected match output, got: %s", out)
	}
}

func TestMain_MismatchedFiles(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP=hello\nDB=world\n")
	target := writeTempEnv(t, "APP=hello\nDB=different\n")
	cmd := exec.Command(bin, "--base", base, "--target", target)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "mismatch") {
		t.Fatalf("expected mismatch output, got: %s", out)
	}
}

func TestMain_FilterByPrefix(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempEnv(t, "APP_NAME=foo\nDB_HOST=localhost\n")
	target := writeTempEnv(t, "APP_NAME=foo\nDB_HOST=remotehost\n")
	cmd := exec.Command(bin, "--base", base, "--target", target, "--prefix", "DB_")
	out, _ := cmd.CombinedOutput()
	if strings.Contains(string(out), "APP_NAME") {
		t.Fatalf("APP_NAME should be filtered out, got: %s", out)
	}
	if !strings.Contains(string(out), "DB_HOST") {
		t.Fatalf("DB_HOST should appear, got: %s", out)
	}
}
