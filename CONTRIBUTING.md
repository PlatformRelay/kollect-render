# Contributing to kollect-render

Thank you for helping improve kollect-render.

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development

1. Install Go (version in `go.mod`) and [Task](https://taskfile.dev/).
2. Run the local gate before opening a change:

   ```bash
   task check
   ```

   That runs `gofmt` check, `go vet`, golangci-lint, go-arch-lint, unit tests with a
   **whole-module coverage floor of ≥90%**, gitleaks, a binary build, the **REWE-trace**
   gate (fails if forbidden internal markers appear in tracked files), and the
   **pinned-actions** gate (every third-party `uses:` in `.github/workflows/` must be a
   full 40-character commit SHA, e.g. `owner/repo@<sha> # vN.M.P`).

3. Keep the tree free of organization-specific identifiers. Runtime configuration belongs
   with callers — not in this repository.

4. When adding or bumping a GitHub Action, pin it to a commit SHA (not a mutable tag) and
   keep the human-readable version in a trailing comment so Renovate can maintain it.

5. Write the failing test first (TDD). Behaviour changes land with tests that failed before
   the implementation.

## Golden-file protocol

Fixture → expected output suites live under `test/golden/`. After an **intentional**
render or format change, regenerate committed goldens:

```bash
UPDATE_GOLDEN=1 go test ./internal/format/ ./internal/render/
```

Review the golden diff like product code before committing. Do not regenerate goldens only
to make CI green without understanding the behaviour change.

## Commits

Use **[Conventional Commits](https://www.conventionalcommits.org/)** with a mandatory ASCII
**[gitmoji](https://gitmoji.dev/)** shortcode prefix (e.g. `:sparkles:`, not Unicode emoji):

```text
:gitmoji: <type>(<optional scope>): <short summary>

<optional body — WHY / trade-offs>
```

Examples:

```text
:sparkles: feat: add validate command
:bug: fix: pointered schema errors
:memo: docs: clarify render flags
:construction_worker: ci: pin action SHAs
```

Types: `feat` `fix` `docs` `style` `refactor` `test` `chore` `ci` `build`.

One logical change per commit; keep the tree green at each. Do not add AI co-author trailers.
Never modify git config as part of a contribution.

## Pull requests

- Keep PRs focused and reviewable.
- Ensure `task check` is green locally before opening.
- **Merge policy: Rebase and merge** — linear history, no squash, no merge commits. Delete the
  branch on merge so each well-formed conventional commit lands on `main`.
- Fill in the PR template checklist (gates, TDD evidence, docs, goldens).

## Security

Report vulnerabilities privately — see [SECURITY.md](SECURITY.md).
