package artifact_test

import (
	"path/filepath"
	"testing"

	"github.com/platformrelay/kollect-render/internal/artifact"
)

func TestWriteEmptyPath(t *testing.T) {
	t.Parallel()
	if _, err := artifact.Write("", []byte("x"), artifact.Meta{Format: "markdown"}); err == nil {
		t.Fatal("Write empty path must error")
	}
}

func TestWriteNilBody(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "empty.md")
	sc, err := artifact.Write(path, nil, artifact.Meta{Format: "markdown"})
	if err != nil {
		t.Fatal(err)
	}
	if sc.ByteLength != 0 {
		t.Fatalf("ByteLength = %d", sc.ByteLength)
	}
	if sc.ContentDigest != artifact.ContentDigest(nil) {
		t.Fatalf("digest = %q", sc.ContentDigest)
	}
}

func TestBuildNilBody(t *testing.T) {
	t.Parallel()
	sc := artifact.Build(nil, artifact.Meta{Format: "markdown"})
	if sc.ByteLength != 0 || sc.ContentDigest != artifact.ContentDigest([]byte{}) {
		t.Fatalf("Build(nil) = %+v", sc)
	}
}

func TestWriteUnwritableDir(t *testing.T) {
	t.Parallel()
	// Parent does not exist → Write must surface write body error.
	path := filepath.Join(t.TempDir(), "missing-dir", "out.md")
	if _, err := artifact.Write(path, []byte("x"), artifact.Meta{Format: "markdown"}); err == nil {
		t.Fatal("want write body error for missing parent")
	}
}
