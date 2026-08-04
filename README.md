<p align="center">
<a href="https://github.com/platformrelay/kollect-render/actions/workflows/ci.yaml"><img src="https://github.com/platformrelay/kollect-render/actions/workflows/ci.yaml/badge.svg" alt="CI"></a>
<a href="https://github.com/platformrelay/kollect-render/actions/workflows/codeql.yaml"><img src="https://github.com/platformrelay/kollect-render/actions/workflows/codeql.yaml/badge.svg" alt="CodeQL"></a>
<a href="https://securityscorecards.dev/viewer/?uri=github.com/PlatformRelay/kollect-render"><img src="https://api.securityscorecards.dev/projects/github.com/PlatformRelay/kollect-render/badge" alt="OpenSSF Scorecard"></a>
<a href="https://github.com/platformrelay/kollect-render/actions/workflows/docs.yaml"><img src="https://github.com/platformrelay/kollect-render/actions/workflows/docs.yaml/badge.svg" alt="Docs CI"></a>
<a href="https://platformrelay.github.io/kollect-render/"><img src="https://img.shields.io/badge/documentation-GitHub%20Pages-2ea44f?logo=readthedocs&logoColor=white" alt="Documentation"></a>
<a href="https://github.com/platformrelay/kollect-render/blob/main/LICENSE"><img src="https://img.shields.io/github/license/platformrelay/kollect-render" alt="License: MIT"></a>
<a href="https://github.com/platformrelay/kollect-render/releases"><img src="https://img.shields.io/github/v/release/platformrelay/kollect-render" alt="Release"></a>
<a href="https://codecov.io/gh/platformrelay/kollect-render"><img src="https://codecov.io/gh/platformrelay/kollect-render/graph/badge.svg" alt="codecov"></a>
<a href="https://github.com/platformrelay/kollect-render/blob/main/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/platformrelay/kollect-render" alt="Go"></a>
<a href="https://github.com/platformrelay/kollect-render/pkgs/container/kollect-render"><img src="https://img.shields.io/badge/ghcr.io-platformrelay%2Fkollect--render-2496ED?logo=docker&logoColor=white" alt="Container"></a>
</p>

# kollect-render

**Credential-free inventory renderer.** Validate versioned inventory documents and render them
to deterministic artifacts (`markdown`, `confluence-storage`, …). File outputs include a
`.meta.json` digest/metadata sidecar for private publishers. This tool never opens network
connections or holds publish credentials.

Collection, scheduling, credentials, and publishing stay in private callers; this repository is
the public Open Source renderer for the PlatformRelay inventory protocol.

**Docs:** **[platformrelay.github.io/kollect-render](https://platformrelay.github.io/kollect-render/)** —
getting started, CLI reference, formats, schema, and development. This README is the front door;
the site is the map.

## Quick start

```bash
task check          # fmt, vet, lint, arch-lint, coverage≥90%, gitleaks, build, REWE-trace, pinned-actions
go run ./cmd/kollect-render version
# Built-in Model → encoder (any registered format):
go run ./cmd/kollect-render render --format markdown --context test/golden/env-inventory-md/context.yaml
go run ./cmd/kollect-render render --format confluence-storage --context test/golden/env-inventory-md/context.yaml
# Custom markdown template ( --format must be markdown ):
go run ./cmd/kollect-render render --format markdown --template test/golden/env-inventory-md/template.md.tmpl \
  --context test/golden/env-inventory-md/context.yaml
```

Requires Go (see `go.mod`), [Task](https://taskfile.dev/), and [gitleaks](https://github.com/gitleaks/gitleaks)
on `PATH` for the secret-scan step of `task check`.

## Install

Release artifacts (binaries + container) are published on each `vX.Y.Z` tag:

```bash
# Binary (example): download from GitHub Releases
# Container:
docker pull ghcr.io/platformrelay/kollect-render:v0.1.0
```

Pin the image or release asset by version; Renovate can track the GitHub Releases / GHCR datasources.

## Documentation

- Site: [platformrelay.github.io/kollect-render](https://platformrelay.github.io/kollect-render/)
- [Getting started](docs/getting-started.md) · [CLI reference](docs/cli-reference.md) ·
  [Formats](docs/formats.md) · [Schema](docs/schema.md) ·
  [Development](docs/development.md) · [Security](docs/security.md)

Golden fixtures live under `test/golden/`. Regenerate after an intentional render change with
`UPDATE_GOLDEN=1 go test ./internal/format/ ./internal/render/`, then review the diff like code.

## Status

Gate-verified facts (not aspirational):

- Local gate `task check`: fmt, vet, lint, arch-lint, module coverage ≥90%, gitleaks, build,
  REWE-trace, pinned Actions.
- CI runs the same gate; Codecov publishes coverage; CodeQL and OpenSSF Scorecard are published.
- Docs build and deploy to GitHub Pages from `.github/workflows/docs.yaml`.
- Tagged releases (`vX.Y.Z`) publish multi-platform binaries and a GHCR image via GoReleaser.
- `validate` and `render` are available; formats register at compile time (`markdown`,
  `confluence-storage`). Writing `--output <file>` emits `<file>.meta.json`.
- SonarCloud analysis runs in CI on every push/PR. The Quality Gate badge is **deferred**: QG is
  still ERROR on `new_duplicated_lines_density` (OD-3 revisit — do not advertise a failing gate).

## License

[MIT](LICENSE)
