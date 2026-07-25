# kollect-render

Credential-free CLI that validates versioned inventory documents and renders them to
deterministic artifacts (markdown, confluence-storage, and more as the format registry grows).
File outputs include a `.meta.json` digest/metadata sidecar for private publishers; this tool
never opens network connections or holds publish credentials.

This repository is the public Open Source renderer for the PlatformRelay inventory protocol.
Collection, scheduling, credentials, and publishing stay in private callers; this tool only
transforms documents you already have.

## Status

v0.0.0-dev. `validate` and `render` are available; output formats register at compile time
(`markdown`, `confluence-storage`). Writing `--output <file>` emits `<file>.meta.json`
(content digest + generation metadata) for the private publisher contract.

## Quick start

```bash
task check          # fmt, vet, test, build, REWE-trace gate
go run ./cmd/kollect-render version
# Built-in Model → encoder (any registered format):
go run ./cmd/kollect-render render --format markdown --context test/golden/env-inventory-md/context.yaml
go run ./cmd/kollect-render render --format confluence-storage --context test/golden/env-inventory-md/context.yaml
# Custom markdown template ( --format must be markdown ):
go run ./cmd/kollect-render render --format markdown --template test/golden/env-inventory-md/template.md.tmpl \
  --context test/golden/env-inventory-md/context.yaml
```

Requires Go (see `go.mod`) and [Task](https://taskfile.dev/).

## CLI

| Command / flag | Role |
| --- | --- |
| `validate <document>` | Schema-validate an inventory document; exit `0` ok / `2` fatal |
| `render` | Deterministic render; exit `0` ok / `2` fatal |
| `--format` | `markdown` or `confluence-storage` (default `markdown`) |
| `--context` | RenderContext YAML/JSON (required) |
| `--template` | Optional markdown `text/template` (markdown format only) |
| `--output` | Body path, or `-` / omit for stdout (file output also writes `.meta.json`) |
| `--upstream-deps` | Replace `Upstream` with a `map[componentID]UpstreamEntry` YAML file |
| `--generated-at` | Override `Generation.GeneratedAt` (RFC3339 UTC) |
| `--report-origin` | Override `Generation.Origin` (e.g. `schedule`, `manual`, `ci`) |

Completeness / catalog policy stays with the private aggregator; this CLI only applies the
flags above to an already-shaped render context.

## Golden tests

Fixture → expected output suites live under `test/golden/`. To regenerate committed goldens after an
intentional render change:

```bash
UPDATE_GOLDEN=1 go test ./internal/format/ ./internal/render/
```

Review the golden diff like code before committing.

## License

[MIT](LICENSE)
