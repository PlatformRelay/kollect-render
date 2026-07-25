package render

import (
	"os"
	"path/filepath"
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

func TestApplyGenerationOverridesInvalid(t *testing.T) {
	t.Parallel()
	_, err := ApplyGenerationOverrides(RenderContext{}, "bogus", "")
	if err == nil {
		t.Fatal("expected error for invalid generated-at")
	}
}
