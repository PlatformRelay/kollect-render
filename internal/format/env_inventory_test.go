package format

import (
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/render"
)

func TestSortedAs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "happy sorts by key",
			run: func(t *testing.T) {
				t.Helper()
				in := []render.Component{
					{ComponentID: "b", Name: "B"},
					{ComponentID: "a", Name: "A"},
				}
				got, err := sortedAs[render.Component]("ComponentID", in)
				if err != nil {
					t.Fatalf("sortedAs: %v", err)
				}
				if len(got) != 2 || got[0].ComponentID != "a" || got[1].ComponentID != "b" {
					t.Fatalf("sortedAs = %#v, want ComponentID order a,b", got)
				}
			},
		},
		{
			name: "wrong type names expected and actual",
			run: func(t *testing.T) {
				t.Helper()
				in := []render.HelmRelease{{Name: "rel"}}
				got, err := sortedAs[render.Component]("Name", in)
				if err == nil {
					t.Fatalf("sortedAs = %#v, want error", got)
				}
				msg := err.Error()
				if !strings.Contains(msg, "[]render.HelmRelease") {
					t.Fatalf("error %q missing actual type", msg)
				}
				if !strings.Contains(msg, "[]render.Component") {
					t.Fatalf("error %q missing expected type", msg)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}
