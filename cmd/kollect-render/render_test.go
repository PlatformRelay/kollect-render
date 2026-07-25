package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/platformrelay/kollect-render/internal/format"
)

func TestRunRenderUsage(t *testing.T) {
	t.Parallel()
	if code := run([]string{"render"}); code != 2 {
		t.Fatalf("render usage exit = %d, want 2", code)
	}
}

func TestRunRenderUnknownFormat(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	tmpl := filepath.Join(root, "test", "golden", "env-inventory-md", "template.md.tmpl")
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	code := run([]string{
		"render",
		"--format", "pdf",
		"--template", tmpl,
		"--context", ctx,
	})
	if code != 2 {
		t.Fatalf("unknown format exit = %d, want 2", code)
	}
}

func TestRunRenderRejectsTemplateWithNonMarkdownFormat(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	tmpl := filepath.Join(root, "test", "golden", "env-inventory-md", "template.md.tmpl")
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.storage.xml")
	code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--template", tmpl,
		"--context", ctx,
		"--output", out,
	})
	if code != 2 {
		t.Fatalf("confluence-storage+template exit = %d, want 2 (must not silently emit markdown)", code)
	}
	if _, err := os.Stat(out); err == nil {
		t.Fatal("must not write output when format+template are incompatible")
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

func TestRunRenderModelConfluenceStorage(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.storage.xml")
	code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--context", ctx,
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render confluence-storage exit = %d, want 0", code)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	want, err := os.ReadFile(filepath.Join(root, "test", "golden", "env-inventory-cs", "expected.storage.xml"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("CLI confluence-storage golden mismatch (got %d bytes, want %d)", len(got), len(want))
	}
}

func TestRunRenderModelMarkdown(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.md")
	code := run([]string{
		"render",
		"--format", format.NameMarkdown,
		"--context", ctx,
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render model-markdown exit = %d, want 0", code)
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
		t.Fatalf("CLI model-markdown golden mismatch")
	}
}
