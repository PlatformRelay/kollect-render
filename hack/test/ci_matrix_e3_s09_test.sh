#!/usr/bin/env bash
# E3-S09: macOS + Linux test matrix and Windows compile-check on Linux leg.
# Run: bash hack/test/ci_matrix_e3_s09_test.sh
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
wf="${root}/.github/workflows/ci.yaml"

fail=0
assert_grep() {
  local desc="$1" pattern="$2"
  if ! grep -Eq "$pattern" "$wf"; then
    echo "FAIL: ${desc} (pattern: ${pattern})" >&2
    fail=1
  else
    echo "ok: ${desc}"
  fi
}

assert_grep "test job key" '^[[:space:]]{2}test:'
assert_grep "fail-fast false" 'fail-fast:[[:space:]]*false'
assert_grep "ubuntu-latest in matrix" 'ubuntu-latest'
assert_grep "macos-latest in matrix" 'macos-latest'
assert_grep "go test ./..." 'go test \./\.\.\.'
assert_grep "GOOS=windows go vet" 'GOOS=windows[[:space:]]+go[[:space:]]+vet[[:space:]]+\./\.\.\.'
assert_grep "untested Windows binaries comment" '[Uu]ntested [Ww]indows'
# Coverage floor must remain on the Linux check job (COVERAGE_MIN), not the matrix test job.
assert_grep "coverage floor on check" 'COVERAGE_MIN:[[:space:]]*"90"'

if [[ "${fail}" -ne 0 ]]; then
  echo "ci_matrix_e3_s09_test: RED" >&2
  exit 1
fi
echo "ci_matrix_e3_s09_test: GREEN"
