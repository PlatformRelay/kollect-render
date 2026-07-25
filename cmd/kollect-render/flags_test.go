package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/artifact"
	"github.com/platformrelay/kollect-render/internal/format"
)

// REQ-E2-S06-01: --generated-at / --report-origin / --upstream-deps on render;
// exit 0 success, 2 fatal.

func TestRunRenderGeneratedAtOverride(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.storage.xml")
	code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--context", ctx,
		"--generated-at", "2026-07-25T12:34:56Z",
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render exit = %d, want 0", code)
	}
	sc := readSidecar(t, out)
	if sc.GeneratedAt != "2026-07-25T12:34:56Z" {
		t.Fatalf("sidecar generatedAt = %q, want override", sc.GeneratedAt)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "2026-07-25T12:34:56Z") {
		t.Fatal("body must use --generated-at override in banner")
	}
	if strings.Contains(string(body), "2026-07-21T05:00:00Z") {
		t.Fatal("body must not keep context Generation.GeneratedAt when overridden")
	}
}

func TestRunRenderReportOriginOverride(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.storage.xml")
	code := run([]string{
		"render",
		"--format", format.NameConfluenceStorage,
		"--context", ctx,
		"--report-origin", "manual",
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render exit = %d, want 0", code)
	}
	sc := readSidecar(t, out)
	if sc.Origin != "manual" {
		t.Fatalf("sidecar origin = %q, want manual", sc.Origin)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "manual") {
		t.Fatal("body must use --report-origin override in banner")
	}
}

func TestRunRenderUpstreamDepsMerge(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	// Context with no upstream block — deps file supplies Upstream map.
	ctxPath := filepath.Join(t.TempDir(), "context.yaml")
	base, err := os.ReadFile(filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	// Drop the upstream: section so only --upstream-deps populates it.
	stripped := stripYAMLSection(string(base), "upstream")
	if err := os.WriteFile(ctxPath, []byte(stripped), 0o600); err != nil {
		t.Fatal(err)
	}
	deps := filepath.Join(t.TempDir(), "upstream-deps.yaml")
	depsBody := "" +
		"event-broker:\n" +
		"  observedVersion: \"9.9.9\"\n" +
		"  sourceURL: https://example.org/event-broker\n" +
		"  query: datasource=github-releases\n" +
		"  retrievedAt: \"2026-07-25T01:00:00Z\"\n" +
		"  evidence: catalog\n" +
		"  status: observed\n"
	if err := os.WriteFile(deps, []byte(depsBody), 0o600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "out.md")
	code := run([]string{
		"render",
		"--format", format.NameMarkdown,
		"--context", ctxPath,
		"--upstream-deps", deps,
		"--output", out,
	})
	if code != 0 {
		t.Fatalf("render exit = %d, want 0", code)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "upstream 9.9.9") {
		t.Fatalf("body missing upstream from --upstream-deps:\n%s", body)
	}
}

func TestRunRenderGeneratedAtInvalid(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	code := run([]string{
		"render",
		"--format", format.NameMarkdown,
		"--context", ctx,
		"--generated-at", "not-a-timestamp",
	})
	if code != 2 {
		t.Fatalf("invalid --generated-at exit = %d, want 2", code)
	}
}

func TestRunRenderUpstreamDepsMissing(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	code := run([]string{
		"render",
		"--format", format.NameMarkdown,
		"--context", ctx,
		"--upstream-deps", filepath.Join(t.TempDir(), "no-such-upstream-deps.yaml"),
	})
	if code != 2 {
		t.Fatalf("missing --upstream-deps exit = %d, want 2", code)
	}
}

func TestRunRenderEqualsFormFlags(t *testing.T) {
	t.Parallel()
	root := repoRoot(t)
	ctx := filepath.Join(root, "test", "golden", "env-inventory-md", "context.yaml")
	out := filepath.Join(t.TempDir(), "out.storage.xml")
	code := run([]string{
		"render",
		"--format=" + format.NameConfluenceStorage,
		"--context=" + ctx,
		"--generated-at=2026-07-25T00:00:00Z",
		"--report-origin=ci",
		"--output=" + out,
	})
	if code != 0 {
		t.Fatalf("equals-form flags exit = %d, want 0", code)
	}
	sc := readSidecar(t, out)
	if sc.GeneratedAt != "2026-07-25T00:00:00Z" || sc.Origin != "ci" {
		t.Fatalf("sidecar overrides = %+v", sc)
	}
}

func readSidecar(t *testing.T, bodyPath string) artifact.Sidecar {
	t.Helper()
	raw, err := os.ReadFile(artifact.SidecarPath(bodyPath))
	if err != nil {
		t.Fatalf("read sidecar: %v", err)
	}
	var sc artifact.Sidecar
	if err := json.Unmarshal(raw, &sc); err != nil {
		t.Fatalf("sidecar JSON: %v", err)
	}
	return sc
}

// stripYAMLSection removes a top-level YAML key and its indented block (naïve; fixtures only).
func stripYAMLSection(doc, key string) string {
	lines := strings.Split(doc, "\n")
	out := make([]string, 0, len(lines))
	skip := false
	prefix := key + ":"
	for _, line := range lines {
		if skip {
			if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && !strings.HasPrefix(line, "#") {
				skip = false
			} else {
				continue
			}
		}
		if strings.HasPrefix(line, prefix) {
			skip = true
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
