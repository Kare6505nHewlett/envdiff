package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/baseline"
)

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "baseline.json")

	keys := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := baseline.Save(dest, ".env.production", keys); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoad_ReturnsCorrectData(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "baseline.json")

	keys := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := baseline.Save(dest, ".env.production", keys); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	b, err := baseline.Load(dest)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if b.SourceFile != ".env.production" {
		t.Errorf("SourceFile = %q, want %q", b.SourceFile, ".env.production")
	}
	if b.Keys["APP_ENV"] != "production" {
		t.Errorf("Keys[APP_ENV] = %q, want %q", b.Keys["APP_ENV"], "production")
	}
	if b.Keys["PORT"] != "8080" {
		t.Errorf("Keys[PORT] = %q, want %q", b.Keys["PORT"], "8080")
	}
}

func TestLoad_SetsCreatedAt(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "baseline.json")

	before := time.Now().UTC().Add(-time.Second)
	if err := baseline.Save(dest, ".env", map[string]string{}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	after := time.Now().UTC().Add(time.Second)

	b, err := baseline.Load(dest)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if b.CreatedAt.Before(before) || b.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v out of expected range [%v, %v]", b.CreatedAt, before, after)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyKeysNotNil(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "baseline.json")

	if err := baseline.Save(dest, ".env", map[string]string{}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	b, err := baseline.Load(dest)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if b.Keys == nil {
		t.Error("Keys should not be nil for empty baseline")
	}
}
