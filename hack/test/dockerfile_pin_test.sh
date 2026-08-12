#!/usr/bin/env bash
# Fail if Dockerfile FROM lines are not pinned to a sha256 digest.
# Scorecard Pinned-Dependencies treats a floating tag as mutable supply chain.
# Run: bash hack/test/dockerfile_pin_test.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DF="${ROOT}/Dockerfile"

fail=0
ok() { echo "ok: $*"; }
bad() { echo "FAIL: $*" >&2; fail=1; }

[[ -f "${DF}" ]] || { echo "FAIL: Dockerfile missing at ${DF}" >&2; exit 1; }

from_count=0
while IFS= read -r line || [[ -n "${line}" ]]; do
  [[ "${line}" =~ ^[[:space:]]*FROM[[:space:]] ]] || continue
  from_count=$((from_count + 1))
  # Drop trailing comments.
  ref="${line#FROM }"
  ref="${ref%%#*}"
  ref="${ref%"${ref##*[![:space:]]}"}"
  # Stage alias (AS builder) is allowed after the digest.
  if [[ ! "${ref}" =~ @sha256:[0-9a-f]{64}([[:space:]]|$) ]]; then
    bad "unpinned FROM: ${line}"
  else
    ok "pinned FROM: ${ref}"
  fi
done <"${DF}"

if [[ "${from_count}" -eq 0 ]]; then
  echo "FAIL: no FROM lines in Dockerfile — assert is vacuous" >&2
  exit 1
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "dockerfile_pin_test: RED" >&2
  exit 1
fi
echo "dockerfile_pin_test: GREEN"
