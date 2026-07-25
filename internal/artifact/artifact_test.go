package artifact_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/artifact"
)

func TestContentDigestSHA256(t *testing.T) {
	t.Parallel()
	body := []byte("hello-artifact\n")
	sum := sha256.Sum256(body)
	want := "sha256:" + hex.EncodeToString(sum[:])
	got := artifact.ContentDigest(body)
	if got != want {
		t.Fatalf("ContentDigest = %q, want %q", got, want)
	}
	if artifact.ContentDigest(nil) != artifact.ContentDigest([]byte{}) {
		t.Fatal("nil and empty body digests must match")
	}
}

func TestBuildSidecarFields(t *testing.T) {
	t.Parallel()
	body := []byte("<p>storage</p>")
	meta := artifact.Meta{
		Format:          "confluence-storage",
		GeneratedAt:     "2026-07-21T12:00:00Z",
		Origin:          "manual",
		SnapshotSHA:     "abc123",
		SourceRepoURL:   "https://example.invalid/snapshots",
		TemplateDigest:  "sha256:" + strings.Repeat("c", 64),
		RendererVersion: "0.0.0-dev",
	}
	sc := artifact.Build(body, meta)
	if sc.SchemaVersion != artifact.SchemaVersion {
		t.Fatalf("SchemaVersion = %q, want %q", sc.SchemaVersion, artifact.SchemaVersion)
	}
	if sc.ContentDigest != artifact.ContentDigest(body) {
		t.Fatalf("ContentDigest = %q, want %q", sc.ContentDigest, artifact.ContentDigest(body))
	}
	if sc.ByteLength != len(body) {
		t.Fatalf("ByteLength = %d, want %d", sc.ByteLength, len(body))
	}
	if sc.Format != meta.Format || sc.GeneratedAt != meta.GeneratedAt || sc.Origin != meta.Origin {
		t.Fatalf("meta fields not copied: %+v", sc)
	}
	if sc.SnapshotSHA != meta.SnapshotSHA || sc.SourceRepoURL != meta.SourceRepoURL {
		t.Fatalf("lineage fields not copied: %+v", sc)
	}
	if sc.TemplateDigest != meta.TemplateDigest || sc.RendererVersion != meta.RendererVersion {
		t.Fatalf("version fields not copied: %+v", sc)
	}
}

func TestMarshalSidecarDeterministic(t *testing.T) {
	t.Parallel()
	body := []byte("same-bytes")
	meta := artifact.Meta{Format: "confluence-storage", GeneratedAt: "2026-07-21T12:00:00Z", Origin: "schedule"}
	a, err := artifact.MarshalSidecar(artifact.Build(body, meta))
	if err != nil {
		t.Fatalf("MarshalSidecar: %v", err)
	}
	b, err := artifact.MarshalSidecar(artifact.Build(body, meta))
	if err != nil {
		t.Fatalf("MarshalSidecar 2: %v", err)
	}
	if string(a) != string(b) {
		t.Fatalf("sidecar JSON not deterministic:\n%s\nvs\n%s", a, b)
	}
	var decoded artifact.Sidecar
	if err := json.Unmarshal(a, &decoded); err != nil {
		t.Fatalf("sidecar must be valid JSON: %v", err)
	}
	if decoded.ContentDigest != artifact.ContentDigest(body) {
		t.Fatalf("decoded digest = %q", decoded.ContentDigest)
	}
}

func TestWriteBodyAndSidecar(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	bodyPath := filepath.Join(dir, "page.storage.xml")
	body := []byte("<ac:structured-macro/>")
	meta := artifact.Meta{
		Format:      "confluence-storage",
		GeneratedAt: "2026-07-21T12:00:00Z",
		Origin:      "manual",
	}
	sc, err := artifact.Write(bodyPath, body, meta)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	gotBody, err := os.ReadFile(bodyPath)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if string(gotBody) != string(body) {
		t.Fatalf("body mismatch")
	}
	sidePath := artifact.SidecarPath(bodyPath)
	if sidePath != bodyPath+artifact.SidecarSuffix {
		t.Fatalf("SidecarPath = %q, want %q", sidePath, bodyPath+artifact.SidecarSuffix)
	}
	raw, err := os.ReadFile(sidePath)
	if err != nil {
		t.Fatalf("read sidecar: %v", err)
	}
	var got artifact.Sidecar
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal sidecar: %v", err)
	}
	if got.ContentDigest != sc.ContentDigest || got.ContentDigest != artifact.ContentDigest(body) {
		t.Fatalf("sidecar digest = %q, want %q", got.ContentDigest, artifact.ContentDigest(body))
	}
	if got.Format != "confluence-storage" || got.ByteLength != len(body) {
		t.Fatalf("sidecar meta = %+v", got)
	}
}

func TestWriteDeterministicAcrossRuns(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	body := []byte("deterministic-body")
	meta := artifact.Meta{Format: "markdown", GeneratedAt: "2026-07-21T12:00:00Z", Origin: "schedule"}

	p1 := filepath.Join(dir, "a.md")
	p2 := filepath.Join(dir, "b.md")
	if _, err := artifact.Write(p1, body, meta); err != nil {
		t.Fatalf("Write 1: %v", err)
	}
	if _, err := artifact.Write(p2, body, meta); err != nil {
		t.Fatalf("Write 2: %v", err)
	}
	s1, err := os.ReadFile(artifact.SidecarPath(p1))
	if err != nil {
		t.Fatal(err)
	}
	s2, err := os.ReadFile(artifact.SidecarPath(p2))
	if err != nil {
		t.Fatal(err)
	}
	if string(s1) != string(s2) {
		t.Fatalf("sidecar bytes differ across writes")
	}
}
