package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/platformrelay/kollect-render/internal/artifact"
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

func TestRunRenderWritesArtifactSidecar(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "page.storage.xml")
	code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--context", ctx,
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render exit = %d, want 0", code)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	sidePath := artifact.SidecarPath(out)
	raw, err := os.ReadFile(sidePath)
	if err != nil {
		t.Fatalf("sidecar missing at %s: %v", sidePath, err)
	}
	var sc artifact.Sidecar
	if err := json.Unmarshal(raw, &sc); err != nil {
		t.Fatalf("sidecar JSON: %v", err)
	}
	sum := sha256.Sum256(body)
	wantDigest := "sha256:" + hex.EncodeToString(sum[:])
	if sc.ContentDigest != wantDigest {
		t.Fatalf("sidecar contentDigest = %q, want %q", sc.ContentDigest, wantDigest)
	}
	if sc.Format != format.NameConfluenceStorage {
		t.Fatalf("sidecar format = %q", sc.Format)
	}
	if sc.ByteLength != len(body) {
		t.Fatalf("sidecar byteLength = %d, want %d", sc.ByteLength, len(body))
	}
	if sc.GeneratedAt != "2026-07-21T05:00:00Z" {
		t.Fatalf("sidecar generatedAt = %q (from fixture Generation)", sc.GeneratedAt)
	}
	if sc.Origin != "schedule" || sc.SnapshotSHA != "abcdef0123456789" {
		t.Fatalf("sidecar lineage = %+v", sc)
	}
	if sc.RendererVersion != version {
		t.Fatalf("sidecar rendererVersion = %q, want %q", sc.RendererVersion, version)
	}
	// Determinism: second render → identical body + sidecar.
	out2 := filepath.Join(t.TempDir(), "page.storage.xml")
	if code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--context", ctx,
		"--output", out2,
	}); code != 0 {
		t.Fatalf("second render exit = %d", code)
	}
	body2, err := os.ReadFile(out2)
	if err != nil {
		t.Fatal(err)
	}
	side2, err := os.ReadFile(artifact.SidecarPath(out2))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != string(body2) {
		t.Fatal("body not deterministic across renders")
	}
	if string(raw) != string(side2) {
		t.Fatal("sidecar not deterministic across renders")
	}
}

func TestRunRenderStdoutSkipsSidecar(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--context", ctx,
		"--output", "-",
	})
	if code != 0 {
		t.Fatalf("stdout render exit = %d, want 0", code)
	}
}
