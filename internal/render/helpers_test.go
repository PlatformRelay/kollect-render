package render_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/platformrelay/kollect-render/internal/render"
)

func TestHelperSetFrozenNames(t *testing.T) {
	t.Parallel()
	want := []string{
		"semverCompare",
		"versionDifference",
		"statusEmoji",
		"fmtTime",
		"fmtDate",
		"anchor",
		"sortBy",
		"groupBy",
		"capNote",
	}
	if len(render.HelperNames) != len(want) {
		t.Fatalf("HelperNames len = %d, want %d", len(render.HelperNames), len(want))
	}
	for i, name := range want {
		if render.HelperNames[i] != name {
			t.Fatalf("HelperNames[%d] = %q, want %q", i, render.HelperNames[i], name)
		}
	}
}

func TestSemverCompare(t *testing.T) {
	t.Parallel()
	cases := []struct {
		a, b string
		want int
	}{
		{"1.2.3", "1.2.3", 0},
		{"v1.2.3", "1.2.3", 0},
		{"1.2.3+build.1", "1.2.3+build.2", 0},
		{"1.2.3", "1.2.4", -1},
		{"2.0.0", "1.9.9", 1},
		{"v0.47.0", "0.46.0", 1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.a+"_vs_"+tc.b, func(t *testing.T) {
			t.Parallel()
			got := render.SemverCompare(tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("SemverCompare(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestVersionDifference(t *testing.T) {
	t.Parallel()
	cases := []struct {
		current, observed, want string
	}{
		{"1.2.3", "1.2.3", "same"},
		{"v1.2.3", "1.2.4", "minor"},
		{"1.0.0", "2.0.0", "major"},
		{"1.2.3", "1.3.0", "minor"},
		{"", "1.0.0", "unknown"},
		{"1.0.0", "", "unknown"},
		{"not-a-version", "1.0.0", "incomparable"},
		{"1.0.0-alpha", "not-semver", "incomparable"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.want+"_"+tc.current+"_"+tc.observed, func(t *testing.T) {
			t.Parallel()
			got := render.VersionDifference(tc.current, tc.observed)
			if got != tc.want {
				t.Fatalf("VersionDifference(%q, %q) = %q, want %q", tc.current, tc.observed, got, tc.want)
			}
		})
	}
}

func TestStatusEmoji(t *testing.T) {
	t.Parallel()
	cases := []struct {
		state, want string
	}{
		{"Ready", "✅"},
		{"current", "✅"},
		{"warning", "⚠️"},
		{"minor", "⚠️"},
		{"stale", "⚠️"},
		{"NotReady", "🔴"},
		{"major", "🔴"},
		{"unknown", "◽"},
		{"fresh", "◽"},
		{"missing", "◽"},
		{"", "◽"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.state, func(t *testing.T) {
			t.Parallel()
			got := render.StatusEmoji(tc.state)
			if got != tc.want {
				t.Fatalf("StatusEmoji(%q) = %q, want %q", tc.state, got, tc.want)
			}
		})
	}
}

func TestFmtTimeAndDate(t *testing.T) {
	t.Parallel()
	ts := time.Date(2026, 7, 21, 5, 0, 0, 0, time.UTC)
	if got := render.FmtTime(ts); got != "2026-07-21T05:00:00Z" {
		t.Fatalf("FmtTime = %q, want RFC3339 UTC", got)
	}
	if got := render.FmtDate(ts); got != "2026-07-21" {
		t.Fatalf("FmtDate = %q, want YYYY-MM-DD", got)
	}
	// Non-UTC input must still emit UTC.
	local := time.Date(2026, 7, 21, 7, 0, 0, 0, time.FixedZone("CET", 2*3600))
	if got := render.FmtTime(local); got != "2026-07-21T05:00:00Z" {
		t.Fatalf("FmtTime(CET) = %q, want 2026-07-21T05:00:00Z", got)
	}
}

func TestAnchor(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
	}{
		{"region-a/cluster-alpha", "region-a-cluster-alpha"},
		{"Hello World", "hello-world"},
		{"Already_OK-1", "already_ok-1"},
		{"  Spaced  ", "spaced"},
		{"a--b", "a-b"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			got := render.Anchor(tc.in)
			if got != tc.want {
				t.Fatalf("Anchor(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSortByStableCaseInsensitive(t *testing.T) {
	t.Parallel()
	type row struct {
		Name string
		ID   int
	}
	in := []row{{Name: "bravo", ID: 1}, {Name: "Alpha", ID: 2}, {Name: "alpha", ID: 3}, {Name: "Charlie", ID: 4}}
	out, err := render.SortBy("Name", in)
	if err != nil {
		t.Fatalf("SortBy: %v", err)
	}
	got, ok := out.([]row)
	if !ok {
		t.Fatalf("SortBy type = %T, want []row", out)
	}
	wantNames := []string{"Alpha", "alpha", "bravo", "Charlie"}
	wantIDs := []int{2, 3, 1, 4} // stable: Alpha before alpha (original order among equal fold)
	for i := range wantNames {
		if got[i].Name != wantNames[i] || got[i].ID != wantIDs[i] {
			t.Fatalf("SortBy[%d] = {%q,%d}, want {%q,%d}", i, got[i].Name, got[i].ID, wantNames[i], wantIDs[i])
		}
	}
}

func TestGroupByStableCaseInsensitive(t *testing.T) {
	t.Parallel()
	type row struct {
		Namespace string
		Name      string
	}
	in := []row{
		{Namespace: "platform", Name: "b"},
		{Namespace: "Messaging", Name: "a"},
		{Namespace: "platform", Name: "a"},
		{Namespace: "messaging", Name: "b"},
	}
	out, err := render.GroupBy("Namespace", in)
	if err != nil {
		t.Fatalf("GroupBy: %v", err)
	}
	groups, ok := out.([]render.Group)
	if !ok {
		t.Fatalf("GroupBy type = %T, want []render.Group", out)
	}
	if len(groups) != 2 {
		t.Fatalf("groups = %d, want 2", len(groups))
	}
	// Keys sorted case-insensitively; first-seen casing preserved.
	if groups[0].Key != "Messaging" {
		t.Fatalf("groups[0].Key = %q, want Messaging", groups[0].Key)
	}
	if groups[1].Key != "platform" {
		t.Fatalf("groups[1].Key = %q, want platform", groups[1].Key)
	}
}

func TestCapNote(t *testing.T) {
	t.Parallel()
	if got := render.CapNote(true, 3); got != "(list capped at 3 entries)" {
		t.Fatalf("CapNote(true,3) = %q", got)
	}
	if got := render.CapNote(false, 3); got != "" {
		t.Fatalf("CapNote(false,3) = %q, want empty", got)
	}
}

func TestGoldenEnvInventoryMarkdown(t *testing.T) {
	t.Parallel()
	dir := goldenDir(t, "env-inventory-md")
	tmpl := mustRead(t, filepath.Join(dir, "template.md.tmpl"))
	want := mustRead(t, filepath.Join(dir, "expected.md"))
	ctx := mustLoadContext(t, filepath.Join(dir, "context.yaml"))

	got1, err := render.Render(tmpl, ctx)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	got2, err := render.Render(tmpl, ctx)
	if err != nil {
		t.Fatalf("Render second pass: %v", err)
	}
	if string(got1) != string(got2) {
		t.Fatalf("Render not deterministic across two in-process runs")
	}
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(filepath.Join(dir, "expected.md"), got1, 0o644); err != nil {
			t.Fatalf("update golden: %v", err)
		}
		return
	}
	if string(got1) != want {
		t.Fatalf("golden mismatch (-want +got):\n%s", diffSnippet(want, string(got1)))
	}
}

func goldenDir(t *testing.T, name string) string {
	t.Helper()
	root := repoRoot(t)
	return filepath.Join(root, "test", "golden", name)
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
