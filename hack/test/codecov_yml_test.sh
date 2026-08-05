#!/usr/bin/env bash
# Assert Codecov upload does not rely on coverage-file search (P2 residual).
# disable_search is a codecov-action / CLI input — not a codecov.yml field
# (codecov.io/validate rejects it as unknown). Guard the CI step.
# Run: bash hack/test/codecov_yml_test.sh
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

assert_grep "codecov-action present" 'uses:[[:space:]]*codecov/codecov-action@'
assert_grep "explicit coverage file" 'files:[[:space:]]*coverage\.out'
assert_grep "disable_search true" 'disable_search:[[:space:]]*true'

if [[ "${fail}" -ne 0 ]]; then
  echo "codecov_yml_test: RED" >&2
  exit 1
fi
echo "codecov_yml_test: GREEN"
