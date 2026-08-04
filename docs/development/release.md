# Releasing kollect-render

Maintainers publish versions by tagging a **proven** commit on protected `main`.
The publishing workflow (`.github/workflows/release.yaml`) is tag-driven and unchanged
by the pre-tag gate — eligibility is a separate, read-only check you run **before**
creating `vX.Y.Z`.

## Pre-tag release gate

[`.github/workflows/release-gate.yaml`](https://github.com/PlatformRelay/kollect-render/blob/main/.github/workflows/release-gate.yaml)
is `workflow_dispatch` only. Permissions are read-only (`contents`, `checks`,
`pull-requests`) — the job never creates a tag, never logs into a registry, and never
uploads release assets.

Given a full 40-character lowercase commit SHA (or empty to use current `main` HEAD)
and a SemVer-like `version` **without** a leading `v`, the gate proves:

1. The SHA is reachable from protected `main` (ancestor or identical).
2. Required exact-SHA checks on that commit are `completed/success` (`check`, `changelog`).
3. `CHANGELOG.md` **at that SHA** contains a Keep a Changelog heading for the version
   (`## [X.Y.Z]` …).

Short SHAs, commits not on `main`, red/missing checks, and a missing changelog section
each fail with a message naming the failed precondition.

Verifier scripts always load from the **default branch** sparse checkout so a candidate
commit cannot supply its own gate.

### Run locally (fixture tests)

```bash
bash hack/release/test-verify-eligibility.sh
```

These tests mock `gh` — no live GitHub API and no credentials beyond what you already
use for local development.

### Run the gate on GitHub

```bash
RELEASE_SHA="$(git rev-parse HEAD)"   # must be full 40-hex, on main
VERSION="0.2.0"                       # must match CHANGELOG.md section

gh workflow run release-gate.yaml \
  -f sha="${RELEASE_SHA}" \
  -f version="${VERSION}"
gh run list --workflow release-gate.yaml --limit 1
```

Omit `sha` to evaluate current protected `main` HEAD:

```bash
gh workflow run release-gate.yaml -f version="${VERSION}"
```

When the gate is green, tag **only that SHA** and push:

```bash
git tag "v${VERSION}" "${RELEASE_SHA}"
git push origin "v${VERSION}"
```

That push triggers Release (GoReleaser). The gate itself does not publish.

## Checklist before tagging

1. `main` is green for the candidate SHA (`check` + `changelog`).
2. `CHANGELOG.md` on that SHA already has `## [X.Y.Z]` (run `task changelog:write` /
   commit as needed).
3. Release gate passed for `${RELEASE_SHA}` + `${VERSION}`.
4. Tag `vX.Y.Z` at that exact SHA — never an arbitrary branch tip.
