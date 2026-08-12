package validate

import (
	"errors"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
)

func TestValidationMessageNil(t *testing.T) {
	t.Parallel()
	if validationMessage(nil) != "validation failed" {
		t.Fatal("nil error")
	}
}

func TestPointeredResultPassesThroughNonValidationError(t *testing.T) {
	t.Parallel()
	in := errors.New("not a schema error")
	if got := pointeredResult(in); !errors.Is(got, in) {
		t.Fatalf("got %v, want passthrough", got)
	}
}

func TestCollectLeavesNilIsNoop(t *testing.T) {
	t.Parallel()
	res := &Result{}
	collectLeaves(nil, res)
	if len(res.Errors) != 0 {
		t.Fatalf("errors = %+v", res.Errors)
	}
}

func TestCollectLeavesNilKindUsesErrorString(t *testing.T) {
	t.Parallel()
	res := &Result{}
	collectLeaves(&jsonschema.ValidationError{ErrorKind: nil}, res)
	if len(res.Errors) != 1 || res.Errors[0].Message != "validation failed" {
		t.Fatalf("errors = %+v", res.Errors)
	}
}

func TestPropertyFromKindEdges(t *testing.T) {
	t.Parallel()
	if _, ok := propertyFromKind(nil); ok {
		t.Fatal("nil kind should not yield a property")
	}
	if _, ok := propertyFromKind(&kind.Required{Missing: []string{"a", "b"}}); ok {
		t.Fatal("multi-missing required should not pick a single property")
	}
	if _, ok := propertyFromKind(&kind.AdditionalProperties{Properties: []string{"a", "b"}}); ok {
		t.Fatal("multi additional properties should not pick a single property")
	}
	if _, ok := propertyFromKind(&kind.Type{Got: "string", Want: []string{"number"}}); ok {
		t.Fatal("type mismatch is not a property name")
	}
	got, ok := propertyFromKind(&kind.Required{Missing: []string{"schemaVersion"}})
	if !ok || got != "schemaVersion" {
		t.Fatalf("required single = %q ok=%v", got, ok)
	}
	got, ok = propertyFromKind(&kind.AdditionalProperties{Properties: []string{"secretLeak"}})
	if !ok || got != "secretLeak" {
		t.Fatalf("additional single = %q ok=%v", got, ok)
	}
}

func TestJSONPointerEscapes(t *testing.T) {
	t.Parallel()
	if jsonPointer(nil) != "" {
		t.Fatal("empty tokens")
	}
	got := jsonPointer([]string{"a/b", "c~d"})
	if got != "/a~1b/c~0d" {
		t.Fatalf("got %q", got)
	}
}

func TestJoinPointerRootAndNested(t *testing.T) {
	t.Parallel()
	if joinPointer("", "x") != "/x" {
		t.Fatal("root")
	}
	if joinPointer("/metadata", "enviromentId") != "/metadata/enviromentId" {
		t.Fatal("nested")
	}
}

func TestPointeredResultFallbackWhenNoLeaves(t *testing.T) {
	t.Parallel()
	// A ValidationError with a cause that is nil is skipped; the parent then
	// has Causes != 0 so collectLeaves adds nothing, and pointeredResult
	// must still emit one error (the empty-leaves fallback).
	err := pointeredResult(&jsonschema.ValidationError{
		Causes: []*jsonschema.ValidationError{nil},
	})
	var res *Result
	if !errors.As(err, &res) || len(res.Errors) != 1 {
		t.Fatalf("got %v", err)
	}
	if res.Errors[0].Message != "validation failed" {
		t.Fatalf("fallback message = %q", res.Errors[0].Message)
	}
}
