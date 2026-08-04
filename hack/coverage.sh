#!/usr/bin/env bash
# Run package tests and enforce a minimum statement-coverage floor on the whole module.
# Override floor with COVERAGE_MIN (default 90).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

MIN="${COVERAGE_MIN:-90}"
export CGO_ENABLED="${CGO_ENABLED:-0}"

rm -f coverage.out coverage-summary.txt

go test -count=1 -coverpkg=./... -coverprofile=coverage.out ./...

if [[ ! -s coverage.out ]] || [[ "$(wc -l < coverage.out)" -lt 2 ]]; then
  echo "coverage.out missing or empty after go test ./..." >&2
  exit 1
fi

go tool cover -func=coverage.out | tee coverage-summary.txt | tail -1
pct="$(awk '/^total:/ {gsub(/%/,"",$3); print $3}' coverage-summary.txt)"
awk -v p="${pct}" -v m="${MIN}" 'BEGIN {
  if (p+0 < m+0) {
    printf "module coverage %.1f%% is below %d%% floor\n", p, m
    exit 1
  }
}'
echo "module coverage ${pct}% (floor ${MIN}%)"
