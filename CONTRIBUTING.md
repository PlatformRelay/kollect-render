# Contributing to kollect-render

Thank you for helping improve kollect-render.

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development

1. Install Go (version in `go.mod`) and [Task](https://taskfile.dev/).
2. Run the local gate before opening a change:

   ```bash
   task check
   ```

   That runs `gofmt` check, `go vet`, golangci-lint, go-arch-lint, unit tests with an
   **internal/ coverage floor of ≥90%**, gitleaks, a binary build, the **REWE-trace**
   gate (fails if forbidden internal markers appear in tracked files), and the
   **pinned-actions** gate (every third-party `uses:` in `.github/workflows/` must be a
   full 40-character commit SHA, e.g. `owner/repo@<sha> # vN.M.P`).

3. Keep the tree free of organization-specific identifiers. Runtime configuration belongs
   with callers — not in this repository.

4. When adding or bumping a GitHub Action, pin it to a commit SHA (not a mutable tag) and
   keep the human-readable version in a trailing comment so Renovate can maintain it.

## Commits

Use conventional commits with an ASCII gitmoji shortcode:

```text
:sparkles: feat: add validate command
:bug: fix: pointered schema errors
:memo: docs: clarify render flags
```

One logical change per commit. Do not add AI co-author trailers.

## Pull requests

- Keep PRs focused and reviewable.
- Ensure `task check` is green.
- Prefer rebase merge so each conventional commit lands on `main`.

## Security

Report vulnerabilities privately — see [SECURITY.md](SECURITY.md).
