package main

import "testing"

func TestRunVersion(t *testing.T) {
	t.Parallel()
	if code := run([]string{"version"}); code != 0 {
		t.Fatalf("version exit code = %d, want 0", code)
	}
}

func TestRunUnknown(t *testing.T) {
	t.Parallel()
	if code := run([]string{"nope"}); code != 2 {
		t.Fatalf("unknown exit code = %d, want 2", code)
	}
}

func TestRunUsage(t *testing.T) {
	t.Parallel()
	if code := run(nil); code != 2 {
		t.Fatalf("usage exit code = %d, want 2", code)
	}
}
