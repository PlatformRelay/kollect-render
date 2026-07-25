package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadUpstreamFile(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "upstream-deps.yaml")
	body := "" +
		"event-broker:\n" +
		"  observedVersion: \"3.10.0\"\n" +
		"  sourceURL: https://example.org/event-broker\n" +
		"  query: datasource=github-releases\n" +
		"  retrievedAt: \"2026-07-21T03:00:00Z\"\n" +
		"  evidence: catalog\n" +
		"  status: observed\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadUpstreamFile(path)
	if err != nil {
		t.Fatalf("LoadUpstreamFile: %v", err)
	}
	e, ok := got["event-broker"]
	if !ok {
		t.Fatalf("missing event-broker: %+v", got)
	}
	if e.ObservedVersion != "3.10.0" || e.Status != "observed" {
		t.Fatalf("entry = %+v", e)
	}
	want := time.Date(2026, 7, 21, 3, 0, 0, 0, time.UTC)
	if !e.RetrievedAt.Equal(want) {
		t.Fatalf("RetrievedAt = %v, want %v", e.RetrievedAt, want)
	}
}

func TestLoadUpstreamFileMissing(t *testing.T) {
	t.Parallel()
	_, err := LoadUpstreamFile(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadUpstreamFileYAMLFail(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "bad-upstream.yaml")
	if err := os.WriteFile(path, []byte("event-broker: [\n  - !!broken"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadUpstreamFile(path)
	if err == nil {
		t.Fatal("expected YAML decode error")
	}
	if !strings.Contains(err.Error(), "decode upstream-deps") {
		t.Fatalf("error = %v, want decode upstream-deps prefix", err)
	}
}

func TestLoadUpstreamFileEmptyDocument(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "empty-upstream.yaml")
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadUpstreamFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("nil map must be initialized to empty")
	}
	if len(got) != 0 {
		t.Fatalf("got = %+v, want empty", got)
	}
}

func TestApplyGenerationOverrides(t *testing.T) {
	t.Parallel()
	ctx := RenderContext{
		Generation: GenerationMeta{
			GeneratedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			Origin:      "schedule",
		},
	}
	got, err := ApplyGenerationOverrides(ctx, "2026-07-25T12:00:00Z", "manual")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	if !got.Generation.GeneratedAt.Equal(want) {
		t.Fatalf("GeneratedAt = %v", got.Generation.GeneratedAt)
	}
	if got.Generation.Origin != "manual" {
		t.Fatalf("Origin = %q", got.Generation.Origin)
	}
}

func TestApplyGenerationOverridesEmptyLeaveUnchanged(t *testing.T) {
	t.Parallel()
	ctx := RenderContext{
		Generation: GenerationMeta{
			GeneratedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			Origin:      "schedule",
		},
	}
	got, err := ApplyGenerationOverrides(ctx, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Generation.GeneratedAt.Equal(ctx.Generation.GeneratedAt) || got.Generation.Origin != "schedule" {
		t.Fatalf("empty overrides mutated context: %+v", got.Generation)
	}
}

func TestApplyGenerationOverridesInvalid(t *testing.T) {
	t.Parallel()
	_, err := ApplyGenerationOverrides(RenderContext{}, "bogus", "")
	if err == nil {
		t.Fatal("expected error for invalid generated-at")
	}
}
