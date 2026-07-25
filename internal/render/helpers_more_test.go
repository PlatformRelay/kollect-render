package render_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/platformrelay/kollect-render/internal/render"
)

func TestRenderParseAndExecuteErrors(t *testing.T) {
	t.Parallel()
	if _, err := render.Render("{{", render.RenderContext{}); err == nil {
		t.Fatal("want parse error")
	}
	if _, err := render.Render("{{.MissingKey.Nope}}", render.RenderContext{}); err == nil {
		t.Fatal("want execute error for missing key")
	}
}

func TestSemverCompareInvalidFallsBack(t *testing.T) {
	t.Parallel()
	got := render.SemverCompare("not-semver-a", "not-semver-b")
	want := strings.Compare("not-semver-a", "not-semver-b")
	if got != want {
		t.Fatalf("SemverCompare invalid = %d, want %d", got, want)
	}
}

func TestFmtTimeAndDateCoercions(t *testing.T) {
	t.Parallel()
	ts := time.Date(2026, 7, 21, 5, 0, 0, 0, time.UTC)
	if got := render.FmtTime(&ts); got != "2026-07-21T05:00:00Z" {
		t.Fatalf("FmtTime(*time.Time) = %q", got)
	}
	if got := render.FmtDate(&ts); got != "2026-07-21" {
		t.Fatalf("FmtDate(*time.Time) = %q", got)
	}
	var nilTS *time.Time
	if got := render.FmtTime(nilTS); got != "" {
		t.Fatalf("FmtTime(nil *time.Time) = %q, want empty", got)
	}
	if got := render.FmtDate(nilTS); got != "" {
		t.Fatalf("FmtDate(nil *time.Time) = %q, want empty", got)
	}
	if got := render.FmtTime("2026-07-21T05:00:00Z"); got != "2026-07-21T05:00:00Z" {
		t.Fatalf("FmtTime(string) = %q", got)
	}
	if got := render.FmtDate("2026-07-21T05:00:00Z"); got != "2026-07-21" {
		t.Fatalf("FmtDate(string) = %q", got)
	}
	if got := render.FmtTime("not-a-time"); got != "" {
		t.Fatalf("FmtTime(bad string) = %q, want empty", got)
	}
	if got := render.FmtDate(42); got != "" {
		t.Fatalf("FmtDate(int) = %q, want empty", got)
	}
	if got := render.FmtTime(42); got != "" {
		t.Fatalf("FmtTime(int) = %q, want empty", got)
	}
}

func TestAnchorEmpty(t *testing.T) {
	t.Parallel()
	if got := render.Anchor("   "); got != "" {
		t.Fatalf("Anchor(blank) = %q", got)
	}
	if got := render.Anchor("---"); got != "" {
		t.Fatalf("Anchor(dashes) = %q", got)
	}
}

func TestSortByErrors(t *testing.T) {
	t.Parallel()
	if _, err := render.SortBy("Name", "not-a-slice"); err == nil {
		t.Fatal("want error for non-slice")
	}
	type row struct{ Name string }
	if _, err := render.SortBy("Missing", []row{{Name: "a"}}); err == nil {
		t.Fatal("want error for missing field")
	}
	if _, err := render.SortBy("Name", []int{1, 2}); err == nil {
		t.Fatal("want error for non-struct element")
	}
}

func TestGroupByErrors(t *testing.T) {
	t.Parallel()
	if _, err := render.GroupBy("Name", map[string]string{}); err == nil {
		t.Fatal("want error for non-slice")
	}
	type row struct{ Name string }
	if _, err := render.GroupBy("Missing", []row{{Name: "a"}}); err == nil {
		t.Fatal("want error for missing field")
	}
}

func TestSortByNonStringField(t *testing.T) {
	t.Parallel()
	type row struct {
		Name string
		N    int
	}
	out, err := render.SortBy("N", []row{{Name: "b", N: 2}, {Name: "a", N: 1}})
	if err != nil {
		t.Fatal(err)
	}
	got := out.([]row)
	if got[0].N != 1 || got[1].N != 2 {
		t.Fatalf("SortBy int field = %+v", got)
	}
}

func TestSortByNilPointerElement(t *testing.T) {
	t.Parallel()
	type row struct{ Name string }
	var nilRow *row
	in := []*row{nilRow, {Name: "z"}, {Name: "a"}}
	out, err := render.SortBy("Name", in)
	if err != nil {
		t.Fatal(err)
	}
	got := out.([]*row)
	if got[0] != nilRow {
		t.Fatalf("nil pointer should sort as empty key first, got %+v", got[0])
	}
	if got[1].Name != "a" || got[2].Name != "z" {
		t.Fatalf("got = %+v", got)
	}
}

func TestFuncMapContainsHelpers(t *testing.T) {
	t.Parallel()
	fm := render.FuncMap()
	for _, name := range render.HelperNames {
		if _, ok := fm[name]; !ok {
			t.Fatalf("FuncMap missing %q", name)
		}
	}
	_ = reflect.TypeOf(fm)
}
