package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRenderUsage(t *testing.T) {
	t.Parallel()
	if code := run([]string{"render"}); code != 2 {
		t.Fatalf("render usage exit = %d, want 2", code)
	}
}

func TestRunRenderGoldenMarkdown(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	tmpl := filepath.Join(root, "test", "golden", "env-inventory-md", "template.md.tmpl")
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.md")
	code := run([]string{
		"render",
		"--format", "markdown",
		"--template", tmpl,
		"--context", ctx,
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render exit = %d, want 0", code)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	want, err := os.ReadFile(filepath.Join(root, "test", "golden", "env-inventory-md", "expected.md"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("CLI render golden mismatch (got %d bytes, want %d)", len(got), len(want))
	}
}
