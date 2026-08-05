#!/usr/bin/env bash
# Sonar githubactions:S8233 — docs workflow: elevated Pages/OIDC perms
# must be job-scoped (deploy only), not workflow-level.
# Run: bash hack/test/docs_perms_s8233_test.sh
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
wf="${root}/.github/workflows/docs.yaml"

fail=0
assert_ok() {
  local desc="$1"
  echo "ok: ${desc}"
}

assert_fail() {
  local desc="$1"
  echo "FAIL: ${desc}" >&2
  fail=1
}

[[ -f "${wf}" ]] || {
  echo "FAIL: missing ${wf}" >&2
  exit 1
}

# Workflow-level permissions block: first ^permissions: until next top-level key.
wf_perms="$(awk '
  /^permissions:[[:space:]]*$/ { in_block=1; next }
  in_block && /^[a-zA-Z]/ { exit }
  in_block { print }
' "${wf}")"

if echo "${wf_perms}" | grep -Eq '^[[:space:]]*contents:[[:space:]]*read[[:space:]]*$'; then
  assert_ok "workflow-level contents: read"
else
  assert_fail "workflow-level must declare contents: read"
fi

if echo "${wf_perms}" | grep -Eq 'pages:'; then
  assert_fail "workflow-level must not declare pages (S8233)"
else
  assert_ok "workflow-level has no pages permission"
fi

if echo "${wf_perms}" | grep -Eq 'id-token:'; then
  assert_fail "workflow-level must not declare id-token (S8233)"
else
  assert_ok "workflow-level has no id-token permission"
fi

# Deploy job keeps elevated perms (required for GitHub Pages).
deploy_block="$(awk '
  /^[[:space:]]{2}deploy:[[:space:]]*$/ { in_job=1; next }
  in_job && /^[[:space:]]{2}[a-zA-Z0-9_-]+:/ { exit }
  in_job { print }
' "${wf}")"

if echo "${deploy_block}" | grep -Eq '^[[:space:]]+pages:[[:space:]]*write[[:space:]]*$'; then
  assert_ok "deploy job pages: write"
else
  assert_fail "deploy job must declare pages: write"
fi

if echo "${deploy_block}" | grep -Eq '^[[:space:]]+id-token:[[:space:]]*write[[:space:]]*$'; then
  assert_ok "deploy job id-token: write"
else
  assert_fail "deploy job must declare id-token: write"
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "docs_perms_s8233_test: RED" >&2
  exit 1
fi
echo "docs_perms_s8233_test: GREEN"
