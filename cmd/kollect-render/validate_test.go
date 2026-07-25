package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunValidateShippedExample(t *testing.T) {
	t.Parallel()
	path := filepath.Join(repoRoot(t), "schema", "examples", "region-a-cluster-alpha.yaml")
	if code := run([]string{"validate", path}); code != 0 {
		t.Fatalf("validate exit = %d, want 0", code)
	}
}

func TestRunValidateSeededViolation(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	content := "schemaVersion: \"9.9.9\"\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	if code := run([]string{"validate", path}); code != 2 {
		t.Fatalf("validate exit = %d, want 2", code)
	}
}

func TestRunValidateUsage(t *testing.T) {
	t.Parallel()
	if code := run([]string{"validate"}); code != 2 {
		t.Fatalf("validate usage exit = %d, want 2", code)
	}
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
	t.Fatal("cannot find schema/inventory-v0.schema.json")
	return ""
}
