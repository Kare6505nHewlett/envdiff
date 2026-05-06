package patcher_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/patcher"
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

func TestApply_AddsNewKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	results, err := patcher.Apply(path, []patcher.Patch{{Key: "BAZ", Value: "qux"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Action != "added" {
		t.Fatalf("expected added, got %+v", results)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "BAZ=qux") {
		t.Errorf("expected BAZ=qux in file, got: %s", data)
	}
}

func TestApply_UpdatesExistingKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=old\nBAR=keep\n")
	results, err := patcher.Apply(path, []patcher.Patch{{Key: "FOO", Value: "new"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Action != "updated" {
		t.Fatalf("expected updated, got %+v", results)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "FOO=new") {
		t.Errorf("expected FOO=new, got: %s", data)
	}
	if !strings.Contains(string(data), "BAR=keep") {
		t.Errorf("expected BAR=keep preserved, got: %s", data)
	}
}

func TestApply_CreatesFilIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.env")
	_, err := patcher.Apply(path, []patcher.Patch{{Key: "KEY", Value: "val"}})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "KEY=val") {
		t.Errorf("expected KEY=val, got: %s", data)
	}
}

func TestDryRun_DoesNotWriteFile(t *testing.T) {
	path := writeTempEnv(t, "EXISTING=yes\n")
	original, _ := os.ReadFile(path)

	results, err := patcher.DryRun(path, []patcher.Patch{
		{Key: "EXISTING", Value: "no"},
		{Key: "NEW", Value: "val"},
	})
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(original) != string(after) {
		t.Error("DryRun must not modify the file")
	}
	actions := map[string]string{}
	for _, r := range results {
		actions[r.Key] = r.Action
	}
	if actions["EXISTING"] != "updated" {
		t.Errorf("expected updated for EXISTING, got %s", actions["EXISTING"])
	}
	if actions["NEW"] != "added" {
		t.Errorf("expected added for NEW, got %s", actions["NEW"])
	}
}

func TestApply_QuotesValueWithSpaces(t *testing.T) {
	path := writeTempEnv(t, "")
	_, err := patcher.Apply(path, []patcher.Patch{{Key: "MSG", Value: "hello world"}})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `MSG=`) {
		t.Errorf("expected MSG= in output, got: %s", data)
	}
}
