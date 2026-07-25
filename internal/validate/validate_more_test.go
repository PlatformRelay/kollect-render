package validate_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/validate"
)

func TestBytesEmptyDocument(t *testing.T) {
	t.Parallel()
	err := validate.Bytes([]byte("   \n"))
	if err == nil {
		t.Fatal("want empty document error")
	}
	var res *validate.Result
	if !errors.As(err, &res) {
		t.Fatalf("want *validate.Result, got %T: %v", err, err)
	}
	if len(res.Errors) == 0 || res.Errors[0].Message != "document is empty" {
		t.Fatalf("errors = %+v", res.Errors)
	}
}

func TestBytesYAMLParseFail(t *testing.T) {
	t.Parallel()
	err := validate.Bytes([]byte("metadata: [\n  - !!broken"))
	if err == nil {
		t.Fatal("want YAML parse error")
	}
	if !strings.Contains(err.Error(), "parse document") {
		t.Fatalf("error = %v, want parse document prefix", err)
	}
}

func TestFileMissing(t *testing.T) {
	t.Parallel()
	err := validate.File(filepath.Join(t.TempDir(), "no-such.yaml"))
	if err == nil {
		t.Fatal("want missing file error")
	}
}

func TestErrorAndResultFormatting(t *testing.T) {
	t.Parallel()
	e := validate.Error{Pointer: "", Message: "bare"}
	if e.Error() != "bare" {
		t.Fatalf("Error() = %q", e.Error())
	}
	e2 := validate.Error{Pointer: "/x", Message: "y"}
	if e2.Error() != "/x: y" {
		t.Fatalf("Error() = %q", e2.Error())
	}
	var nilRes *validate.Result
	if nilRes.Error() != "" {
		t.Fatalf("nil Result.Error = %q", nilRes.Error())
	}
	empty := &validate.Result{}
	if empty.Error() != "" {
		t.Fatalf("empty Result.Error = %q", empty.Error())
	}
	one := &validate.Result{Errors: []validate.Error{{Pointer: "/a", Message: "one"}}}
	if one.Error() != "/a: one" {
		t.Fatalf("single Result.Error = %q", one.Error())
	}
	multi := &validate.Result{Errors: []validate.Error{
		{Pointer: "/a", Message: "one"},
		{Pointer: "/b", Message: "two"},
	}}
	got := multi.Error()
	if !strings.Contains(got, "2 validation errors") || !strings.Contains(got, "/a: one") || !strings.Contains(got, "/b: two") {
		t.Fatalf("multi Result.Error = %q", got)
	}
}
