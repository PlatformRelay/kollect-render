# Development

## Prerequisites

- Go version from `go.mod`
- [Task](https://taskfile.dev/)
- [gitleaks](https://github.com/gitleaks/gitleaks) on `PATH` (pinned version noted in `Taskfile.yml`)
- [uv](https://docs.astral.sh/uv/) for the docs toolchain (optional unless you build docs)

## Local gate matrix

| Task | What it enforces |
| --- | --- |
| `task fmt-check` | `gofmt` clean |
| `task vet` | `go vet ./...` |
| `task lint` | golangci-lint |
| `task arch-lint` | `go-arch-lint` against `.go-arch-lint.yml` |
| `task coverage` | Whole-module coverage floor (`COVERAGE_MIN`, default 90) |
| `task gitleaks` | Secret scan with `.github/gitleaks.toml` |
| `task build` | `bin/kollect-render` |
| `task rewe-trace` | No forbidden internal markers in tracked files |
| `task pinned-actions` | Every workflow `uses:` is a full commit SHA |
| `task check` | All of the above, in order |

Run `task check` before opening a PR. See also [CONTRIBUTING.md](https://github.com/PlatformRelay/kollect-render/blob/main/CONTRIBUTING.md).

## Documentation site

| Task | Role |
| --- | --- |
| `task docs:build` | Install the locked MkDocs toolchain and build `site/` with `--strict` |
| `task docs:check` | Same as `docs:build` — broken nav/page links fail the build |
| `task docs:lock` | Regenerate `docs/requirements-docs.txt` from `docs/requirements-docs.in` |

## Golden protocol

Fixture → expected output suites live under `test/golden/`. After an intentional render change:

```bash
UPDATE_GOLDEN=1 go test ./internal/format/ ./internal/render/
```

Review the golden diff like code before committing. Do not regenerate goldens to “make CI green”
without understanding the behaviour change.

## Releases

Before tagging `vX.Y.Z`, run the read-only [release gate](development/release.md). It proves
the candidate SHA is on protected `main`, required checks are green, and `CHANGELOG.md` has a
matching section — without creating a tag or publishing artifacts.

## Commits

Conventional commits with an ASCII gitmoji shortcode (no Unicode emoji, no AI co-author trailers).
One logical change per commit; prefer rebase merge so each commit lands on `main`.
