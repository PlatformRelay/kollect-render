package validate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/platformrelay/kollect-render/internal/validate"
)

// FuzzBytes exercises the YAML/JSON decode + schema-validate path. The contract
// under test is total: validate.Bytes must return (nil or an error) for ANY
// input and must never panic. Seeds are the shipped example inventories plus
// malformed shapes that historically break YAML/JSON bridges.
func FuzzBytes(f *testing.F) {
	for _, name := range []string{
		"region-a-cluster-alpha.yaml",
		"region-b-cluster-beta.yaml",
	} {
		raw, err := os.ReadFile(filepath.Join(fuzzRepoRoot(f), "schema", "examples", name))
		if err != nil {
			f.Fatalf("read seed %s: %v", name, err)
		}
		f.Add(raw)
	}
	// Non-golden seeds: empty, whitespace, scalar, alias-heavy YAML, deep
	// nesting, NaN/Inf (unrepresentable in JSON), and a huge-int overflow.
	for _, s := range []string{
		"", "   \n\t ", "null", "42", "[]", "{}", "- a\n- b",
		"a: !!binary |\n  aGVsbG8=",
		"x: &a [1]\ny: *a",
		"a: .inf", "b: .nan",
		"n: 100000000000000000000000000000",
		"{\"apiVersion\":\"v0\"}",
	} {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, raw []byte) {
		// Total function: any return is acceptable, a panic is not.
		_ = validate.Bytes(raw)
	})
}

func fuzzRepoRoot(f *testing.F) string {
	f.Helper()
	wd, err := os.Getwd()
	if err != nil {
		f.Fatalf("getcwd: %v", err)
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
	f.Fatal("cannot find schema/inventory-v0.schema.json from test cwd")
	return ""
}
