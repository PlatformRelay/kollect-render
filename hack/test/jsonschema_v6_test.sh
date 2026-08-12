#!/usr/bin/env bash
# E6-S03: go.mod must require jsonschema/v6 only (no leftover v5).
# Run: bash hack/test/jsonschema_v6_test.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GOMOD="${ROOT}/go.mod"

fail=0
ok() { echo "ok: $*"; }
bad() { echo "FAIL: $*" >&2; fail=1; }

[[ -f "${GOMOD}" ]] || { echo "FAIL: go.mod missing" >&2; exit 1; }

if grep -Eq 'github.com/santhosh-tekuri/jsonschema/v6 v6\.' "${GOMOD}"; then
  ok "go.mod requires jsonschema/v6"
else
  bad "go.mod does not require github.com/santhosh-tekuri/jsonschema/v6"
fi

if grep -Eq 'github.com/santhosh-tekuri/jsonschema/v5' "${GOMOD}"; then
  bad "go.mod still mentions jsonschema/v5 — dual-require is the broken Renovate shape"
else
  ok "go.mod has no jsonschema/v5"
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "jsonschema_v6_test: RED" >&2
  exit 1
fi
echo "jsonschema_v6_test: GREEN"
