#!/usr/bin/env bash
# REWE-trace gate: fail if forbidden internal markers appear in tracked files.
# Patterns are assembled at runtime so this script does not contain the
# forbidden substrings as contiguous literals (the gate scans itself too).
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$root"

# shellcheck disable=SC2034
patterns=(
  "esc""-cli"
  "TITI""SEVEN"
  "prod""165"
  "rewe"".""cloud"
  "rewe"" digital"
)

# Case-insensitive fixed-string search across all tracked files.
fail=0
while IFS= read -r -d '' file; do
  for p in "${patterns[@]}"; do
    if grep -Fqi -- "$p" "$file"; then
      echo "REWE-trace: forbidden marker '${p}' found in ${file}" >&2
      fail=1
    fi
  done
done < <(git ls-files -z)

if [[ "$fail" -ne 0 ]]; then
  echo "REWE-trace: FAIL — remove internal markers from tracked files" >&2
  exit 1
fi

echo "REWE-trace: ok"
