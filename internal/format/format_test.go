package format_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/format"
	"github.com/platformrelay/kollect-render/internal/render"
)

func TestRegistryHasBuiltinFormats(t *testing.T) {
	t.Parallel()
	names := format.Names()
	if len(names) < 2 {
		t.Fatalf("Names() = %v, want at least markdown + confluence-storage", names)
	}
	for _, want := range []string{format.NameMarkdown, format.NameConfluenceStorage} {
		if _, ok := format.Lookup(want); !ok {
			t.Fatalf("Lookup(%q) missing", want)
		}
	}
	if _, ok := format.Lookup("nope"); ok {
		t.Fatal("Lookup(nope) should be false")
	}
}

func TestSameModelEncodesToBothGoldens(t *testing.T) {
	t.Parallel()
	ctx := mustLoadContext(t, filepath.Join(goldenRoot(t), "env-inventory-md", "context.yaml"))
	model, err := format.EnvInventoryModel(ctx)
	if err != nil {
		t.Fatalf("EnvInventoryModel: %v", err)
	}

	mdEnc, ok := format.Lookup(format.NameMarkdown)
	if !ok {
		t.Fatal("markdown encoder missing")
	}
	csEnc, ok := format.Lookup(format.NameConfluenceStorage)
	if !ok {
		t.Fatal("confluence-storage encoder missing")
	}

	gotMD, err := mdEnc.Encode(model)
	if err != nil {
		t.Fatalf("markdown Encode: %v", err)
	}
	gotCS, err := csEnc.Encode(model)
	if err != nil {
		t.Fatalf("confluence-storage Encode: %v", err)
	}

	// Determinism: encode twice.
	gotMD2, err := mdEnc.Encode(model)
	if err != nil {
		t.Fatalf("markdown Encode 2: %v", err)
	}
	gotCS2, err := csEnc.Encode(model)
	if err != nil {
		t.Fatalf("confluence-storage Encode 2: %v", err)
	}
	if string(gotMD) != string(gotMD2) {
		t.Fatal("markdown encode not deterministic")
	}
	if string(gotCS) != string(gotCS2) {
		t.Fatal("confluence-storage encode not deterministic")
	}

	mdGolden := filepath.Join(goldenRoot(t), "env-inventory-md", "expected.md")
	csGolden := filepath.Join(goldenRoot(t), "env-inventory-cs", "expected.storage.xml")

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(mdGolden, gotMD, 0o644); err != nil {
			t.Fatalf("update md golden: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(csGolden), 0o755); err != nil {
			t.Fatalf("mkdir cs golden: %v", err)
		}
		if err := os.WriteFile(csGolden, gotCS, 0o644); err != nil {
			t.Fatalf("update cs golden: %v", err)
		}
		return
	}

	wantMD := mustRead(t, mdGolden)
	if string(gotMD) != wantMD {
		t.Fatalf("markdown golden mismatch (-want +got):\n%s", diffSnippet(wantMD, string(gotMD)))
	}
	wantCS := mustRead(t, csGolden)
	if string(gotCS) != wantCS {
		t.Fatalf("confluence-storage golden mismatch (-want +got):\n%s", diffSnippet(wantCS, string(gotCS)))
	}
}

func TestConfluenceEscapesHostileMarkup(t *testing.T) {
	t.Parallel()
	ctx := mustLoadContext(t, filepath.Join(goldenRoot(t), "env-inventory-md", "context.yaml"))
	model, err := format.EnvInventoryModel(ctx)
	if err != nil {
		t.Fatalf("EnvInventoryModel: %v", err)
	}
	enc, _ := format.Lookup(format.NameConfluenceStorage)
	out, err := enc.Encode(model)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "<script>") {
		t.Fatal("confluence-storage must not emit raw <script> (contextual escaping required)")
	}
	if !strings.Contains(s, "&lt;script&gt;") {
		t.Fatal("expected escaped script tag in confluence-storage output")
	}
}

func TestConfluenceStoragePublishReadyShape(t *testing.T) {
	t.Parallel()
	ctx := mustLoadContext(t, filepath.Join(goldenRoot(t), "env-inventory-md", "context.yaml"))
	model, err := format.EnvInventoryModel(ctx)
	if err != nil {
		t.Fatalf("EnvInventoryModel: %v", err)
	}
	enc, _ := format.Lookup(format.NameConfluenceStorage)
	out, err := enc.Encode(model)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	s := string(out)
	checks := []string{
		`ac:name="panel"`,
		`<table`,
		`ac:name="status"`,
		`<h1>`,
	}
	for _, c := range checks {
		if !strings.Contains(s, c) {
			t.Fatalf("publish-ready shape missing %q", c)
		}
	}
	if err := format.ValidateStorageFragment(s); err != nil {
		t.Fatalf("storage fragment not well-formed: %v", err)
	}
}

func TestValidateStorageFragmentRejectsBroken(t *testing.T) {
	t.Parallel()
	if err := format.ValidateStorageFragment(`<p>unclosed`); err == nil {
		t.Fatal("expected error for unclosed tag")
	}
	if err := format.ValidateStorageFragment(""); err == nil {
		t.Fatal("expected error for empty fragment")
	}
}

func TestNoRawSafeBypassInAPI(t *testing.T) {
	t.Parallel()
	// Allowed Inline concrete types only — extending requires an escaped encoder arm.
	allowed := map[string]struct{}{
		"Text": {}, "Strong": {}, "Code": {}, "Link": {}, "Emoji": {},
	}
	for _, in := range []format.Inline{
		format.Text{}, format.Strong{}, format.Code{}, format.Link{}, format.Emoji{},
	} {
		name := reflect.TypeOf(in).Name()
		if _, ok := allowed[name]; !ok {
			t.Fatalf("unexpected Inline concrete type %q", name)
		}
		delete(allowed, name)
	}
	if len(allowed) != 0 {
		t.Fatalf("missing Inline coverage: %v", allowed)
	}

	// Source guard: model.go must not grow a Raw/HTML/Unsafe bypass type.
	src := mustRead(t, filepath.Join(repoRoot(t), "internal", "format", "model.go"))
	for _, bad := range []string{
		"type Raw ", "type HTML ", "type Unsafe", "template.HTML", "type RawHTML",
	} {
		if strings.Contains(src, bad) {
			t.Fatalf("raw-safe bypass marker %q found in model.go", bad)
		}
	}
}

func goldenRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "test", "golden")
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func mustLoadContext(t *testing.T, path string) render.RenderContext {
	t.Helper()
	ctx, err := render.LoadContextFile(path)
	if err != nil {
		t.Fatalf("LoadContextFile: %v", err)
	}
	return ctx
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	dir := wd
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "go.mod")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("cannot find go.mod")
	return ""
}

func diffSnippet(want, got string) string {
	const limit = 800
	w, g := want, got
	if len(w) > limit {
		w = w[:limit] + "…"
	}
	if len(g) > limit {
		g = g[:limit] + "…"
	}
	return "want:\n" + w + "\n\ngot:\n" + g
}
