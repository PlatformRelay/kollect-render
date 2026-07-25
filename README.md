# kollect-render

Credential-free CLI that validates versioned inventory documents and renders them to
deterministic artifacts (markdown, confluence-storage, and more as the format registry grows).

This repository is the public Open Source renderer for the PlatformRelay inventory protocol.
Collection, scheduling, credentials, and publishing stay in private callers; this tool only
transforms documents you already have.

## Status

Scaffold (v0.0.0-dev). Schema validate/render commands land in follow-up work.

## Quick start

```bash
task check          # fmt, vet, test, build, REWE-trace gate
go run ./cmd/kollect-render version
```

Requires Go (see `go.mod`) and [Task](https://taskfile.dev/).

## License

[MIT](LICENSE)
