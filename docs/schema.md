# Inventory schema

kollect-render validates **inventory evidence envelopes** against the committed JSON Schema at
[`schema/inventory-v0.schema.json`](https://github.com/PlatformRelay/kollect-render/blob/main/schema/inventory-v0.schema.json)
(draft v0, `$id` under this repository).

## Envelope shape (v0)

Required top-level fields (`schemaVersion` is the const `0.2.0`):

| Field | Role |
| --- | --- |
| `schemaVersion` | Must be `"0.2.0"` |
| `metadata` | `environmentId`, `sourceId`, `collectedAt` (optional `stale`) |
| `provenance` | Collection run identity, digests, tool versions |
| `targetOutcomes` | Per-target collection status |
| `components` | Platform / operator component instances |
| `helmReleases` | Helm release inventory |
| `workloads` | Workload kind/name/namespace/images |
| `nodes` | Node count and Kubernetes versions |

`additionalProperties` is false at every object level the schema defines — unknown keys fail
validation.

## Committed examples

Two synthetic, organization-neutral fixtures live under `schema/examples/`:

| File | Environment / source |
| --- | --- |
| [`region-a-cluster-alpha.yaml`](https://github.com/PlatformRelay/kollect-render/blob/main/schema/examples/region-a-cluster-alpha.yaml) | `region-a` / `cluster-alpha` |
| [`region-b-cluster-beta.yaml`](https://github.com/PlatformRelay/kollect-render/blob/main/schema/examples/region-b-cluster-beta.yaml) | `region-b` / `cluster-beta` |

Validate either from the repo root:

```bash
go run ./cmd/kollect-render validate schema/examples/region-a-cluster-alpha.yaml
go run ./cmd/kollect-render validate schema/examples/region-b-cluster-beta.yaml
```

## Render context vs inventory document

`validate` checks a **single inventory document**. `render` loads a **RenderContext** (manifest +
one or more docs + generation metadata + optional upstream/copy maps) — see the golden fixture at
`test/golden/env-inventory-md/context.yaml`. Private aggregators assemble that context; this CLI
does not invent missing sources.
