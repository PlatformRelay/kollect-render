# CLI reference

`kollect-render` is a small command tree: `validate`, `render`, and `version`. Unknown commands
and fatal failures exit `2` with a message on stderr; success exits `0`.

## Commands

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
| `version` | Print the binary version string |

## Usage shapes

```text
kollect-render validate <document>
kollect-render render --format <markdown|confluence-storage> --context <file> \
  [--template <file>] [--output <file>] [--upstream-deps <file>] \
  [--generated-at <RFC3339>] [--report-origin <label>]
kollect-render version
```

## Behaviour notes

- Completeness / catalog policy stays with the private aggregator; this CLI only applies the
  flags above to an already-shaped render context.
- `--template` with `--format confluence-storage` is rejected (exit `2`, no output file).
- `--upstream-deps` **replaces** the whole `Upstream` map from context; it does not merge.
- `--generated-at` must be RFC3339 UTC; invalid values exit `2`.
- File output writes a sidecar `<path>.meta.json` with content digest (`sha256:…`) and generation
  metadata for private publishers.

See [Getting started](getting-started.md) for runnable examples against
`test/golden/env-inventory-md/context.yaml`.
