package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_BasicKeyValues(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, env, "APP_ENV", "production")
	assertKey(t, env, "DB_HOST", "localhost")
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Fatalf("expected 1 key, got %d", len(env))
	}
	assertKey(t, env, "KEY", "value")
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, env, "SECRET", "my secret value")
	assertKey(t, env, "TOKEN", "abc123")
}

func TestParseFile_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export REGION=us-east-1\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, env, "REGION", "us-east-1")
}

func TestParseFile_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # http port\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertKey(t, env, "PORT", "8080")
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func assertKey(t *testing.T, env EnvMap, key, want string) {
	t.Helper()
	got, ok := env[key]
	if !ok {
		t.Errorf("key %q not found in env map", key)
		return
	}
	if got != want {
		t.Errorf("key %q: got %q, want %q", key, got, want)
	}
}
