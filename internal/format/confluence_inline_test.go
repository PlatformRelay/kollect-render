package format

import (
	"strings"
	"testing"
)

type rogueInline struct{}

func (rogueInline) inline() {}

func TestWriteConfluenceInlinesFailClosed(t *testing.T) {
	t.Parallel()
	var b strings.Builder
	err := writeConfluenceInlines(&b, []Inline{rogueInline{}})
	if err == nil {
		t.Fatal("expected error for unknown Inline")
	}
	if !strings.Contains(err.Error(), "unknown inline") {
		t.Fatalf("error = %v, want unknown inline", err)
	}
}
