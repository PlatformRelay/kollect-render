# Getting started

Requires Go (see `go.mod`), [Task](https://taskfile.dev/), and
[gitleaks](https://github.com/gitleaks/gitleaks) on `PATH` for the secret-scan step of
`task check`.

## Build and version

From the repository root:

```bash
task check
go run ./cmd/kollect-render version
```

`task check` runs the full local gate (format, vet, lint, arch-lint, coverage ≥90%, gitleaks,
build, REWE-trace, pinned-actions). The `version` command prints the binary version
(`0.0.0-dev` for local `go run` builds).

## Validate an inventory document

Committed schema examples are valid inventory envelopes:

```bash
go run ./cmd/kollect-render validate schema/examples/region-a-cluster-alpha.yaml
```

Exit `0` on success; exit `2` on schema or I/O failure (message on stderr).

## Render against the golden context

The markdown golden fixture at `test/golden/env-inventory-md/context.yaml` is a complete
RenderContext. Use it to exercise both registered formats:

```bash
go run ./cmd/kollect-render render --format markdown --context test/golden/env-inventory-md/context.yaml
```

```bash
go run ./cmd/kollect-render render --format confluence-storage --context test/golden/env-inventory-md/context.yaml
```

Custom markdown templates are supported when `--format` is `markdown`:

```bash
go run ./cmd/kollect-render render --format markdown --template test/golden/env-inventory-md/template.md.tmpl --context test/golden/env-inventory-md/context.yaml
```

Omit `--output` (or pass `-`) to write the body to stdout. Writing `--output <file>` also emits
`<file>.meta.json` beside the body for the private publisher contract.

## Install from a release

Release artifacts (binaries + container) are published on each `vX.Y.Z` tag. Pull a pinned image
from GHCR, for example `ghcr.io/platformrelay/kollect-render:v0.1.0`, or download the matching
binary from GitHub Releases. See the [CLI reference](cli-reference.md) for every command and flag.
