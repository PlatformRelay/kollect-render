#!/usr/bin/env bash
# Fail if any GitHub Actions workflow uses a mutable (non-SHA) action ref.
# Allowed exceptions: local composites (uses: ./...) and docker:// images.
# Expected pin form: owner/repo@<40-hex> # vN.M.P
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$root"

shopt -s nullglob
workflows=(.github/workflows/*.yaml .github/workflows/*.yml)
if [[ ${#workflows[@]} -eq 0 ]]; then
  echo "pinned-actions: no workflows under .github/workflows/" >&2
  exit 1
fi

fail=0
count=0
for wf in "${workflows[@]}"; do
  lineno=0
  while IFS= read -r line || [[ -n "$line" ]]; do
    lineno=$((lineno + 1))
    # Match "uses:" at YAML indent (skip comments / unrelated text).
    if ! [[ "$line" =~ ^[[:space:]]*-?[[:space:]]*uses:[[:space:]]* ]]; then
      continue
    fi
    ref="${line#*uses:}"
    ref="${ref#"${ref%%[![:space:]]*}"}" # ltrim
    ref="${ref%%[[:space:]]*#*}"         # drop trailing comment
    ref="${ref%"${ref##*[![:space:]]}"}" # rtrim

    # Local composite / reusable workflow — nothing to pin.
    if [[ "$ref" == ./* ]]; then
      continue
    fi
    # Container actions — digest pinning is out of scope for this gate.
    if [[ "$ref" == docker://* ]]; then
      continue
    fi

    count=$((count + 1))
    if [[ ! "$ref" =~ ^[^@[:space:]]+@[0-9a-f]{40}$ ]]; then
      echo "pinned-actions: unpinned action at ${wf}:${lineno}: ${ref}" >&2
      fail=1
    fi
  done <"$wf"
done

if [[ "$count" -eq 0 ]]; then
  echo "pinned-actions: no 'uses:' action refs found — assert is vacuous" >&2
  exit 1
fi

if [[ "$fail" -ne 0 ]]; then
  echo "pinned-actions: FAIL — pin every third-party action to a full 40-char commit SHA" >&2
  exit 1
fi

echo "pinned-actions: ok (${count} refs)"
