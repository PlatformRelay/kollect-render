#!/usr/bin/env bash
# Fail if CHANGELOG.md is stale relative to git history and hack/release/cliff.toml.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT}"

CLIFF="${GIT_CLIFF_BIN:-${ROOT}/bin/git-cliff}"
if [[ ! -x "${CLIFF}" ]]; then
  echo "verify-changelog: installing git-cliff into bin/" >&2
  bash hack/install-git-cliff.sh v2.14.1 bin/git-cliff
  CLIFF="${ROOT}/bin/git-cliff"
fi

CLIFF_CONFIG="${CLIFF_CONFIG:-hack/release/cliff.toml}"

if [[ ! -f CHANGELOG.md ]]; then
  echo "verify-changelog: CHANGELOG.md missing — run 'task changelog:write' and commit" >&2
  exit 1
fi

# A release-prep commit carries "## [X.Y.Z]" before vX.Y.Z is tagged (the release
# gate requires it). Render that section with --tag so it is not seen as drift.
cliff_args=(--config "${CLIFF_CONFIG}")
pending=""
top_version="$(sed -n -E 's/^## \[([0-9][^]]*)\].*/\1/p' CHANGELOG.md | head -n 1)"
if [[ -n "${top_version}" ]] && ! git rev-parse -q --verify "refs/tags/v${top_version}" >/dev/null; then
  pending="${top_version}"
  cliff_args+=(--tag "v${pending}")
  echo "verify-changelog: v${pending} not tagged yet — verifying as pending release" >&2
fi

scratch="$(mktemp)"
trap 'rm -f "${scratch}"' EXIT

"${CLIFF}" "${cliff_args[@]}" -o "${scratch}"

# git-cliff dates an untagged --tag section with the current day; keep the committed
# date for that one heading so a later re-run is not red. Tagged sections stay strict.
if [[ -n "${pending}" ]]; then
  want="$(grep -m1 -F "## [${pending}]" CHANGELOG.md)"
  got="$(grep -m1 -F "## [${pending}]" "${scratch}")"
  if [[ "${want% - *}" == "${got% - *}" && "${want}" != "${got}" ]]; then
    fixed="$(mktemp)"
    trap 'rm -f "${scratch}" "${fixed}"' EXIT
    while IFS= read -r line || [[ -n "${line}" ]]; do
      if [[ "${line}" == "${got}" ]]; then line="${want}"; fi
      printf '%s\n' "${line}"
    done <"${scratch}" >"${fixed}"
    mv "${fixed}" "${scratch}"
  fi
fi

if ! diff -u CHANGELOG.md "${scratch}"; then
  echo "verify-changelog: CHANGELOG.md drift — run 'task changelog:write' and commit" >&2
  exit 1
fi

echo "verify-changelog: ok"
