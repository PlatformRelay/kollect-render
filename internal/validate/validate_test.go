package validate_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/platformrelay/kollect-render/internal/validate"
)

func TestShippedExamplesValidate(t *testing.T) {
	t.Parallel()
	for _, name := range []string{
		"region-a-cluster-alpha.yaml",
		"region-b-cluster-beta.yaml",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			raw, err := os.ReadFile(examplePath(t, name))
			if err != nil {
				t.Fatalf("read example: %v", err)
			}
			if err := validate.Bytes(raw); err != nil {
				t.Fatalf("validate.Bytes: %v", err)
			}
			if err := validate.File(examplePath(t, name)); err != nil {
				t.Fatalf("validate.File: %v", err)
			}
		})
	}
}

func TestSeededViolationsPointered(t *testing.T) {
	t.Parallel()
	base := mustRead(t, examplePath(t, "region-a-cluster-alpha.yaml"))

	cases := []struct {
		name       string
		mutate     func(string) string
		wantPtrSub string
	}{
		{
			name: "missing-required-schemaVersion",
			mutate: func(s string) string {
				return strings.Replace(s, "schemaVersion: \"0.2.0\"\n", "", 1)
			},
			wantPtrSub: "/schemaVersion",
		},
		{
			name: "unsupported-schemaVersion",
			mutate: func(s string) string {
				return strings.Replace(s, "schemaVersion: \"0.2.0\"", "schemaVersion: \"1.0.0\"", 1)
			},
			wantPtrSub: "/schemaVersion",
		},
		{
			name: "wrong-type-nodes-count",
			mutate: func(s string) string {
				return strings.Replace(s, "count: 3", "count: \"three\"", 1)
			},
			wantPtrSub: "/nodes/count",
		},
		{
			name: "unknown-field-outside-extensions",
			mutate: func(s string) string {
				return s + "\nsecretLeak: nope\n"
			},
			wantPtrSub: "/secretLeak",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validate.Bytes([]byte(tc.mutate(base)))
			if err == nil {
				t.Fatal("validate.Bytes: want error, got nil")
			}
			assertPointered(t, err, tc.wantPtrSub)
		})
	}
}

func TestExtensionsUnknownKeysTolerated(t *testing.T) {
	t.Parallel()
	base := mustRead(t, examplePath(t, "region-a-cluster-alpha.yaml"))
	doc := base + "\nextensions:\n  futureOptional:\n    note: ignored-by-readers\n"
	if err := validate.Bytes([]byte(doc)); err != nil {
		t.Fatalf("extensions tolerance: %v", err)
	}
}

func TestTypoOutsideExtensionsRefused(t *testing.T) {
	t.Parallel()
	base := mustRead(t, examplePath(t, "region-a-cluster-alpha.yaml"))
	doc := strings.Replace(base, "metadata:\n", "metadata:\n  enviromentId: typo\n", 1)
	err := validate.Bytes([]byte(doc))
	if err == nil {
		t.Fatal("want refusal of typo field outside extensions")
	}
	assertPointered(t, err, "/metadata/enviromentId")
}

func assertPointered(t *testing.T, err error, wantPtrSub string) {
	t.Helper()
	var res validate.Result
	if !errors.As(err, &res) {
		t.Fatalf("want validate.Result with JSON pointers, got %T: %v", err, err)
	}
	if len(res.Errors) == 0 {
		t.Fatal("validate.Result has no Errors")
	}
	joined := ""
	for _, e := range res.Errors {
		joined += e.Pointer + " " + e.Error() + "\n"
		if e.Pointer == "" {
			t.Errorf("error missing pointer: %+v", e)
		}
	}
	if !strings.Contains(joined, wantPtrSub) {
		t.Fatalf("want pointer containing %q, got:\n%s", wantPtrSub, joined)
	}
}

func examplePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "schema", "examples", name)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	dir := wd
	for i := 0; i < 8; i++ {
		candidate := filepath.Join(dir, "schema", "inventory-v0.schema.json")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("cannot find schema/inventory-v0.schema.json from test cwd")
	return ""
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
