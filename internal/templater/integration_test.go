package templater_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/templater"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
	return p
}

func TestIntegration_TemplateFromRealFiles(t *testing.T) {
	dir := t.TempDir()

	baseFile := writeEnv(t, dir, "base.env", "DB_HOST=localhost\nAPI_KEY=secret\nPORT=8080\n")
	otherFile := writeEnv(t, dir, "other.env", "DB_HOST=prod-host\nPORT=8080\n")

	base, err := parser.ParseFile(baseFile)
	if err != nil {
		t.Fatalf("parse base: %v", err)
	}
	other, err := parser.ParseFile(otherFile)
	if err != nil {
		t.Fatalf("parse other: %v", err)
	}

	results := diff.Compare(base, other, otherFile)

	out := filepath.Join(dir, ".env.template")
	err = templater.WriteFile(out, results, templater.Options{
		Placeholder:     "FILL_ME",
		IncludeComments: true,
	})
	if err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)

	for _, key := range []string{"API_KEY", "DB_HOST", "PORT"} {
		if !strings.Contains(content, key+"=FILL_ME") {
			t.Errorf("expected %s=FILL_ME in template:\n%s", key, content)
		}
	}

	if !strings.Contains(content, "# status:") {
		t.Errorf("expected status comments in template:\n%s", content)
	}
}
