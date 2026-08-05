#!/usr/bin/env bash
# Assert hack/install-git-cliff.sh pins curl to HTTPS (Sonar shell:S6506).
# Redirects must not be able to downgrade the protocol.
# Run: bash hack/test/install_git_cliff_https_test.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SCRIPT="${ROOT}/hack/install-git-cliff.sh"

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

pass() {
  echo "ok - $*"
}

[[ -f "${SCRIPT}" ]] || fail "install-git-cliff.sh not found at ${SCRIPT}"

curl_lines="$(grep -E 'curl[[:space:]]' "${SCRIPT}" || true)"
[[ -n "${curl_lines}" ]] || fail "no curl invocation in ${SCRIPT}"

while IFS= read -r line; do
  [[ -n "${line}" ]] || continue
  if ! printf '%s\n' "${line}" | grep -qE -- "--proto[= ]['\"]?=https"; then
    fail "curl missing --proto '=https': ${line}"
  fi
  if ! printf '%s\n' "${line}" | grep -qE -- '--tlsv1\.2'; then
    fail "curl missing --tlsv1.2: ${line}"
  fi
done <<<"${curl_lines}"

pass "install-git-cliff.sh curl enforces HTTPS (--proto '=https' --tlsv1.2)"
echo "install_git_cliff_https_test: GREEN"
