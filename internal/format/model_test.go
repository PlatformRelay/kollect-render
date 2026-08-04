package format

import (
	"strings"
	"testing"
)

// Marker methods exist only for interface satisfaction; call them so coverage
// counts the statements (type switches never invoke them).
func TestModelMarkerMethods(t *testing.T) {
	t.Parallel()
	blocks := []Block{
		Banner{}, StatusLegend{}, Heading{}, Paragraph{}, BulletList{},
		Table{}, Blank{}, CapNote{}, Footnote{},
	}
	for _, b := range blocks {
		b.block()
	}
	inlines := []Inline{Text{}, Strong{}, Code{}, Link{}, Emoji{}}
	for _, in := range inlines {
		in.inline()
	}
}

func TestEncodeUnsupportedBlock(t *testing.T) {
	t.Parallel()
	bad := Model{Blocks: []Block{unsupportedBlock{}}}
	if _, err := (markdownEncoder{}).Encode(bad); err == nil {
		t.Fatal("markdown: want unsupported block error")
	}
	_, err := (confluenceEncoder{}).Encode(bad)
	if err == nil {
		t.Fatal("confluence: want unsupported block error")
	}
	const want = "confluence-storage: unsupported block"
	if got := err.Error(); !strings.Contains(got, want) {
		t.Fatalf("confluence error = %q, want substring %q", got, want)
	}
}

type unsupportedBlock struct{}

func (unsupportedBlock) block() {}

func TestMustRegisterPanicOnDuplicate(t *testing.T) {
	// Not parallel: mutates the global format registry.
	defer func() {
		if recover() == nil {
			t.Fatal("MustRegister duplicate must panic")
		}
	}()
	MustRegister(markdownEncoder{})
}

func TestNamesIncludesExtrasSorted(t *testing.T) {
	// Not parallel: mutates the global format registry.
	name := "zzz-test-format"
	MustRegister(namedEncoder{name: name})
	t.Cleanup(func() { delete(registry, name) })
	names := Names()
	found := false
	for _, n := range names {
		if n == name {
			found = true
		}
	}
	if !found {
		t.Fatalf("Names() missing %q: %v", name, names)
	}
	// Prefer builtins first.
	if names[0] != NameMarkdown || names[1] != NameConfluenceStorage {
		t.Fatalf("Names() = %v, want builtins first", names)
	}
}

type namedEncoder struct{ name string }

func (e namedEncoder) Name() string                 { return e.name }
func (e namedEncoder) Encode(Model) ([]byte, error) { return nil, nil }
