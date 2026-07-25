# Contributing to kollect-render

Thank you for helping improve kollect-render.

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development

1. Install Go (version in `go.mod`) and [Task](https://taskfile.dev/).
2. Run the local gate before opening a change:

   ```bash
   task check
   ```

   That runs `gofmt` check, `go vet`, unit tests, a binary build, and the **REWE-trace**
   gate (fails if forbidden internal markers appear in tracked files).

3. Keep the tree free of organization-specific identifiers. Runtime configuration belongs
   with callers — not in this repository.

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
