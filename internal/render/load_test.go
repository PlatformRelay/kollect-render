package render_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/render"
)

func TestDecodeContextYAMLFail(t *testing.T) {
	t.Parallel()
	_, err := render.DecodeContext([]byte("upstream: [\n  - !!broken"))
	if err == nil {
		t.Fatal("DecodeContext: want YAML decode error, got nil")
	}
	if !strings.Contains(err.Error(), "decode context") {
		t.Fatalf("error = %v, want decode context prefix", err)
	}
}

func TestLoadContextFileMissing(t *testing.T) {
	t.Parallel()
	_, err := render.LoadContextFile(filepath.Join(t.TempDir(), "missing-context.yaml"))
	if err == nil {
		t.Fatal("LoadContextFile: want error for missing file")
	}
}

func TestDecodeContextNilMapsInitialized(t *testing.T) {
	t.Parallel()
	ctx, err := render.DecodeContext([]byte("generation:\n  origin: schedule\n"))
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Upstream == nil || ctx.Copy == nil {
		t.Fatalf("nil maps must be initialized: upstream=%v copy=%v", ctx.Upstream, ctx.Copy)
	}
}

func TestLoadContextFileRoundTrip(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "ctx.yaml")
	body := "generation:\n  origin: manual\nupstream: {}\ncopy: {}\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, err := render.LoadContextFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Generation.Origin != "manual" {
		t.Fatalf("Origin = %q", ctx.Generation.Origin)
	}
}
