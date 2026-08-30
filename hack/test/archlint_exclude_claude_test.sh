#!/usr/bin/env bash
# .go-arch-lint.yml must exclude .claude/ so nested agent worktrees
# (created at .claude/worktrees/<id>/) do not produce spurious
# "not attached to any component" notices in `task arch-lint`.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"

if ! grep -qE '^\s*-\s*\.claude(/\*\*)?\s*$' .go-arch-lint.yml; then
  echo "archlint_exclude_claude: .go-arch-lint.yml exclude list is missing '.claude/**'" >&2
  echo "  Without it, a git worktree under .claude/worktrees/ makes arch-lint fail locally." >&2
  exit 1
fi

echo "archlint_exclude_claude: ok"
