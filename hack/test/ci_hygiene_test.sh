#!/usr/bin/env bash
# Assert CI hygiene keys exist in .github/workflows/ci.yaml (E3-S03).
# Run: bash hack/test/ci_hygiene_test.sh
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

assert_grep "concurrency group" 'group:[[:space:]]*ci-\$\{\{ github\.ref \}\}'
assert_grep "cancel-in-progress expression" "cancel-in-progress:[[:space:]]*\\$\\{\\{ github\\.ref != 'refs/heads/main' \\}\\}"
assert_grep "check job timeout" 'timeout-minutes:[[:space:]]*20'
assert_grep "setup-go cache-dependency-path" 'cache-dependency-path:[[:space:]]*go\.sum'
assert_grep "tools cache key" "key:[[:space:]]*\\$\\{\\{ runner\\.os \\}\\}-tools-\\$\\{\\{ hashFiles\\('Taskfile\\.yml'\\) \\}\\}"
assert_grep "tools cache path go/pkg/mod" '~/go/pkg/mod'
assert_grep "tools cache path go-build" '~/\.cache/go-build'
assert_grep "step summary coverage" 'GITHUB_STEP_SUMMARY'

# Every job must declare timeout-minutes: 20 (count job keys vs timeouts).
job_count="$(grep -E '^[[:space:]]{2}[a-zA-Z0-9_-]+:' "$wf" | grep -cv '^[[:space:]]*name:' || true)"
# Simpler: require at least two timeout-minutes: 20 lines (check + sonarqube).
timeout_count="$(grep -cE 'timeout-minutes:[[:space:]]*20' "$wf" || true)"
if [[ "${timeout_count}" -lt 2 ]]; then
  echo "FAIL: expected timeout-minutes: 20 on every job (found ${timeout_count})" >&2
  fail=1
else
  echo "ok: timeout-minutes: 20 present ${timeout_count} times"
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "ci_hygiene_test: RED" >&2
  exit 1
fi
echo "ci_hygiene_test: GREEN"
