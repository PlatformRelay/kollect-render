# kollect-render

Credential-free CLI that validates versioned inventory documents and renders them to
deterministic artifacts (markdown, confluence-storage, and more as the format registry grows).
File outputs include a `.meta.json` digest/metadata sidecar for private publishers; this tool
never opens network connections or holds publish credentials.

This repository is the public Open Source renderer for the PlatformRelay inventory protocol.
Collection, scheduling, credentials, and publishing stay in private callers; this tool only
transforms documents you already have.

## Status

OSS renderer with kollect-parity quality gates (lint, arch-lint, coverage ≥90%, gitleaks,
REWE-trace). SonarCloud analysis runs in CI on every push/PR (badge deferred). Tagged releases
(`vX.Y.Z`) publish multi-platform binaries and a GHCR container image via GoReleaser.
`validate` and `render` are available; output formats register at compile time
(`markdown`, `confluence-storage`). Writing `--output <file>` emits `<file>.meta.json`
(content digest + generation metadata) for the private publisher contract.

## Quick start

```bash
task check          # fmt, vet, lint, arch-lint, coverage≥90%, gitleaks, build, REWE-trace
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

## CLI

Command and flag reference: [docs/cli-reference.md](docs/cli-reference.md).

## Golden tests

Fixture → expected output suites live under `test/golden/`. To regenerate committed goldens after an
intentional render change:

```bash
UPDATE_GOLDEN=1 go test ./internal/format/ ./internal/render/
```

Review the golden diff like code before committing.

## License

[MIT](LICENSE)
