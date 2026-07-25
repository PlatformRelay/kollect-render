package format_test

import (
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/format"
)

func TestMarkdownEncodesAllBlockKinds(t *testing.T) {
	t.Parallel()
	model := format.Model{Blocks: []format.Block{
		format.Heading{Level: 2, Text: "plain"},
		format.Heading{Level: 2, Text: "anchored", Anchor: "a1"},
		format.Banner{
			GeneratedAt: "2026-07-21T05:00:00Z", GeneratedDate: "2026-07-21",
			Origin: "manual", SourceRepoURL: "https://example.org/r", SnapshotSHA: "abc",
			HowToChange: "edit upstream",
		},
		format.StatusLegend{},
		format.Paragraph{Inlines: []format.Inline{
			format.Text{S: "t"}, format.Strong{S: "s"}, format.Code{S: "c"},
			format.Link{Text: "l", URL: "https://example.org"}, format.Emoji{S: "✅"},
		}},
		format.BulletList{Items: [][]format.Inline{{format.Text{S: "item"}}}},
		format.Table{Headers: []string{"H"}, Rows: [][]string{{"R"}}},
		format.Blank{},
		format.CapNote{N: 5},
		format.Footnote{Text: "note"},
	}}
	enc, ok := format.Lookup(format.NameMarkdown)
	if !ok {
		t.Fatal("markdown missing")
	}
	out, err := enc.Encode(model)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{"## plain\n", "{#a1}", "**Status legend:**", "| H |", "(list capped at 5 entries)", "_note_"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in:\n%s", want, s)
		}
	}
}

func TestConfluenceEncodesAllBlockKinds(t *testing.T) {
	t.Parallel()
	model := format.Model{Blocks: []format.Block{
		format.Heading{Level: 2, Text: "plain"},
		format.Heading{Level: 2, Text: "anchored", Anchor: "a1"},
		format.Banner{
			GeneratedAt: "2026-07-21T05:00:00Z", GeneratedDate: "2026-07-21",
			Origin: "manual", SourceRepoURL: "https://example.org/r", SnapshotSHA: "abc",
			HowToChange: "edit upstream",
		},
		format.StatusLegend{},
		format.Paragraph{Inlines: []format.Inline{
			format.Text{S: "t"}, format.Strong{S: "s"}, format.Code{S: "c"},
			format.Link{Text: "l", URL: "https://example.org"}, format.Emoji{S: "✅"},
		}},
		format.BulletList{Items: [][]format.Inline{{format.Text{S: "item"}}}},
		format.Table{Headers: []string{"H"}, Rows: [][]string{{"R"}}},
		format.Blank{},
		format.CapNote{N: 5},
		format.Footnote{Text: "note"},
	}}
	enc, ok := format.Lookup(format.NameConfluenceStorage)
	if !ok {
		t.Fatal("confluence-storage missing")
	}
	out, err := enc.Encode(model)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{"<h2>", `ac:name="anchor"`, `ac:name="panel"`, "<ul>", "<table", "list capped at 5", "<em>note</em>"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q in:\n%s", want, s)
		}
	}
	if err := format.ValidateStorageFragment(s); err != nil {
		t.Fatalf("fragment: %v", err)
	}
}
