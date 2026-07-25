package artifact_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Forbidden imports for the pure render path (REQ-E2-S05-01 network denial).
var forbiddenNetworkPrefixes = []string{
	"net",
	"net/http",
	"net/url",
	"crypto/tls",
	"golang.org/x/net",
}

// Packages that must remain offline (render → artifact for private publisher).
var renderPathDirs = []string{
	".",
	"../render",
	"../format",
}

func TestRenderPathDeniesNetworkImports(t *testing.T) {
	t.Parallel()
	fset := token.NewFileSet()
	for _, dir := range renderPathDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("readdir %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
				continue
			}
			path := filepath.Join(dir, e.Name())
			f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			for _, imp := range f.Imports {
				pathLit := strings.Trim(imp.Path.Value, `"`)
				for _, bad := range forbiddenNetworkPrefixes {
					if pathLit == bad || strings.HasPrefix(pathLit, bad+"/") {
						t.Errorf("%s imports forbidden network package %q", path, pathLit)
					}
				}
			}
		}
	}
}
