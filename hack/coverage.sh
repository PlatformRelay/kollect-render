#!/usr/bin/env bash
# Run package tests and enforce a minimum statement-coverage floor on ./internal/...
# Override floor with COVERAGE_MIN (default 90).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

MIN="${COVERAGE_MIN:-90}"
export CGO_ENABLED="${CGO_ENABLED:-0}"

rm -f coverage.out coverage-summary.txt

# Non-internal packages (cmd, schema) still run — no coverage floor on them.
other_pkgs="$(go list ./... | grep -v '/internal/' || true)"
if [[ -n "${other_pkgs}" ]]; then
  # shellcheck disable=SC2086
  go test -count=1 ${other_pkgs}
fi

go test -count=1 -coverpkg=./internal/... -coverprofile=coverage.out ./internal/...

if [[ ! -s coverage.out ]] || [[ "$(wc -l < coverage.out)" -lt 2 ]]; then
  echo "coverage.out missing or empty after go test ./internal/..." >&2
  exit 1
fi

go tool cover -func=coverage.out | tee coverage-summary.txt | tail -1
pct="$(awk '/^total:/ {gsub(/%/,"",$3); print $3}' coverage-summary.txt)"
awk -v p="${pct}" -v m="${MIN}" 'BEGIN {
  if (p+0 < m+0) {
    printf "internal/ coverage %.1f%% is below %d%% floor\n", p, m
    exit 1
  }
}'
echo "internal/ coverage ${pct}% (floor ${MIN}%)"
