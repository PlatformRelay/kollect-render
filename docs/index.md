# kollect-render

Credential-free CLI that validates versioned inventory documents and renders them to
deterministic artifacts (`markdown`, `confluence-storage`, and more as the format registry
grows). File outputs include a `.meta.json` digest/metadata sidecar for private publishers;
this tool never opens network connections or holds publish credentials.

This repository is the public Open Source renderer for the PlatformRelay inventory protocol.
Collection, scheduling, credentials, and publishing stay in private callers; this tool only
transforms documents you already have.

## Start here

- [Getting started](getting-started.md) — install, validate, and render against golden fixtures
- [CLI reference](cli-reference.md) — commands and flags (single source of truth)
- [Output formats](formats.md) — markdown and confluence-storage guarantees
- [Inventory schema](schema.md) — envelope v0 and committed examples
- [Development](development.md) — local gates and golden protocol
- [Security](security.md) — threat model and reporting

## Status

OSS renderer with kollect-parity quality gates (lint, arch-lint, coverage ≥90%, gitleaks,
REWE-trace, pinned Actions). Tagged releases (`vX.Y.Z`) publish multi-platform binaries and a
GHCR container image via GoReleaser.
