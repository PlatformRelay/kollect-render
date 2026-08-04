#!/usr/bin/env bash
# E3-S07: Renovate config, workflow, and annotated version pins.
# Run: bash hack/test/renovate_e3_s07_test.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RENOVATE_JSON="${ROOT}/renovate.json"
RENOVATE_CFG="${ROOT}/.github/renovate-config.json"
RENOVATE_WF="${ROOT}/.github/workflows/renovate.yaml"
TASKFILE="${ROOT}/Taskfile.yml"
CI_WF="${ROOT}/.github/workflows/ci.yaml"

fail=0

ok() { echo "ok: $*"; }
bad() { echo "FAIL: $*" >&2; fail=1; }

assert_file() {
  [[ -f "$1" ]] && ok "exists $(basename "$1")" || bad "missing $1"
}

assert_grep() {
  local file="$1" desc="$2" pattern="$3"
  if ! grep -Eq "$pattern" "$file"; then
    bad "${desc} in $(basename "$file") (pattern: ${pattern})"
  else
    ok "${desc}"
  fi
}

assert_annotation_above() {
  local file="$1" var="$2" dep="$3"
  # Require a renovate annotation comment on the line immediately above the var pin.
  if ! awk -v var="$var" -v dep="$dep" '
    $0 ~ ("# renovate: datasource=") && $0 ~ ("depName=" dep) { pending=1; next }
    pending && $0 ~ ("^[[:space:]]*" var ":") { found=1; exit }
    { pending=0 }
    END { exit found ? 0 : 1 }
  ' "$file"; then
    bad "annotation for ${var} (depName=${dep}) missing or not immediately above pin in $(basename "$file")"
  else
    ok "annotation ${var} → ${dep} in $(basename "$file")"
  fi
}

assert_file "${RENOVATE_JSON}"
assert_file "${RENOVATE_CFG}"
assert_file "${RENOVATE_WF}"

# renovate.json essentials
assert_grep "${RENOVATE_JSON}" "config:recommended" '"config:recommended"'
assert_grep "${RENOVATE_JSON}" ":disableDependencyDashboard" '":disableDependencyDashboard"'
assert_grep "${RENOVATE_JSON}" "helpers:pinGitHubActionDigests" '"helpers:pinGitHubActionDigests"'
assert_grep "${RENOVATE_JSON}" "commitMessagePrefix" '"commitMessagePrefix"[[:space:]]*:[[:space:]]*":arrow_up:"'
assert_grep "${RENOVATE_JSON}" "commitMessageAction" '"commitMessageAction"[[:space:]]*:[[:space:]]*"chore\(deps\): update"'
assert_grep "${RENOVATE_JSON}" "ignorePaths bin" '"bin/\*\*"'
assert_grep "${RENOVATE_JSON}" "ignorePaths dist" '"dist/\*\*"'
assert_grep "${RENOVATE_JSON}" "ignorePaths .task" '"\.task/\*\*"'
assert_grep "${RENOVATE_JSON}" "customManagers regex" '"customType"[[:space:]]*:[[:space:]]*"regex"'
assert_grep "${RENOVATE_JSON}" "Taskfile manager pattern" 'Taskfile\\\\\.yml'
assert_grep "${RENOVATE_JSON}" "ci.yaml manager pattern" 'ci\\\\\.yaml'

# Annotations — three Taskfile pins + GITLEAKS in ci.yaml
assert_annotation_above "${TASKFILE}" "GO_ARCH_LINT_VERSION" "fe3dback/go-arch-lint"
assert_annotation_above "${TASKFILE}" "GOLANGCI_LINT_VERSION" "golangci/golangci-lint"
assert_annotation_above "${TASKFILE}" "GITLEAKS_VERSION" "gitleaks/gitleaks"
assert_annotation_above "${CI_WF}" "GITLEAKS_VERSION" "gitleaks/gitleaks"

# Both GITLEAKS pins must share the same version literal
task_ver="$(grep -E '^[[:space:]]*GITLEAKS_VERSION:' "${TASKFILE}" | head -1 | sed -E 's/.*GITLEAKS_VERSION:[[:space:]]*"?([^"]+)"?.*/\1/')"
ci_ver="$(grep -E '^[[:space:]]*GITLEAKS_VERSION:' "${CI_WF}" | head -1 | sed -E 's/.*GITLEAKS_VERSION:[[:space:]]*"?([^"]+)"?.*/\1/')"
if [[ -z "${task_ver}" || -z "${ci_ver}" ]]; then
  bad "could not read GITLEAKS_VERSION from Taskfile/ci.yaml"
elif [[ "${task_ver}" != "${ci_ver}" ]]; then
  bad "GITLEAKS_VERSION diverged: Taskfile=${task_ver} ci.yaml=${ci_ver}"
else
  ok "GITLEAKS_VERSION aligned (${task_ver})"
fi

# Workflow: schedule + dispatch, SHA-pinned actions, RENOVATE_TOKEN warning fallback
assert_grep "${RENOVATE_WF}" "weekly schedule" 'cron:'
assert_grep "${RENOVATE_WF}" "workflow_dispatch" 'workflow_dispatch:'
assert_grep "${RENOVATE_WF}" "RENOVATE_TOKEN fallback" 'secrets\.RENOVATE_TOKEN'
assert_grep "${RENOVATE_WF}" "github.token fallback" 'github\.token'
assert_grep "${RENOVATE_WF}" "::warning:: when token unset" '::warning::'
assert_grep "${RENOVATE_WF}" "configurationFile" 'configurationFile:[[:space:]]*\.github/renovate-config\.json'

# Every uses: in renovate.yaml must be SHA-pinned (delegate to gate script after files exist)
if [[ -f "${RENOVATE_WF}" ]]; then
  while IFS= read -r line; do
    ref="${line#*uses:}"
    ref="${ref#"${ref%%[![:space:]]*}"}"
    ref="${ref%%[[:space:]]*#*}"
    ref="${ref%"${ref##*[![:space:]]}"}"
    if [[ ! "$ref" =~ ^[^@[:space:]]+@[0-9a-f]{40}$ ]]; then
      bad "unpinned action in renovate.yaml: ${ref}"
    fi
  done < <(grep -E '^[[:space:]]*-?[[:space:]]*uses:' "${RENOVATE_WF}" || true)
  ok "renovate.yaml action refs are SHA-pinned (or none yet)"
fi

# Must NOT enable Dependabot version updates
if [[ -f "${ROOT}/.github/dependabot.yml" ]] || [[ -f "${ROOT}/.github/dependabot.yaml" ]]; then
  if grep -Eq 'package-ecosystem' "${ROOT}/.github/dependabot."* 2>/dev/null; then
    bad "Dependabot version-updates config present — E3-S07 keeps Dependabot security-only"
  fi
else
  ok "no Dependabot version-updates config"
fi

if [[ "${fail}" -ne 0 ]]; then
  echo "renovate_e3_s07_test: RED" >&2
  exit 1
fi
echo "renovate_e3_s07_test: GREEN"
